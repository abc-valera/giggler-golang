package userData

import (
	"context"

	"github.com/abc-valera/giggler-golang/src/components/data"
	"github.com/abc-valera/giggler-golang/src/components/env"
	"github.com/abc-valera/giggler-golang/src/components/gormgen/gormQuery"
	"github.com/abc-valera/giggler-golang/src/features/user/userModel"
	"github.com/abc-valera/giggler-golang/src/shared/ds"
)

var Query = initQuery()

type iQuery interface {
	GetByID(ctx context.Context, id string) (userModel.User, error)
	GetByEmail(ctx context.Context, email string) (userModel.User, error)
	GetAll(ctx context.Context, s ds.Selector) ([]userModel.User, error)
}

func initQuery() func(dataStore ds.IDS) iQuery {
	switch ds.DbVal {
	case ds.DbVariantGorm:
		return func(dataStore ds.IDS) iQuery {
			return &gormQ{gormQuery.Use(ds.GormDS(dataStore))}
		}
	default:
		panic(env.ErrInvalidEnvValue)
	}
}

type gormQ struct{ *gormQuery.Query }

func (q gormQ) GetByID(ctx context.Context, id string) (userModel.User, error) {
	dto, err := q.WithContext(ctx).User.Where(q.User.ID.Eq(id)).First()
	return userModel.FromGormDTO(dto), data.GormQueryError(err)
}

func (q gormQ) GetByEmail(ctx context.Context, email string) (userModel.User, error) {
	dto, err := q.WithContext(ctx).User.Where(q.User.Email.Eq(email)).First()
	return userModel.FromGormDTO(dto), data.GormQueryError(err)
}

func (q gormQ) GetAll(ctx context.Context, s ds.Selector) ([]userModel.User, error) {
	dtos, err := q.WithContext(ctx).User.Limit(int(s.PagingLimit)).Offset(int(s.PagingOffset)).Find()
	return userModel.FromGormDTOs(dtos), data.GormQueryError(err)
}
