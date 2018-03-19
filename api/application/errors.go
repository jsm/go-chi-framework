package application

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/pkg/errors"
)

type HTTPError interface {
	HTTPCode() int
	HTTPResponse() []error
}

// ErrInternalServer Default error for internal failures
var ErrInternalServer = fmt.Errorf("Internal Server Error")

// ErrBodyExpected Default error when expecting a request body
var ErrBodyExpected = fmt.Errorf("Expected Request Body")

// InternalServerErrorResponse Default internal server error message
var InternalServerErrorResponse = string(CreateResponseJSON(false, []error{ErrInternalServer}, nil))

func formatError(mainErr error, context map[string]string) string {
	if context != nil {
		contextJSON, err := json.Marshal(context)
		if err == nil {
			return fmt.Sprintf("Error: %+v\nError Context: %s", mainErr, string(contextJSON))
		}
		return fmt.Sprintf("Error: %+v\nError Context: %s", mainErr, "PARSINGFAILED")
	}
	return fmt.Sprintf("Error: %+v", mainErr)
}

// CaptureWarning and route to appropriate mechanisms
func CaptureWarning(mainErr error, context map[string]string) {
	// Locally, print out the warning
	if Env.IsLocal {
		Instance.Log.Warning(formatError(mainErr, context))
	}
}

// CaptureWarnings and route to appropriate mechanisms
func CaptureWarnings(errs []error, context map[string]string) {
	for _, err := range errs {
		CaptureWarning(err, context)
	}
}

// CaptureError and route to appropriate mechanisms
func CaptureError(mainErr error, context map[string]string) {
	Instance.CaptureError(mainErr, context)
}

// CaptureErrors and route to appropriate mechanisms
func CaptureErrors(errs []error, context map[string]string) {
	Instance.CaptureErrors(errs, context)
}

// HandleJSONDecodeError default response for JSON Decode errors
func HandleJSONDecodeError(err error, errContext map[string]string, w http.ResponseWriter) {
	// Return 400 if it's an expected Decode error
	switch err.(type) {
	case *json.UnmarshalFieldError, *json.UnmarshalTypeError, *json.InvalidUnmarshalError:
		jsonResponse := CreateResponseJSON(false, []error{err}, nil)

		w.WriteHeader(http.StatusBadRequest)
		w.Write(jsonResponse)
		CaptureError(err, errContext)

		return

	case *json.SyntaxError:
		jsonResponse := CreateResponseJSON(false, []error{fmt.Errorf("Invalid JSON for request body")}, nil)

		w.WriteHeader(http.StatusBadRequest)
		w.Write(jsonResponse)
		CaptureError(err, errContext)

		return

	default:
		HandleInternalServerError(err, errContext, w)
	}
}

// HandleInternalServerError default response for internal server errors
func HandleInternalServerError(err error, errContext map[string]string, w http.ResponseWriter) {
	Instance.CaptureError(err, errContext)
	http.Error(w, InternalServerErrorResponse, http.StatusInternalServerError)
}

// HandleInternalServerErrors default response for internal server errors
func HandleInternalServerErrors(errs []error, errContext map[string]string, w http.ResponseWriter) {
	Instance.CaptureErrors(errs, errContext)
	http.Error(w, InternalServerErrorResponse, http.StatusInternalServerError)
}

func HandleError(err error, errContext map[string]string, w http.ResponseWriter) {
	if httpErr, ok := errors.Cause(err).(HTTPError); ok {
		respJ := CreateResponseJSON(false, httpErr.HTTPResponse(), nil)

		w.WriteHeader(httpErr.HTTPCode())
		w.Write(respJ)
	} else {
		HandleInternalServerError(err, errContext, w)
	}
}
