package interfaces

import (
	"time"

	"github.com/jsm/gode/models"
)

type AuthService interface {
	Connect()
	IsEmailRegistered(email string) (bool, error)
	RegisterEmail(email string, password string) (models.User, string, time.Time, error)
	LoginEmail(email string, password string) (models.User, string, time.Time, error)
	LoginOrSignupSSO(token string, provider string) (models.User, string, time.Time, error)
}
