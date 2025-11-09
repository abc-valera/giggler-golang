package internal

import (
	"context"

	"giggler-golang/src/features/user/JWT"
	"giggler-golang/src/features/user/userData"
	"giggler-golang/src/shared/data"
	"giggler-golang/src/shared/data/gormgenQuery"
	"giggler-golang/src/shared/otel"
)

func Profile(ctx context.Context) (*userData.User, error) {
	ctx, span := otel.Trace(ctx)
	defer span.End()

	p := JWT.Get(ctx)

	q := gormgenQuery.Use(data.GetDb())
	user, err := q.WithContext(ctx).User.Where(q.User.ID.Eq(p.UserID)).First()
	if err != nil {
		return nil, err
	}

	return user, nil
}
