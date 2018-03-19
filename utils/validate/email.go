package validate

import (
	"regexp"
)

var emailFormat = regexp.MustCompile(`^[^\s@]+@[^\s@]+\.[^\s@]+$`)

// Email validation
func Email(email string) (errs []error) {
	if !emailFormat.MatchString(email) {
		errs = append(errs, validationError("Email address format is invalid"))
	}

	return errs
}
