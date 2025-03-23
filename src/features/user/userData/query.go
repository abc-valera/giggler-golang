package userData

import (
	"context"

	"giggler-golang/src/features/user/userDTO"
	"giggler-golang/src/features/user/userModel"
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

func (q query) GetByID(ctx context.Context, id string) (userModel.User, error) {
	dto, err := q.dbQuery.WithContext(ctx).User.Where(q.dbQuery.User.ID.Eq(id)).First()
	return userDTO.FromGorm(dto), dbDTO.QueryError(err)
}

func (q query) GetByEmail(ctx context.Context, email string) (userModel.User, error) {
	dto, err := q.dbQuery.WithContext(ctx).User.Where(q.dbQuery.User.Email.Eq(email)).First()
	return userDTO.FromGorm(dto), dbDTO.QueryError(err)
}

func (q query) GetAll(ctx context.Context, s dataModel.Selector) ([]userModel.User, error) {
	dtos, err := q.dbQuery.WithContext(ctx).User.Limit(int(s.Limit)).Offset(int(s.Offset)).Find()
	return userDTO.FromGorms(dtos), dbDTO.QueryError(err)
}
