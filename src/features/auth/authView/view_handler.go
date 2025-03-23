package authView

import (
	"context"

	"github.com/abc-valera/giggler-golang/src/features/auth/internal/authLogic"
	"github.com/abc-valera/giggler-golang/src/features/user/userModel"
	"github.com/abc-valera/giggler-golang/src/shared/otel"
	"github.com/abc-valera/giggler-golang/src/shared/view/viewgen"
	"go.opentelemetry.io/otel/trace"
)

type View struct{}

func (View) AuthRegisterPost(ctx context.Context, req *viewgen.AuthRegisterPostReq) error {
	var span trace.Span
	ctx, span = otel.Trace(ctx)
	defer span.End()

	_, err := authLogic.AuthRegister(ctx, authLogic.AuthRegisterReq{
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password,
	})
	return err
}

func (View) AuthLoginPost(ctx context.Context, req *viewgen.AuthLoginPostReq) (*viewgen.AuthLoginPostOK, error) {
	var span trace.Span
	ctx, span = otel.Trace(ctx)
	defer span.End()

	loggedUser, err := authLogic.AuthLogin(ctx, authLogic.AuthLoginReq{
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		return nil, err
	}

	access, err := viewCreateAccessToken(loggedUser.ID)
	if err != nil {
		return nil, err
	}

	refresh, err := viewCreateRefreshToken(loggedUser.ID)
	if err != nil {
		return nil, err
	}

	return &viewgen.AuthLoginPostOK{
		UserResponse: *userModel.NewViewDTO(loggedUser),
		AccessToken:  access,
		RefreshToken: refresh,
	}, nil
}

func (View) AuthRefreshPost(ctx context.Context, req *viewgen.AuthRefreshPostReq) (*viewgen.AuthRefreshPostOK, error) {
	_, span := otel.Trace(ctx)
	defer span.End()

	payload, err := ViewVerifyToken(req.RefreshToken)
	if err != nil {
		return nil, err
	}

	if !payload.IsRefresh {
		return nil, errProvidedAccessToken
	}

	access, err := viewCreateAccessToken(payload.UserID)
	if err != nil {
		return nil, err
	}

	return &viewgen.AuthRefreshPostOK{
		AccessToken: access,
	}, nil
}
