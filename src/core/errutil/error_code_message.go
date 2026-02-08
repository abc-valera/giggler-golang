package errutil

import (
	"bytes"
	"fmt"
)

type codeError struct {
	code Code  // code provides general information about the error
	err  error // err is a nested error
}

func NewCode(code Code, err error) error {
	return codeError{
		code: code,
		err:  err,
	}
}

func (e codeError) Error() string {
	var buf bytes.Buffer

	if e.code != "" {
		fmt.Fprintf(&buf, "<%s> ", e.code)
	}
	buf.WriteString(e.err.Error())

	return buf.String()
}
