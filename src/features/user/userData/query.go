package userData

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

func (q query) GetByID(ctx context.Context, id string) (*gormModel.User, error) {
	model, err := q.dbQuery.WithContext(ctx).User.Where(q.dbQuery.User.ID.Eq(id)).First()
	return model, dbDTO.QueryError(err)
}

func (q query) GetByEmail(ctx context.Context, email string) (*gormModel.User, error) {
	model, err := q.dbQuery.WithContext(ctx).User.Where(q.dbQuery.User.Email.Eq(email)).First()
	return model, dbDTO.QueryError(err)
}

func (q query) GetAll(ctx context.Context, s dataModel.Selector) ([]*gormModel.User, error) {
	models, err := q.dbQuery.WithContext(ctx).User.Limit(int(s.Limit)).Offset(int(s.Offset)).Find()
	return models, dbDTO.QueryError(err)
}
