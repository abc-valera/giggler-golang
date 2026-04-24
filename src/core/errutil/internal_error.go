package errutil

import (
	"errors"
)

var (
	ErrorInternalServer     = errors.New("internal logic error")
	ErrorNotImplemented     = errors.New("not implemented")
	ErrorServiceUnavailable = errors.New("service unavailable")
	ErrorResourceExhausted  = errors.New("resource exhausted")
)

// TODO: all the errors in the app should be created from these using
// either errors.Join(FundamentalError, SpecificError)
// or by wrapping them using fmt.Errorf("%w: user with such ID already exists", FundamentalError)
