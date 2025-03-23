package internal

import (
	"context"

	"go.opentelemetry.io/otel/trace"

	"giggler-golang/src/features/auth/authJWT"
	"giggler-golang/src/features/auth/authPassword"
	"giggler-golang/src/features/user/userData"
	"giggler-golang/src/shared/data"
	"giggler-golang/src/shared/data/dbgen/gormModel"
	"giggler-golang/src/shared/otel"
)

type LoginReq struct {
	Email    string
	Password string
}

type LoginResp struct {
	User         *gormModel.User
	AccessToken  string
	RefreshToken string
}

func Login(ctx context.Context, req LoginReq) (LoginResp, error) {
	var span trace.Span
	ctx, span = otel.Trace(ctx)
	defer span.End()

	foundUser, err := userData.NewQuery(data.DB()).GetByEmail(ctx, req.Email)
	if err != nil {
		return LoginResp{}, err
	}

	if !authPassword.IsReal(ctx, req.Password, foundUser.HashedPassword) {
		return LoginResp{}, authPassword.ErrInvalidPass
	}

	access, err := authJWT.CreateToken(authJWT.Payload{
		UserID:    foundUser.ID,
		IsRefresh: false,
	})
	if err != nil {
		return LoginResp{}, err
	}

	refresh, err := authJWT.CreateToken(authJWT.Payload{
		UserID:    foundUser.ID,
		IsRefresh: true,
	})
	if err != nil {
		return LoginResp{}, err
	}

	return LoginResp{
		User:         foundUser,
		AccessToken:  access,
		RefreshToken: refresh,
	}, nil
}
