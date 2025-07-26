package internal

import (
	"context"

	"giggler-golang/src/features/user/JWT"
	"giggler-golang/src/features/user/userData"
	"giggler-golang/src/features/user/userPassword"
	"giggler-golang/src/shared/data"
	"giggler-golang/src/shared/data/gormgenQuery"
	"giggler-golang/src/shared/otel"
)

type LoginIn struct {
	Email    string
	Password string
}

type LoginOut struct {
	User         *userData.User
	AccessToken  string
	RefreshToken string
}

func Login(ctx context.Context, req LoginIn) (LoginOut, error) {
	ctx, span := otel.Trace(ctx)
	defer span.End()

	q := gormgenQuery.Use(data.DB())
	foundUser, err := q.WithContext(ctx).User.Where(q.User.Email.Eq(req.Email)).First()
	if err != nil {
		return LoginOut{}, err
	}

	if !userPassword.IsReal(ctx, req.Password, foundUser.HashedPassword) {
		return LoginOut{}, userPassword.ErrInvalidPass
	}

	payload := JWT.Payload{
		UserID: foundUser.ID,
	}

	return LoginOut{
		User:         foundUser,
		AccessToken:  JWT.NewAccessToken(payload),
		RefreshToken: JWT.NewRefreshToken(payload),
	}, nil
}
