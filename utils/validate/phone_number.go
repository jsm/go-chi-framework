package validate

import "fmt"

// PhoneNumber validation
func PhoneNumber(phoneNumber uint) (errs []error) {
	if phoneNumber == 0 {
		errs = append(errs, validationError("Phone number is empty"))
		return errs
	}

	if len(fmt.Sprint(phoneNumber)) != 10 {
		errs = append(errs, validationError("Phone number must be of length 10"))
	}

	return errs
}
