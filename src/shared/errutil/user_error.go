package errutil

import (
	"errors"
	"fmt"
)

var (
	ErrAuthentication = errors.New("authentication failed")
	ErrAuthorization  = errors.New("access denied")
	ErrValidation     = errors.New("invalid input")
	ErrConflict       = errors.New("resource conflict")
	ErrNotFound       = errors.New("resource not found")
	ErrRateLimit      = errors.New("rate limit exceeded")
)

func Wrap(base error, specific error) error {
	return fmt.Errorf("%v: %w", specific, base)
}

func WrapStr(baseStr string, specific error) error {
	return fmt.Errorf("%v: %w", specific, errors.New(baseStr))
}
