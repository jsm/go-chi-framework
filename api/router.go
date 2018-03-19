package main

import (
	"fmt"
	"net/http"
	"os"
	"runtime/debug"
	"strings"

	"github.com/getsentry/raven-go"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"

	"github.com/jsm/gode/api/application"
	"github.com/jsm/gode/api/v1"
)

const maxBytes = 50000000

func router() http.Handler {
	// Initialize Main Router
	r := chi.NewRouter()

	c := cors.New(cors.Options{
		AllowedOrigins:   strings.Split(os.Getenv("ALLOWED_ORIGINS"), ","),
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
	})

	// Setup Middleware
	r.Use(limitMaxBytes)
	r.Use(c.Handler)
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.StripSlashes)
	r.Use(setContentJSON)
	r.Use(recoverer)

	// Custom NotFound Handler
	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		respJ := application.CreateResponseJSON(false, []error{fmt.Errorf("Route Not Found")}, map[string]string{
			"path": r.URL.Path,
		})

		w.WriteHeader(http.StatusNotFound)
		w.Write(respJ)
	})

	// Define routes
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write(application.DefaultOKResponse)
	})

	r.Mount("/v1", v1.Router())

	// Return router
	return r
}

func limitMaxBytes(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		r.Body = http.MaxBytesReader(w, r.Body, maxBytes)
		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}

func setContentJSON(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}

// Custom recovery middleware. Combines Sentry Handler as well as Chi's recoverer handler, with our custom json response.
// https://github.com/go-chi/chi/blob/master/middleware/recoverer.go
// https://github.com/getsentry/raven-go/blob/master/http.go#L70
func recoverer(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rvr := recover(); rvr != nil {

				logEntry := middleware.GetLogEntry(r)
				if logEntry != nil {
					logEntry.Panic(rvr, debug.Stack())
				} else {
					fmt.Fprintf(os.Stderr, "Panic: %+v\n", rvr)
					debug.PrintStack()
				}

				// Send to Sentry
				if application.Instance.IsSentryEnabled {
					rvrStr := fmt.Sprint(rvr)
					packet := raven.NewPacket(rvrStr, raven.NewException(fmt.Errorf(rvrStr), raven.NewStacktrace(2, 3, nil)), raven.NewHttp(r))
					raven.Capture(packet, nil)
				}

				http.Error(w, application.InternalServerErrorResponse, http.StatusInternalServerError)
			}
		}()

		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}
