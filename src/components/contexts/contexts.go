package contexts

import (
	"context"

	"github.com/abc-valera/giggler-golang/src/components/errutil"
)

const userIDKey key = "user-id"

type key string

func SetUserID(ctx context.Context, userID string) context.Context {
	return context.WithValue(ctx, userIDKey, userID)
}

func GetUserID(ctx context.Context) (string, error) {
	userID, ok := ctx.Value(userIDKey).(string)
	if !ok {
		return "", errutil.NewInternalString("UserID not found in context")
	}
	return userID, nil
}
