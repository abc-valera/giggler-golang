package data

import (
	"context"

	"gorm.io/gorm"
)

type (
	IGenericCommandCreate[Model any] interface {
		Create(context.Context, Model) error
	}

	IGenericCommandUpdate[Model any] interface {
		Update(context.Context, Model) error
	}

	IGenericCommandDelete[Model any] interface {
		Delete(context.Context, Model) error
	}

	IGenericCommandCreateDelete[Model any] interface {
		IGenericCommandCreate[Model]
		IGenericCommandDelete[Model]
	}

	IGenericCommandUpdateDelete[Model any] interface {
		IGenericCommandUpdate[Model]
		IGenericCommandDelete[Model]
	}

	IGenericCommandCreateUpdateDelete[Model any] interface {
		IGenericCommandCreate[Model]
		IGenericCommandUpdate[Model]
		IGenericCommandDelete[Model]
	}
)

type gormGenericCommand[DomainModel, GormModel any] struct {
	db  *gorm.DB
	dto func(DomainModel) GormModel
}

func NewGormGenericCommand[DomainModel, GormModel any](
	db *gorm.DB,
	newDto func(DomainModel) GormModel,
) IGenericCommandCreateUpdateDelete[DomainModel] {
	return gormGenericCommand[DomainModel, GormModel]{
		db:  db,
		dto: newDto,
	}
}

func (command gormGenericCommand[DomainModel, GormModel]) Create(ctx context.Context, req DomainModel) error {
	gormModel := command.dto(req)
	res := command.db.WithContext(ctx).Create(&gormModel)
	return GormCommandError(res)
}

func (command gormGenericCommand[DomainModel, GormModel]) Update(ctx context.Context, req DomainModel) error {
	gormModel := command.dto(req)
	res := command.db.WithContext(ctx).Save(&gormModel)
	return GormCommandError(res)
}

func (command gormGenericCommand[DomainModel, GormModel]) Delete(ctx context.Context, req DomainModel) error {
	gormModel := command.dto(req)
	res := command.db.WithContext(ctx).Delete(&gormModel)
	return GormCommandError(res)
}
