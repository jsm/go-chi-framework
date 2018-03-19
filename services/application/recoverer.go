package application

import (
	"fmt"

	"github.com/jsm/gode/services/errors"
	"github.com/pkg/errors"
)

func Recoverer(errChan chan error) {
	if r := recover(); r != nil {
		errChan <- errors.WithStack(serviceerrors.RoutinePanicError(fmt.Sprintf("%#v", r)))
	}
}
