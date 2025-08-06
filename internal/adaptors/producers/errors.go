package producers

import "golang-skeleton/pkg/errors"

// ProducerError represent the type of errors returned by the producers package.
type ProducerError struct {
	error
	Code int
	Msg  string
}

var (
	ErrMsgValEncode     = errors.Msg(errors.CodeValEncode)
	ErrMsgKeyEncode     = errors.Msg(errors.CodeKeyEncode)
	ErrMsgResourceWrite = errors.Msg(errors.CodeResWFailed)
)
