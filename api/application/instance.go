package application

import (
	"github.com/getsentry/raven-go"
	"github.com/op/go-logging"
	"github.com/pkg/errors"

	"github.com/jsm/gode/services/errors"
)

// CaptureErrors and route to appropriate mechanisms
func (app *Application) CaptureErrors(errs []error, context map[string]string) {
	for _, err := range errs {
		app.CaptureError(err, context)
	}
}

// CaptureError and route to appropriate mechanisms
func (app *Application) CaptureError(mainErr error, context map[string]string) {
	if mainErr == nil {
		return
	}

	if errs, ok := errors.Cause(mainErr).(serviceerrors.MultiError); ok {
		app.CaptureErrors(errs, context)
		return
	}

	// Send to Sentry
	if app.IsSentryEnabled {
		raven.CaptureError(mainErr, context)
	}

	// Locally, print out the error
	if Env.IsLocal {
		app.Log.Error(formatError(mainErr, context))
	}
}

// Logger for app instance
func (app *Application) Logger() *logging.Logger {
	return app.Log
}
