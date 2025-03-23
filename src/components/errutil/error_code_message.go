package errutil

import (
	"bytes"
	"fmt"
)

type codeMessage struct {
	code    Code   // Code provides general information about the error
	message string // Message provides additional context in human-readable form
}

func NewCodeMessage(code Code, message string) error {
	return &codeMessage{
		code:    code,
		message: message,
	}
}

func NewCodeError(code Code, err error) error {
	return &codeMessage{
		code:    code,
		message: err.Error(),
	}
}

func NewCodeMessageError(code Code, message string, err error) error {
	return &codeMessage{
		code:    code,
		message: fmt.Sprintf("%s: %s", message, err.Error()),
	}
}

// Error returns the string representation of the Error
func (e codeMessage) Error() string {
	var buf bytes.Buffer

	if e.code != "" {
		fmt.Fprintf(&buf, "<%s> ", e.code)
	}
	buf.WriteString(e.message)

	return buf.String()
}
