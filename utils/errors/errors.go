package errors

import (
	"github.com/getsentry/raven-go"
	"github.com/pkg/errors"
)

type Extra raven.Extra

var New = errors.New
var WithStack = errors.WithStack

func WithExtra(err error, extra map[string]interface{}) error {
	return WithStack(raven.WrapWithExtra(err, extra))
}
