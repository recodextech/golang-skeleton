package uuid

import (
	"golang-skeleton/pkg/errors"
)

// ParseError represent the type of errors returned by the uuid package.
type ParseError struct {
	error
	Msg string
}

var (
	ErrMsgInvalidUUID = errors.Msg(errors.CodeMalformedReq)
)
