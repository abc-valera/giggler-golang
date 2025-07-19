package dataModel

import (
	"context"

	"gorm.io/gorm"

	"giggler-golang/src/shared/data/dbDto"
)

type (
	IGenericCommandCreate[Model any] interface {
		Create(context.Context, *Model) error
	}

	IGenericCommandUpdate[Model any] interface {
		Update(context.Context, *Model) error
	}

	IGenericCommandDelete[Model any] interface {
		Delete(context.Context, *Model) error
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

type genericCommand[GormModel any] struct {
	db *gorm.DB
}

func NewGenericCommand[GormModel any](db *gorm.DB) genericCommand[GormModel] {
	return genericCommand[GormModel]{
		db: db,
	}
}

func (command genericCommand[GormModel]) Create(ctx context.Context, model *GormModel) error {
	return dbDto.CommandError(command.db.WithContext(ctx).Create(model))
}

func (command genericCommand[GormModel]) Update(ctx context.Context, model *GormModel) error {
	return dbDto.CommandError(command.db.WithContext(ctx).Save(model))
}

func (command genericCommand[GormModel]) Delete(ctx context.Context, model *GormModel) error {
	return dbDto.CommandError(command.db.WithContext(ctx).Delete(model))
}
