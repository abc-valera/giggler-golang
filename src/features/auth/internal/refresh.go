package internal

import (
	"context"

	"giggler-golang/src/features/auth/authJWT"
	"giggler-golang/src/shared/otel"
)

func Refresh(ctx context.Context, refreshToken string) (string, error) {
	_, span := otel.Trace(ctx)
	defer span.End()

	payload, err := authJWT.VerifyToken(refreshToken)
	if err != nil {
		return "", err
	}

	if !payload.IsRefresh {
		return "", authJWT.ErrProvidedAccessToken
	}

	access, err := authJWT.CreateToken(authJWT.Payload{
		UserID:    payload.UserID,
		IsRefresh: false,
	})
	if err != nil {
		return "", err
	}

	return access, nil
}
