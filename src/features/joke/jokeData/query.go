package jokeData

import (
	"context"

	"github.com/abc-valera/giggler-golang/src/features/joke/jokeModel"
	"github.com/abc-valera/giggler-golang/src/shared/data"
	"github.com/abc-valera/giggler-golang/src/shared/data/gormgen/gormQuery"
	"github.com/abc-valera/giggler-golang/src/shared/env"
)

var Query = initQuery()

type iQuery interface {
	GetByID(ctx context.Context, id string) (jokeModel.Joke, error)
	GetAllByUserID(ctx context.Context, userID string, s data.Selector) ([]jokeModel.Joke, error)
}

func initQuery() func(data.IDS) iQuery {
	switch data.DbVal {
	case data.DbVariantGorm:
		return func(dataStore data.IDS) iQuery {
			return &gormQ{gormQuery.Use(data.GormDS(dataStore))}
		}
	default:
		panic(env.ErrInvalidEnvValue)
	}
}

type gormQ struct{ *gormQuery.Query }

func (q gormQ) GetByID(ctx context.Context, id string) (jokeModel.Joke, error) {
	dto, err := q.WithContext(ctx).Joke.Where(q.Joke.ID.Eq(id)).First()
	return jokeModel.FromGormDTO(dto), data.GormQueryError(err)
}

func (q gormQ) GetAllByUserID(ctx context.Context, userID string, s data.Selector) ([]jokeModel.Joke, error) {
	dtos, err := q.WithContext(ctx).Joke.Where(q.Joke.UserID.Eq(userID)).Find()
	return jokeModel.FromGormDTOs(dtos), data.GormQueryError(err)
}
