package auth

import (
	"encoding/base64"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"

	"github.com/jsm/gode/models"
	"github.com/jsm/gode/services/application"
	"github.com/jsm/gode/services/interfaces"
	"github.com/jsm/gode/utils/errors"
	"github.com/jsm/gode/utils/pge"
)

type instance struct {
	app           application.Application
	db            *gorm.DB
	jwtSecret     []byte
	jwtExpiration time.Duration
	services      *services
}

type services struct {
}

// Initialize the Auth Service
func Initialize(app application.Application, db *gorm.DB) interfaces.AuthService {
	jwtSecretString := os.Getenv("JWT_SECRET")
	jwtSecret, err := base64.StdEncoding.DecodeString(jwtSecretString)
	if err != nil {
		panic(err)
	}

	jwtExpirationString := os.Getenv("JWT_EXPIRATION_MINUTES")
	jwtExpirationNum, err := strconv.ParseInt(jwtExpirationString, 10, 64)
	if err != nil {
		panic(err)
	}
	jwtExpiration := time.Minute * time.Duration(jwtExpirationNum)

	return &instance{
		app:           app,
		db:            db,
		jwtSecret:     jwtSecret,
		jwtExpiration: jwtExpiration,
	}
}

func (i *instance) Connect() {
	i.services = &services{}
}

func (i *instance) generateJWT(user models.User) (string, time.Time, error) {
	expiresAt := time.Now().Add(i.jwtExpiration)

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, jwt.MapClaims{
		"userID":    user.ID,
		"expiresAt": expiresAt.Unix(),
	})

	tokenString, err := token.SignedString(i.jwtSecret)
	if err != nil {
		return "", expiresAt, errors.WithStack(err)
	}

	return tokenString, expiresAt, err
}

func (i *instance) IsEmailRegistered(email string) (bool, error) {
	var emailCount uint
	if err := i.db.Model(&models.EmailLogin{}).Where("email=?", email).Count(&emailCount).Error; err != nil {
		return false, errors.WithStack(err)
	}

	return emailCount > 0, nil
}

func (i *instance) RegisterEmail(email string, password string) (models.User, string, time.Time, error) {
	var user models.User
	var expiresAt time.Time

	alreadyRegistered, err := i.IsEmailRegistered(email)
	if err != nil {
		return user, "", expiresAt, err
	}

	if alreadyRegistered {
		return user, "", expiresAt, ErrEmailAlreadyRegistered(email)
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return user, "", expiresAt, errors.WithStack(err)
	}

	if err := i.db.Create(&user).Error; err != nil {
		return user, "", expiresAt, errors.WithStack(err)
	}

	emailLogin := models.EmailLogin{
		UserID:       user.ID,
		Email:        email,
		PasswordHash: string(hashedPassword),
	}
	if err := i.db.Create(&emailLogin).Error; err != nil {
		return user, "", expiresAt, errors.WithStack(err)
	}

	tokenString, expiresAt, err := i.generateJWT(user)
	if err != nil {
		return user, "", expiresAt, err
	}

	return user, tokenString, expiresAt, nil
}

func (i *instance) LoginEmail(email string, password string) (models.User, string, time.Time, error) {
	var user models.User
	var expiresAt time.Time

	var emailLogin models.EmailLogin
	if err := i.db.Where("email=?", email).First(&emailLogin).Error; err != nil {
		if strings.Contains(err.Error(), pge.RecordNotFound) {
			return user, "", expiresAt, ErrNoUserForEmail(email)
		}
		return user, "", expiresAt, errors.WithStack(err)
	}

	err := bcrypt.CompareHashAndPassword([]byte(emailLogin.PasswordHash), []byte(password))
	if err != nil {
		return user, "", expiresAt, errors.WithStack(ErrWrongPassword{})
	}

	tokenString, expiresAt, err := i.generateJWT(user)
	if err != nil {
		return user, "", expiresAt, err
	}

	if err := i.db.First(&user, emailLogin.UserID).Error; err != nil {
		return user, "", expiresAt, errors.WithStack(err)
	}

	return user, tokenString, expiresAt, nil
}

func (i *instance) IsSSORegistered(ssoID string, ssoType string) (bool, uint, error) {
	var ssoLogin models.SSOLogin
	if err := i.db.Where("sso_id=? AND sso_type=?", ssoID, ssoType).First(&ssoLogin).Error; err != nil {
		if strings.Contains(err.Error(), pge.RecordNotFound) {
			return false, 0, nil
		}
		return false, 0, errors.WithStack(err)
	}

	return true, ssoLogin.UserID, nil
}

func (i *instance) loginOrSignupSSO(ssoID string, ssoType string) (models.User, string, time.Time, error) {
	var user models.User
	var expiresAt time.Time

	alreadyRegistered, userID, err := i.IsSSORegistered(ssoID, ssoType)
	if err != nil {
		return user, "", expiresAt, err
	}

	if alreadyRegistered {
		if err := i.db.First(&user, userID).Error; err != nil {
			return user, "", expiresAt, errors.WithStack(err)
		}
	} else {
		tx := i.db.Begin()

		if err := tx.Create(&user).Error; err != nil {
			tx.Rollback()
			return user, "", expiresAt, errors.WithStack(err)
		}

		ssoLogin := models.SSOLogin{
			UserID:  user.ID,
			SSOID:   ssoID,
			SSOType: ssoType,
		}
		if err := tx.Create(&ssoLogin).Error; err != nil {
			tx.Rollback()
			return user, "", expiresAt, errors.WithStack(err)
		}

		if err := tx.Commit().Error; err != nil {
			return user, "", expiresAt, errors.WithStack(err)
		}
	}

	tokenString, expiresAt, err := i.generateJWT(user)
	if err != nil {
		return user, "", expiresAt, err
	}

	return user, tokenString, expiresAt, nil
}

func (i *instance) loginOrSignupGoogle(token string) (models.User, string, time.Time, error) {
	var user models.User

	tokenInfo, err := getGoogleTokenInfo(token)
	if err != nil {
		return user, "", time.Time{}, err
	}

	return i.loginOrSignupSSO(tokenInfo.UserID, models.SSOLoginTypes.Google)
}

func (i *instance) loginOrSignupLinkedIn(token string) (models.User, string, time.Time, error) {
	var user models.User

	userInfo, err := getLinkedInUserInfo(token)
	if err != nil {
		return user, "", time.Time{}, err
	}

	return i.loginOrSignupSSO(userInfo.ID, models.SSOLoginTypes.LinkedIn)
}

func (i *instance) LoginOrSignupSSO(token string, provider string) (models.User, string, time.Time, error) {
	switch provider {
	case models.SSOLoginTypes.Google:
		return i.loginOrSignupGoogle(token)
	case models.SSOLoginTypes.LinkedIn:
		return i.loginOrSignupLinkedIn(token)
	}

	return models.User{}, "", time.Time{}, ErrInvalidSSOProvider(provider)
}
