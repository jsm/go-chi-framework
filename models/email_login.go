package models

type EmailLogin struct {
	Base
	User         User
	UserID       uint
	Email        string
	PasswordHash string
}
