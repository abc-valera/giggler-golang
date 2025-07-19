package jokeRepo

import (
	"context"

	"gorm.io/gorm"

	"giggler-golang/src/features/joke/jokeData"
	"giggler-golang/src/shared/data/dataModel"
	"giggler-golang/src/shared/data/dbDto"
	"giggler-golang/src/shared/data/gormgenQuery"
)

type query struct {
	dbQuery *gormgenQuery.Query
}

func NewQuery(db *gorm.DB) query {
	return query{
		dbQuery: gormgenQuery.Use(db),
	}
}

func (q query) GetByID(ctx context.Context, id string) (*jokeData.Joke, error) {
	model, err := q.dbQuery.WithContext(ctx).Joke.Where(q.dbQuery.Joke.ID.Eq(id)).First()
	return model, dbDto.QueryError(err)
}

func (q query) GetAllByUserID(ctx context.Context, userID string, s dataModel.Selector) ([]*jokeData.Joke, error) {
	models, err := q.dbQuery.WithContext(ctx).Joke.Where(q.dbQuery.Joke.UserID.Eq(userID)).Find()
	return models, dbDto.QueryError(err)
}
