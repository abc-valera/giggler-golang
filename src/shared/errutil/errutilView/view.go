package errutilView

import (
	"context"

	"giggler-golang/src/shared/errutil"
	"giggler-golang/src/shared/log"
	"giggler-golang/src/shared/view/viewgen"
)

type Handler struct{}

func (Handler) NewError(ctx context.Context, err error) *viewgen.ErrorSchemaStatusCode {
	var codeError viewgen.ErrorSchema
	if errutil.ErrorCode(err) == errutil.CodeInternal {
		codeError = viewgen.ErrorSchema{ErrorMessage: "Internal error"}
	} else {
		codeError = viewgen.ErrorSchema{ErrorMessage: err.Error()}
	}

	switch errutil.ErrorCode(err) {
	case errutil.CodeInvalidArgument, errutil.CodeNotFound, errutil.CodeAlreadyExists:
		return &viewgen.ErrorSchemaStatusCode{
			StatusCode: 400,
			Response:   codeError,
		}
	case errutil.CodeUnauthenticated:
		return &viewgen.ErrorSchemaStatusCode{
			StatusCode: 401,
			Response:   codeError,
		}
	case errutil.CodePermissionDenied:
		return &viewgen.ErrorSchemaStatusCode{
			StatusCode: 403,
			Response:   codeError,
		}
	default:
		log.Error("REQUEST_ERROR", "err", err.Error())
		return &viewgen.ErrorSchemaStatusCode{
			StatusCode: 500,
			Response:   codeError,
		}
	}
}
