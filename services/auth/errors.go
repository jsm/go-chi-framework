package auth

import (
	"fmt"

	"github.com/jsm/gode/utils"
)

type ErrEmailAlreadyRegistered string

func (err ErrEmailAlreadyRegistered) Error() string {
	return fmt.Sprintf("The Email %s is already registered", string(err))
}

type ErrNoUserForEmail string

func (err ErrNoUserForEmail) Error() string {
	return fmt.Sprintf("There is no users associated with the email: %s", string(err))
}

type ErrWrongPassword utils.EmptyStruct

func (err ErrWrongPassword) Error() string {
	return "Wrong Password"
}

type ErrInvalidSSOProvider string

func (err ErrInvalidSSOProvider) Error() string {
	return fmt.Sprintf("The provider %s is invalid", string(err))
}
