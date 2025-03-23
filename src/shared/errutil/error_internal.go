package errutil

import (
	"bytes"
	"fmt"
)

type internalError struct {
	caller string // caller provides additional context about error's location
	err    error  // err is a nested error
}

func NewInternal(err error) error {
	return internalError{
		caller: caller(2),
		err:    err,
	}
}

func (e internalError) Error() string {
	var buf bytes.Buffer

	fmt.Fprintf(&buf, "%s ", e.caller)

	fmt.Fprintf(&buf, "<%s> ", CodeInternal)
	if e.err != nil {
		buf.WriteString(e.err.Error())
	}

	return buf.String()
}
