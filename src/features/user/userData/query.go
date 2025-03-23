package userData

import (
	"context"

	"giggler-golang/src/features/user/internal"
	"giggler-golang/src/features/user/userModel"
	"giggler-golang/src/shared/data"
	"giggler-golang/src/shared/data/gormgen/gormQuery"
	"giggler-golang/src/shared/env"
)

var Query = initQuery()

type iQuery interface {
	GetByID(ctx context.Context, id string) (userModel.User, error)
	GetByEmail(ctx context.Context, email string) (userModel.User, error)
	GetAll(ctx context.Context, s data.Selector) ([]userModel.User, error)
}

func initQuery() func(dataStore data.IDS) iQuery {
	switch data.DbEnvVal {
	case data.DbVariantPostgres:
		return func(dataStore data.IDS) iQuery {
			return internal.GormQ{Query: gormQuery.Use(data.GormDS(dataStore))}
		}
	default:
		panic(env.ErrInvalidEnvValue)
	}
}
