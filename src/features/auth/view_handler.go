package auth

import (
	"context"

	"github.com/abc-valera/giggler-golang/gen/view"
	"github.com/abc-valera/giggler-golang/src/features/user"
	"github.com/abc-valera/giggler-golang/src/shared/otel"
	"go.opentelemetry.io/otel/trace"
)

type View struct{}

func (View) AuthRegisterPost(ctx context.Context, req *view.AuthRegisterPostReq) error {
	var span trace.Span
	ctx, span = otel.Trace(ctx)
	defer span.End()

	_, err := authRegister(ctx, authRegisterRequest{
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password,
	})
	return err
}

func (View) AuthLoginPost(ctx context.Context, req *view.AuthLoginPostReq) (*view.AuthLoginPostOK, error) {
	var span trace.Span
	ctx, span = otel.Trace(ctx)
	defer span.End()

	loggedUser, err := authLogin(ctx, authLoginRequest{
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

	return &view.AuthLoginPostOK{
		UserResponse: *user.ViewNewUserDTO(loggedUser),
		AccessToken:  access,
		RefreshToken: refresh,
	}, nil
}

func (View) AuthRefreshPost(ctx context.Context, req *view.AuthRefreshPostReq) (*view.AuthRefreshPostOK, error) {
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

	return &view.AuthRefreshPostOK{
		AccessToken: access,
	}, nil
}
