package errutil

import (
	"errors"
)

var (
	ErrInternalServer     = errors.New("internal logic error")
	ErrNotImplemented     = errors.New("not implemented")
	ErrServiceUnavailable = errors.New("service unavailable")
	ErrResourceExhausted  = errors.New("resource exhausted")
)
