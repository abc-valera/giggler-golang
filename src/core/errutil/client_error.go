package errutil

import (
	"errors"
	"fmt"
)

var (
	ErrorAuthentication = errors.New("authentication failed")
	ErrorAuthorization  = errors.New("access denied")
	ErrorValidation     = errors.New("invalid input")
	ErrorConflict       = errors.New("resource conflict")
	ErrorNotFound       = errors.New("resource not found")
	ErrorRateLimit      = errors.New("rate limit exceeded")
)

func Wrap(wrapped error, wrapper error) error {
	return fmt.Errorf("%v: %w", wrapper, wrapped)
}

func WrapStr(wrappedStr string, wrapper error) error {
	return fmt.Errorf("%v: %w", wrapper, errors.New(wrappedStr))
}

// Usage example:
// 	ormError := errors.New("user with id 42 doesn't exist")
// 	return fmt.Errorf("%v: %w", NotFoundError, returnedErr)

// TODO: all the errors in the app should be created from these using
// either errors.Join(FundamentalError, SpecificError)
// or by wrapping them using fmt.Errorf("%w: user with such ID already exists", FundamentalError)

// The errors will then be checked using errors.Is(err, errutil.AuthenticationError)
// to handler specific fundamental cases

// Also check how to use errors.Unwrap method
