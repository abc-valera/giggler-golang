package authView

import (
	"context"

	"go.opentelemetry.io/otel/trace"

	"giggler-golang/src/features/auth/authJWT"
	"giggler-golang/src/features/auth/internal"
	"giggler-golang/src/features/user/userDTO"
	"giggler-golang/src/shared/otel"
	"giggler-golang/src/shared/view/viewgen"
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

	access, err := authJWT.CreateToken(authJWT.Payload{
		UserID:    loggedUser.ID,
		IsRefresh: false,
	})
	if err != nil {
		return nil, err
	}

	refresh, err := authJWT.CreateToken(authJWT.Payload{
		UserID:    loggedUser.ID,
		IsRefresh: true,
	})
	if err != nil {
		return nil, err
	}

	return &viewgen.AuthLoginPostOK{
		UserResponse: *userDTO.NewView(loggedUser),
		AccessToken:  access,
		RefreshToken: refresh,
	}, nil
}

func (Handler) AuthRefreshPost(ctx context.Context, req *viewgen.AuthRefreshPostReq) (*viewgen.AuthRefreshPostOK, error) {
	_, span := otel.Trace(ctx)
	defer span.End()

	payload, err := authJWT.VerifyToken(req.RefreshToken)
	if err != nil {
		return nil, err
	}

	if !payload.IsRefresh {
		return nil, authJWT.ErrProvidedAccessToken
	}

	access, err := authJWT.CreateToken(authJWT.Payload{
		UserID:    payload.UserID,
		IsRefresh: false,
	})
	if err != nil {
		return nil, err
	}

	return &viewgen.AuthRefreshPostOK{
		AccessToken: access,
	}, nil
}
