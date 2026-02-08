package userViewWebapi

import (
	"context"
	"net/http"

	"github.com/danielgtaylor/huma/v2"

	"giggler-golang/src/features/user/userData"
	"giggler-golang/src/features/user/userUsecase"
)

func InitRoutes(api huma.API, usecase userUsecase.Usecase) {
	h := handler{usecase: usecase}

	huma.Post(api, "/register", h.registerHandler)
	huma.Post(api, "/login", h.loginHandler)
	huma.Post(api, "/refresh", h.refreshHandler)
	// TODO: add a helper function for such registration
	huma.Register(api, huma.Operation{
		OperationID: "profile",
		Method:      http.MethodGet,
		Path:        "/profile",
		Middlewares: huma.Middlewares{authMiddleware},
		Security: []map[string][]string{
			{string(userData.Basic): {}},
		},
	}, h.profileHandler)
}

type handler struct {
	usecase userUsecase.Usecase
}

type registerIn struct {
	Body struct {
		Username string `example:"valeriy"`
		Email    string `example:"valeriy@example.com"`
		Password string `example:"QWERTY1235*"`
	}
}

type registerOut struct{}

func (h handler) registerHandler(ctx context.Context, input *registerIn) (*registerOut, error) {
	if err := h.usecase.Register(ctx, userUsecase.RegisterIn{
		Username: input.Body.Username,
		Email:    input.Body.Email,
		Password: input.Body.Password,
	}); err != nil {
		return nil, err
	}
	return &registerOut{}, nil
}

type loginIn struct {
	Body struct {
		Email    string `example:"valeriy@example.com"`
		Password string `example:"QWERTY1235*"`
	}
}

type loginOut struct {
	Body struct {
		User         *Model
		AccessToken  string `example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9"`
		RefreshToken string `example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9"`
	}
}

func (h handler) loginHandler(ctx context.Context, input *loginIn) (*loginOut, error) {
	loginResp, err := h.usecase.Login(ctx, userUsecase.LoginIn{
		Email:    input.Body.Email,
		Password: input.Body.Password,
	})
	if err != nil {
		return nil, err
	}

	var out loginOut
	out.Body.User = NewModel(loginResp.User)
	out.Body.AccessToken = loginResp.AccessToken
	out.Body.RefreshToken = loginResp.RefreshToken
	return &out, nil
}

type refreshIn struct {
	Body struct {
		RefreshToken string `example:"some-refresh-token"`
	}
}

type refreshOut struct {
	Body struct {
		AccessToken string `example:"some-access-token"`
	}
}

func (h handler) refreshHandler(ctx context.Context, input *refreshIn) (*refreshOut, error) {
	accessToken, err := h.usecase.Refresh(ctx, input.Body.RefreshToken)
	if err != nil {
		return nil, err
	}

	var resp refreshOut
	resp.Body.AccessToken = accessToken
	return &resp, nil
}

type profileIn struct{}

type profileOut struct {
	Body struct {
		User *Model
	}
}

func (h handler) profileHandler(ctx context.Context, _ *profileIn) (*profileOut, error) {
	user, err := h.usecase.Profile(ctx)
	if err != nil {
		return nil, err
	}

	var out profileOut
	out.Body.User = NewModel(user)
	return &out, nil
}
