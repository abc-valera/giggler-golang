package errutil

import (
	"bytes"
	"errors"
	"fmt"
)

type internal struct {
	caller string // Caller provides additional context about error's location
	err    error  // Err is a nested error
}

func NewInternalErr(err error) error {
	return &internal{
		caller: caller(2),
		err:    err,
	}
}

func NewInternalString(err string) error {
	return &internal{
		caller: caller(2),
		err:    errors.New(err),
	}
}

func (e *internal) Error() string {
	var buf bytes.Buffer

	fmt.Fprintf(&buf, "%s ", e.caller)

	fmt.Fprintf(&buf, "<%s> ", CodeInternal)
	if e.err != nil {
		buf.WriteString(e.err.Error())
	}

	return buf.String()
}
