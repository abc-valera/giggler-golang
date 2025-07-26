package userView

import (
	"context"

	"giggler-golang/src/features/user/internal"
	"giggler-golang/src/shared/otel"
)

type RegisterIn struct {
	Body struct {
		Username string `example:"valeriy"`
		Email    string `example:"valeriy@example.com"`
		Password string `example:"QWERTY1235*"`
	}
}

type RegisterOut struct{}

func RegisterHandler(ctx context.Context, input *RegisterIn) (*RegisterOut, error) {
	ctx, span := otel.Trace(ctx)
	defer span.End()

	if err := internal.Register(ctx, internal.RegisterIn{
		Username: input.Body.Username,
		Email:    input.Body.Email,
		Password: input.Body.Password,
	}); err != nil {
		return nil, err
	}
	return &RegisterOut{}, nil
}

type LoginIn struct {
	Body struct {
		Email    string `example:"valeriy@example.com"`
		Password string `example:"QWERTY1235*"`
	}
}

type LoginOut struct {
	Body struct {
		User         *UserModel
		AccessToken  string `example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9"`
		RefreshToken string `example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9"`
	}
}

func LoginHandler(ctx context.Context, input *LoginIn) (*LoginOut, error) {
	ctx, span := otel.Trace(ctx)
	defer span.End()

	loginResp, err := internal.Login(ctx, internal.LoginIn{
		Email:    input.Body.Email,
		Password: input.Body.Password,
	})
	if err != nil {
		return nil, err
	}

	var out LoginOut
	out.Body.User = NewUserDTO(loginResp.User)
	out.Body.AccessToken = loginResp.AccessToken
	out.Body.RefreshToken = loginResp.RefreshToken
	return &out, nil
}

type RefreshIn struct {
	Body struct {
		RefreshToken string `example:"some-refresh-token"`
	}
}

type RefreshOut struct {
	Body struct {
		AccessToken string `example:"some-access-token"`
	}
}

func RefreshHandler(ctx context.Context, input *RefreshIn) (*RefreshOut, error) {
	ctx, span := otel.Trace(ctx)
	defer span.End()

	accessToken, err := internal.Refresh(ctx, input.Body.RefreshToken)
	if err != nil {
		return nil, err
	}

	var resp RefreshOut
	resp.Body.AccessToken = accessToken
	return &resp, nil
}

type ProfileIn struct{}

type ProfileOut struct {
	Body struct {
		User *UserModel
	}
}

func ProfileHandler(ctx context.Context, _ *ProfileIn) (*ProfileOut, error) {
	ctx, span := otel.Trace(ctx)
	defer span.End()

	user, err := internal.Profile(ctx)
	if err != nil {
		return nil, err
	}

	var out ProfileOut
	out.Body.User = NewUserDTO(user)
	return &out, nil
}
