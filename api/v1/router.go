package v1

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/jsm/gode/api/v1/auth"
)

// Router for /v1
func Router() http.Handler {
	// Initialize Router
	r := chi.NewRouter()

	// Define routes
	r.Mount("/auth", auth.Router())

	// Return router
	return r
}
