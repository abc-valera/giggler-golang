package userViewRestapi

import (
	"net/http"

	"github.com/danielgtaylor/huma/v2"

	"giggler-golang/src/features/user/JWT"
)

func AuthMiddleware(ctx huma.Context, next func(huma.Context)) {
	authHeader := ctx.Header("Authorization")

	payload, err := JWT.VerifyAccessToken(authHeader)
	if err != nil {
		ctx.SetStatus(http.StatusUnauthorized)
		ctx.BodyWriter().Write([]byte(err.Error()))
		return
	}

	ctx = huma.WithContext(ctx, JWT.Set(ctx.Context(), payload))

	next(ctx)
}
