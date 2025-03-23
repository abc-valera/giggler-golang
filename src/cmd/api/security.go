package main

import (
	"context"

	"giggler-golang/src/features/auth/authView"
	"giggler-golang/src/shared/contexts"
	"giggler-golang/src/shared/view/viewgen"
)

type securityHandler struct{}

func (securityHandler) HandleBearerAuth(ctx context.Context, _ string, t viewgen.BearerAuth) (context.Context, error) {
	payload, err := authView.ViewVerifyToken(t.Token)
	if err != nil {
		return ctx, err
	}

	return contexts.SetUserID(ctx, payload.UserID), nil
}
