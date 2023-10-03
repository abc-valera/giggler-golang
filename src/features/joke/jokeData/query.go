package jokeData

import (
	"context"

	"github.com/abc-valera/giggler-golang/src/components/env"
	"github.com/abc-valera/giggler-golang/src/components/gormgen/gormQuery"
	"github.com/abc-valera/giggler-golang/src/features/joke/jokeModel"
	"github.com/abc-valera/giggler-golang/src/shared/data"
)

var Query func(data.IDS) iQuery

func init() {
	switch data.DbVal {
	case data.DbVariantGorm:
		Query = func(dataStore data.IDS) iQuery {
			return &gormQ{gormQuery.Use(data.GormDS(dataStore))}
		}
	default:
		panic(env.ErrInvalidEnvValue)
	}
}

type iQuery interface {
	GetByID(ctx context.Context, id string) (jokeModel.Joke, error)
	GetAllByUserID(ctx context.Context, userID string, s data.Selector) ([]jokeModel.Joke, error)
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
