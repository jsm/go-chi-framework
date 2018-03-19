package validate

import (
	"fmt"
	"net/http"
	"strings"
)

// Generate a validation error
func validationError(message string) error {
	return fmt.Errorf("ValidationError: %s", message)
}

//InvalidArgumentsError checks for invalid arguments
type InvalidArgumentsError []error

func (errors InvalidArgumentsError) Error() string {
	var invalidArgs []string
	for _, err := range errors {
		invalidArgs = append(invalidArgs, err.Error())
	}
	return fmt.Sprintf("Invalid arguments: %s", strings.Join(invalidArgs, ", "))
}

func (errors InvalidArgumentsError) HTTPCode() int {
	return http.StatusBadRequest
}

func (errors InvalidArgumentsError) HTTPResponse() []error {
	return errors
}
