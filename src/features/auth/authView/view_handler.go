package authView

import (
	"context"

	"giggler-golang/src/features/auth/internal"
	"giggler-golang/src/features/user/userDto"
	"giggler-golang/src/shared/otel"
	"giggler-golang/src/shared/view/viewgen"

	"go.opentelemetry.io/otel/trace"
)

type Handler struct{}

func (Handler) AuthRegisterPost(ctx context.Context, req *viewgen.AuthRegisterPostReq) error {
	var span trace.Span
	ctx, span = otel.Trace(ctx)
	defer span.End()

	_, err := internal.AuthRegister(ctx, internal.AuthRegisterReq{
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password,
	})
	return err
}

func (Handler) AuthLoginPost(ctx context.Context, req *viewgen.AuthLoginPostReq) (*viewgen.AuthLoginPostOK, error) {
	var span trace.Span
	ctx, span = otel.Trace(ctx)
	defer span.End()

	loggedUser, err := internal.AuthLogin(ctx, internal.AuthLoginReq{
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
		UserResponse: *userDto.NewViewDTO(loggedUser),
		AccessToken:  access,
		RefreshToken: refresh,
	}, nil
}

func (Handler) AuthRefreshPost(ctx context.Context, req *viewgen.AuthRefreshPostReq) (*viewgen.AuthRefreshPostOK, error) {
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
