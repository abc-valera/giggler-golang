package userView

import (
	"context"

	"giggler-golang/src/features/user/userData/userRepo"
	"giggler-golang/src/shared/contexts"
	"giggler-golang/src/shared/data"
	"giggler-golang/src/shared/otel"
)

type userGetInput struct{}

type userGetOutput struct{}

func userGetHandler(ctx context.Context, input *userGetInput) (*userGetOutput, error) {
	ctx, span := otel.Trace(ctx)
	defer span.End()

	userID, err := contexts.GetUserID(ctx)
	if err != nil {
		return nil, err
	}

	user, err := userRepo.NewQuery(data.DB()).GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	return NewUserModel(user), nil
}
