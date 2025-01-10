package view

import (
	"context"

	"github.com/abc-valera/giggler-golang/gen/view"
	"github.com/abc-valera/giggler-golang/src/components/errutil"
	"github.com/abc-valera/giggler-golang/src/components/viewutil"
	"github.com/abc-valera/giggler-golang/src/features/auth"
	"github.com/abc-valera/giggler-golang/src/shared/logger"
)

type handler struct{}

func (handler) NewError(ctx context.Context, err error) *view.CodeErrorStatusCode {
	var codeError view.CodeError
	if errutil.ErrorCode(err) == errutil.CodeInternal {
		codeError = view.CodeError{
			ErrorMessage: "Internal error",
		}
	} else {
		codeError = view.CodeError{
			ErrorMessage: err.Error(),
		}
	}

	switch errutil.ErrorCode(err) {
	case errutil.CodeInvalidArgument, errutil.CodeNotFound, errutil.CodeAlreadyExists:
		return &view.CodeErrorStatusCode{
			StatusCode: 400,
			Response:   codeError,
		}
	case errutil.CodeUnauthenticated:
		return &view.CodeErrorStatusCode{
			StatusCode: 401,
			Response:   codeError,
		}
	case errutil.CodePermissionDenied:
		return &view.CodeErrorStatusCode{
			StatusCode: 403,
			Response:   codeError,
		}
	default:
		logger.Error("REQUEST_ERROR", "err", err.Error())
		return &view.CodeErrorStatusCode{
			StatusCode: 500,
			Response:   codeError,
		}
	}
}

func (handler) HandleBearerAuth(ctx context.Context, _ string, t view.BearerAuth) (context.Context, error) {
	payload, err := auth.ViewVerifyToken(t.Token)
	if err != nil {
		return ctx, err
	}
	return viewutil.SetUserID(ctx, payload.UserID), nil
}
