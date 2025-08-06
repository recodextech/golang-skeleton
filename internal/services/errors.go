package services

import (
	"golang-skeleton/pkg/errors"
)

// ServiceError represent the type of custom errors returned by the services.
type ServiceError struct {
	error
	Code int
	Msg  string
}

var (
	ErrMsgResNotFound = errors.Msg(errors.CodeResNotFound)
)
