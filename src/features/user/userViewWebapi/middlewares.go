package userViewWebapi

import (
	"net/http"

	"github.com/danielgtaylor/huma/v2"

	"giggler-golang/src/features/user/jwt"
)

func authMiddleware(ctx huma.Context, next func(huma.Context)) {
	authHeader := ctx.Header("Authorization")

	payload, err := jwt.VerifyAccessToken(authHeader)
	if err != nil {
		ctx.SetStatus(http.StatusUnauthorized)
		ctx.BodyWriter().Write([]byte(err.Error()))
		return
	}

	ctx = huma.WithContext(ctx, jwt.Set(ctx.Context(), payload))

	next(ctx)
}
