package errutil

import "errors"

const (
	// VALIDATION CODES

	// CodeInvalidArgument is a validation error. Should be returned with all failed validation checks.
	CodeInvalidArgument Code = "invalid_argument"

	// DATASOURCE CODES

	// CodeNotFound means requested entity was not found
	CodeNotFound Code = "not_found"

	// CodeAlreadyExists means an attempt to create an entity failed because one already exists
	CodeAlreadyExists Code = "already_exists"

	// ROLES CODES

	// CodeUnauthenticated indicates the request does not have valid authentication credentials for the operation
	CodeUnauthenticated Code = "unauthenticated"

	// CodePermissionDenied indicates the caller does not have permission to execute the specified operation
	CodePermissionDenied Code = "permission_denied"

	// INTERNAL CODES

	// CodeCanceled means the operation was canceled (either by the caller or the system)
	CodeCanceled Code = "canceled"

	// CodeInternal means an internal error occured
	CodeInternal Code = "internal"
)

type Code string

// ErrorCode returns the code of the root error, if available, otherwise returns Internal.
func ErrorCode(err error) Code {
	var messageErrTarget *codeError
	var internalErrTarget *internalError
	if err == nil {
		return ""
	} else if e, ok := err.(*codeError); ok && e.code != "" {
		return e.code
	} else if _, ok := err.(*internalError); ok {
		return CodeInternal
	} else if errors.As(err, &messageErrTarget) {
		return messageErrTarget.code
	} else if errors.As(err, &internalErrTarget) {
		return CodeInternal
	}
	return CodeInternal
}
