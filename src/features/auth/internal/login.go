package internal

import (
	"context"

	"giggler-golang/src/features/auth/authJWT"
	"giggler-golang/src/features/user/userData"
	"giggler-golang/src/features/user/userData/userRepo"
	"giggler-golang/src/features/user/userPassword"
	"giggler-golang/src/shared/data"
	"giggler-golang/src/shared/otel"
)

type LoginReq struct {
	Email    string
	Password string
}

type LoginResp struct {
	User         *userData.User
	AccessToken  string
	RefreshToken string
}

func Login(ctx context.Context, req LoginReq) (LoginResp, error) {
	ctx, span := otel.Trace(ctx)
	defer span.End()

	foundUser, err := userRepo.NewQuery(data.DB()).GetByEmail(ctx, req.Email)
	if err != nil {
		return LoginResp{}, err
	}

	if !userPassword.IsReal(ctx, req.Password, foundUser.HashedPassword) {
		return LoginResp{}, userPassword.ErrInvalidPass
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
