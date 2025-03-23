package jokeData

import (
	"context"

	"giggler-golang/src/features/joke/internal"
	"giggler-golang/src/features/joke/jokeModel"
	"giggler-golang/src/shared/data"
	"giggler-golang/src/shared/data/gormgen/gormQuery"
	"giggler-golang/src/shared/env"
)

var Query = initQuery()

type iQuery interface {
	GetByID(ctx context.Context, id string) (jokeModel.Joke, error)
	GetAllByUserID(ctx context.Context, userID string, s data.Selector) ([]jokeModel.Joke, error)
}

func initQuery() func(data.IDS) iQuery {
	switch data.DbEnvVal {
	case data.DbVariantPostgres:
		return func(dataStore data.IDS) iQuery {
			return &internal.GormQ{Query: gormQuery.Use(data.GormDS(dataStore))}
		}
	default:
		panic(env.ErrInvalidEnvValue)
	}
}
