package internal

import (
	"context"

	"giggler-golang/src/features/user/JWT"
	"giggler-golang/src/shared/otel"
)

func Refresh(ctx context.Context, refreshToken string) (string, error) {
	_, span := otel.Trace(ctx)
	defer span.End()

	payload, err := JWT.VerifyRefreshToken(refreshToken)
	if err != nil {
		return "", err
	}

	return JWT.NewAccessToken(JWT.Payload{
		UserID: payload.UserID,
	}), nil
}
