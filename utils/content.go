package utils

import (
	"fmt"
	"io"

	"github.com/pkg/errors"
	filetype "gopkg.in/h2non/filetype.v1"
	"gopkg.in/h2non/filetype.v1/types"
)

const defaultType = "application/octet-stream"

var errUnknownType = fmt.Errorf("Unknown file type")

type readSeekFile interface {
	io.Reader
	io.ReaderAt
	io.Seeker
}

// DetectFileContentType returns a MIME type detected from the given file
func DetectFileContentType(file readSeekFile) (string, error) {
	kind, err := filetype.MatchReader(file)

	if err != nil {
		return defaultType, errors.WithStack(err)
	}

	if kind == types.Unknown {
		return defaultType, errors.WithStack(errUnknownType)
	}

	return kind.MIME.Value, nil
}
