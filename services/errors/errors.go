package serviceerrors

import (
	"fmt"
	"net/http"
	"strings"
)

type InvalidFirebaseKeyError string

func (err InvalidFirebaseKeyError) Error() string {
	return fmt.Sprintf("Invalid firebase key: %s", string(err))
}

type RoutinePanicError string

func (err RoutinePanicError) Error() string {
	return fmt.Sprintf("Panic in Goroutine: %s", string(err))
}

type MultiError []error

func (errs MultiError) Error() string {
	errStrings := []string{"Multiple Errors Occured"}
	for i, err := range errs {
		errStrings = append(errStrings, fmt.Sprintf("[Error %d]", i))
		errStrings = append(errStrings, err.Error())
	}
	return strings.Join(errStrings, "\n")
}

// NoUpdatesError for cases when no updates were sent
type NoUpdatesError string

func (err NoUpdatesError) Error() string {
	return fmt.Sprintf("No updates sent for %s", string(err))
}

func (err NoUpdatesError) HTTPCode() int {
	return http.StatusBadRequest
}

func (err NoUpdatesError) HTTPResponse() []error {
	return []error{err}
}

type ForbiddenUser uint

func (err ForbiddenUser) Error() string {
	return "Not authorized"
}

func (err ForbiddenUser) HTTPCode() int {
	return http.StatusForbidden
}

func (err ForbiddenUser) HTTPResponse() []error {
	return []error{err}
}
