package jokeData

import (
	"context"

	"gorm.io/gorm"

	"giggler-golang/src/shared/data/dataModel"
	"giggler-golang/src/shared/data/dbDTO"
	"giggler-golang/src/shared/data/dbgen/gormModel"
	"giggler-golang/src/shared/data/dbgen/gormQuery"
)

type query struct {
	dbQuery *gormQuery.Query
}

func NewQuery(db *gorm.DB) query {
	return query{
		dbQuery: gormQuery.Use(db),
	}
}

func (q query) GetByID(ctx context.Context, id string) (*gormModel.Joke, error) {
	model, err := q.dbQuery.WithContext(ctx).Joke.Where(q.dbQuery.Joke.ID.Eq(id)).First()
	return model, dbDTO.QueryError(err)
}

func (q query) GetAllByUserID(ctx context.Context, userID string, s dataModel.Selector) ([]*gormModel.Joke, error) {
	models, err := q.dbQuery.WithContext(ctx).Joke.Where(q.dbQuery.Joke.UserID.Eq(userID)).Find()
	return models, dbDTO.QueryError(err)
}
