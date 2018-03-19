package auth

import (
	"net/http"
	"time"

	"github.com/jsm/gode/api/application"
	"github.com/jsm/gode/api/v1/requests"
	"github.com/jsm/gode/models"
	"github.com/jsm/gode/services"
)

func testHandler(w http.ResponseWriter, r *http.Request) {
	authorizations, success := Authenticate(w, r)
	if !success {
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(application.CreateResponseJSON(true, nil, authorizations))
}

func userResponse(user models.User) map[string]interface{} {
	return map[string]interface{}{
		"id": user.ID,
	}
}

func loginOrSignupEmailHandler(w http.ResponseWriter, r *http.Request) {
	var request requests.AuthLoginOrSignupEmail
	if !application.HandleJSONDecode(&request, w, r, nil) {
		return
	}

	errContext := map[string]string{
		"email": request.Email,
	}

	alreadyRegistered, err := services.Auth.IsEmailRegistered(request.Email)
	if err != nil {
		application.HandleError(err, errContext, w)
		return
	}

	var action string
	if alreadyRegistered {
		action = "login"
	} else {
		action = "signup"
	}

	responseJSON := application.CreateResponseJSON(true, nil, map[string]interface{}{
		"email":  request.Email,
		"action": action,
	})

	w.WriteHeader(http.StatusOK)
	w.Write(responseJSON)
}

func successfulLogin(user models.User, jwt string, expiresAt time.Time, w http.ResponseWriter) {
	responseJSON := application.CreateResponseJSON(true, nil, map[string]interface{}{
		"user": userResponse(user),
	})

	cookie := &http.Cookie{
		Name:    "auth_token",
		Value:   jwt,
		Expires: expiresAt,
	}

	http.SetCookie(w, cookie)
	w.WriteHeader(http.StatusOK)
	w.Write(responseJSON)
}

func signupEmailHandler(w http.ResponseWriter, r *http.Request) {
	var request requests.AuthSignupEmail
	if !application.HandleJSONDecode(&request, w, r, nil) {
		return
	}

	errContext := map[string]string{
		"email": request.Email,
	}

	user, jwt, expiresAt, err := services.Auth.RegisterEmail(request.Email, request.Password)
	if err != nil {
		application.HandleError(err, errContext, w)
		return
	}

	successfulLogin(user, jwt, expiresAt, w)
}

func loginEmailHandler(w http.ResponseWriter, r *http.Request) {
	var request requests.AuthLoginEmail
	if !application.HandleJSONDecode(&request, w, r, nil) {
		return
	}

	errContext := map[string]string{
		"email": request.Email,
	}

	user, jwt, expiresAt, err := services.Auth.LoginEmail(request.Email, request.Password)
	if err != nil {
		application.HandleError(err, errContext, w)
		return
	}

	successfulLogin(user, jwt, expiresAt, w)
}

func loginOrSignupSSOHandler(w http.ResponseWriter, r *http.Request) {
	var request requests.AuthLoginOrSignupSSO
	if !application.HandleJSONDecode(&request, w, r, nil) {
		return
	}

	user, jwt, expiresAt, err := services.Auth.LoginOrSignupSSO(request.Token, request.Provider)
	if err != nil {
		application.HandleError(err, nil, w)
		return
	}

	successfulLogin(user, jwt, expiresAt, w)
}
