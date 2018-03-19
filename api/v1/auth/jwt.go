package auth

import (
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/jsm/gode/api/application"
)

var jwtSecret []byte
var errNoToken = errors.New("This endpoint requires Token authorization")
var errParsingToken = errors.New("There was an error while parsing the provided Token")

type Authorizations struct {
	UserID    uint
	ExpiresAt time.Time
}

func initializeJWTSecret() {
	jwtSecretString := os.Getenv("JWT_SECRET")

	var err error
	jwtSecret, err = base64.StdEncoding.DecodeString(jwtSecretString)
	if err != nil {
		panic(err)
	}
}

func Authenticate(w http.ResponseWriter, r *http.Request) (Authorizations, bool) {
	var authorizations Authorizations

	cookie, err := r.Cookie("auth_token")
	if err != nil {
		jsonResponse := application.CreateResponseJSON(false, []error{errNoToken}, nil)
		w.WriteHeader(http.StatusUnauthorized)
		w.Write(jsonResponse)
		return authorizations, false
	}

	tokenString := cookie.Value
	fmt.Println(tokenString)
	if tokenString == "" {
		jsonResponse := application.CreateResponseJSON(false, []error{errNoToken}, nil)
		w.WriteHeader(http.StatusUnauthorized)
		w.Write(jsonResponse)
		return authorizations, false
	}

	// Parse token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return jwtSecret, nil
	})

	if err != nil {
		jsonResponse := application.CreateResponseJSON(false, []error{errParsingToken}, nil)
		w.WriteHeader(http.StatusUnauthorized)
		w.Write(jsonResponse)
		return authorizations, false
	}

	// Extract Claims
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		jsonResponse := application.CreateResponseJSON(false, []error{errParsingToken}, nil)
		w.WriteHeader(http.StatusUnauthorized)
		w.Write(jsonResponse)
		return authorizations, false
	}

	// Parse Authorizations into struct
	authorizations.UserID = uint(claims["userID"].(float64))
	authorizations.ExpiresAt = time.Unix(int64(claims["expiresAt"].(float64)), 0)

	// Check Expiration
	if authorizations.ExpiresAt.Before(time.Now()) {
		return authorizations, false
	}

	return authorizations, true
}
