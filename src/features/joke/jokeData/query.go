package jokeData

import (
	"context"

	"giggler-golang/src/features/joke/jokeDTO"
	"giggler-golang/src/features/joke/jokeModel"
	"giggler-golang/src/shared/data"
	"giggler-golang/src/shared/data/dataModel"
	"giggler-golang/src/shared/data/dbDTO"
	"giggler-golang/src/shared/data/dbgen/gormQuery"
)

type query struct {
	dbQuery *gormQuery.Query
}

func NewQuery(ds data.IDB) query {
	return query{
		dbQuery: gormQuery.Use(getDB(ds)),
	}
}

func (q query) GetByID(ctx context.Context, id string) (jokeModel.Joke, error) {
	dto, err := q.dbQuery.WithContext(ctx).Joke.Where(q.dbQuery.Joke.ID.Eq(id)).First()
	return jokeDTO.FromGorm(dto), dbDTO.QueryError(err)
}

func (q query) GetAllByUserID(ctx context.Context, userID string, s dataModel.Selector) ([]jokeModel.Joke, error) {
	dtos, err := q.dbQuery.WithContext(ctx).Joke.Where(q.dbQuery.Joke.UserID.Eq(userID)).Find()
	return jokeDTO.FromGorms(dtos), dbDTO.QueryError(err)
}
