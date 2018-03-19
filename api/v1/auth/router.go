package auth

import (
	"net/http"

	"github.com/go-chi/chi"
)

// Router for /v1/auth
func Router() http.Handler {
	// Initialize router
	r := chi.NewRouter()
	initializeJWTSecret()

	r.Get("/test", testHandler)
	r.Post("/login_or_signup_email", loginOrSignupEmailHandler)
	r.Post("/signup_email", signupEmailHandler)
	r.Post("/login_email", loginEmailHandler)
	r.Post("/login_or_signup_sso", loginOrSignupSSOHandler)

	// Return router
	return r
}
