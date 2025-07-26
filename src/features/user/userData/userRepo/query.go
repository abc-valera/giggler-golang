package userRepo

import (
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"giggler-golang/src/features/user/userData"
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

func (q query) GetByID(ctx context.Context, id uuid.UUID) (*userData.User, error) {
	model, err := q.dbQuery.WithContext(ctx).User.Where(q.dbQuery.User.ID.Eq(id)).First()
	return model, dbDto.QueryError(err)
}

func (q query) GetByEmail(ctx context.Context, email string) (*userData.User, error) {
	model, err := q.dbQuery.WithContext(ctx).User.Where(q.dbQuery.User.Email.Eq(email)).First()
	return model, dbDto.QueryError(err)
}

func (q query) GetAll(ctx context.Context, s dataModel.Selector) ([]*userData.User, error) {
	models, err := q.dbQuery.WithContext(ctx).User.Limit(int(s.Limit)).Offset(int(s.Offset)).Find()
	return models, dbDto.QueryError(err)
}
