package internal

import (
	"context"

	"giggler-golang/src/features/user/userDto"
	"giggler-golang/src/features/user/userModel"
	"giggler-golang/src/shared/data"
	"giggler-golang/src/shared/data/gormgen/gormQuery"
)

type GormQ struct{ *gormQuery.Query }

func (q GormQ) GetByID(ctx context.Context, id string) (userModel.User, error) {
	dto, err := q.WithContext(ctx).User.Where(q.User.ID.Eq(id)).First()
	return userDto.FromGormDTO(dto), data.GormQueryError(err)
}

func (q GormQ) GetByEmail(ctx context.Context, email string) (userModel.User, error) {
	dto, err := q.WithContext(ctx).User.Where(q.User.Email.Eq(email)).First()
	return userDto.FromGormDTO(dto), data.GormQueryError(err)
}

func (q GormQ) GetAll(ctx context.Context, s data.Selector) ([]userModel.User, error) {
	dtos, err := q.WithContext(ctx).User.Limit(int(s.PagingLimit)).Offset(int(s.PagingOffset)).Find()
	return userDto.FromGormDTOs(dtos), data.GormQueryError(err)
}
