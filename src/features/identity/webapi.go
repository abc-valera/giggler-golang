package identity

import (
	"context"
	"net/http"
	"time"

	"github.com/danielgtaylor/huma/v2"
	"github.com/google/uuid"
)

func InitRoutes(api huma.API, usecase Usecase) {
	h := handler{usecase: usecase}

	huma.Post(api, "/register", h.registerHandler)
	huma.Post(api, "/login", h.loginHandler)
	huma.Post(api, "/refresh", h.refreshHandler)
	huma.Register(api, huma.Operation{
		OperationID: "profile",
		Method:      http.MethodGet,
		Path:        "/profile",
		Middlewares: huma.Middlewares{h.AuthMiddleware},
	}, h.profileHandler)
}

type WebapiDTO struct {
	ID        uuid.UUID
	Username  string
	Email     string
	Fullname  *string
	Status    *string
	CreatedAt time.Time
}

func NewModel(user User) *WebapiDTO {
	return &WebapiDTO{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		Fullname:  user.Fullname,
		Status:    user.Status,
		CreatedAt: user.CreatedAt,
	}
}

type handler struct {
	usecase Usecase
}

type registerInput struct {
	Body struct {
		Username string `example:"valeriy"`
		Email    string `example:"valeriy@example.com"`
		Password string `example:"QWERTY1235*"`
	}
}

type registerOutput struct{}

func (h handler) registerHandler(ctx context.Context, input *registerInput) (*registerOutput, error) {
	if err := h.usecase.Register(ctx, RegisterInput{
		Username: input.Body.Username,
		Email:    input.Body.Email,
		Password: input.Body.Password,
	}); err != nil {
		return nil, err
	}
	return &registerOutput{}, nil
}

type loginInput struct {
	Body struct {
		Email    string `example:"valeriy@example.com"`
		Password string `example:"QWERTY1235*"`
	}
}

type loginOutput struct {
	Body struct {
		User         *WebapiDTO
		AccessToken  string `example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9"`
		RefreshToken string `example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9"`
	}
}

func (h handler) loginHandler(ctx context.Context, input *loginInput) (*loginOutput, error) {
	loginResp, err := h.usecase.Login(ctx, LoginInput{
		Email:    input.Body.Email,
		Password: input.Body.Password,
	})
	if err != nil {
		return nil, err
	}

	var out loginOutput
	out.Body.User = NewModel(loginResp.User)
	out.Body.AccessToken = loginResp.AccessTokenJWT
	out.Body.RefreshToken = loginResp.RefreshTokenJWT
	return &out, nil
}

type refreshInput struct {
	Body struct {
		RefreshToken string `example:"some-refresh-token"`
	}
}

type refreshOutput struct {
	Body struct {
		AccessToken string `example:"some-access-token"`
	}
}

func (h handler) refreshHandler(ctx context.Context, input *refreshInput) (*refreshOutput, error) {
	accessToken, err := h.usecase.Refresh(ctx, input.Body.RefreshToken)
	if err != nil {
		return nil, err
	}

	var resp refreshOutput
	resp.Body.AccessToken = accessToken
	return &resp, nil
}

type profileInput struct{}

type profileOutput struct {
	Body struct {
		User *WebapiDTO
	}
}

func (h handler) profileHandler(ctx context.Context, _ *profileInput) (*profileOutput, error) {
	user, err := h.usecase.Profile(ctx)
	if err != nil {
		return nil, err
	}

	var out profileOutput
	out.Body.User = NewModel(user)
	return &out, nil
}

var HttpAuthHeaderName = "Authorization"

func (h handler) AuthMiddleware(ctx huma.Context, next func(huma.Context)) {
	payload, err := h.usecase.authTokenService.VerifyAccessTokenJWT(ctx.Header(HttpAuthHeaderName))
	if err != nil {
		ctx.SetStatus(http.StatusUnauthorized)
		ctx.BodyWriter().Write([]byte(err.Error()))
		return
	}

	ctx = huma.WithContext(ctx, SetAuthToken(ctx.Context(), payload))

	next(ctx)
}
