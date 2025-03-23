package authView

import (
	"context"

	"giggler-golang/src/features/auth/internal"
	"giggler-golang/src/features/user/userView"
	"giggler-golang/src/shared/otel"
	"giggler-golang/src/shared/view/viewgen"
)

type Handler struct{}

func (Handler) AuthRegisterPost(ctx context.Context, req *viewgen.AuthRegisterPostReq) error {
	ctx, span := otel.Trace(ctx)
	defer span.End()

	_, err := internal.Register(ctx, internal.RegisterReq{
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password,
	})

	return err
}

func (Handler) AuthLoginPost(ctx context.Context, req *viewgen.AuthLoginPostReq) (*viewgen.AuthLoginPostOK, error) {
	ctx, span := otel.Trace(ctx)
	defer span.End()

	authLoginResp, err := internal.Login(ctx, internal.LoginReq{
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		return nil, err
	}

	return &viewgen.AuthLoginPostOK{
		UserResponse: *userView.NewUserModel(authLoginResp.User),
		AccessToken:  authLoginResp.AccessToken,
		RefreshToken: authLoginResp.RefreshToken,
	}, nil
}

func (Handler) AuthRefreshPost(ctx context.Context, req *viewgen.AuthRefreshPostReq) (*viewgen.AuthRefreshPostOK, error) {
	ctx, span := otel.Trace(ctx)
	defer span.End()

	accessToken, err := internal.Refresh(ctx, req.RefreshToken)
	if err != nil {
		return nil, err
	}

	return &viewgen.AuthRefreshPostOK{
		AccessToken: accessToken,
	}, nil
}
