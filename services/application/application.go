package application

import "github.com/op/go-logging"

// Application Interface
type Application interface {
	CaptureError(mainErr error, context map[string]string)
	CaptureErrors(errs []error, context map[string]string)
	Logger() *logging.Logger
}
