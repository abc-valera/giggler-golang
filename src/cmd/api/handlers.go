package main

import (
	"context"

	"github.com/abc-valera/giggler-golang/src/features/auth/authView"
	"github.com/abc-valera/giggler-golang/src/shared/contexts"
	"github.com/abc-valera/giggler-golang/src/shared/errutil"
	"github.com/abc-valera/giggler-golang/src/shared/log"
	"github.com/abc-valera/giggler-golang/src/shared/view/viewgen"
)

type errorHandler struct{}

func (errorHandler) NewError(ctx context.Context, err error) *viewgen.ErrorSchemaStatusCode {
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

type securityHandler struct{}

func (securityHandler) HandleBearerAuth(ctx context.Context, _ string, t viewgen.BearerAuth) (context.Context, error) {
	payload, err := authView.ViewVerifyToken(t.Token)
	if err != nil {
		return ctx, err
	}

	return contexts.SetUserID(ctx, payload.UserID), nil
}
