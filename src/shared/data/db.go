package data

import (
	"context"

	"gorm.io/gorm"

	"giggler-golang/src/shared/data/dbDTO"
	"giggler-golang/src/shared/errutil"
)

type IDB interface {
	WithinTX(ctx context.Context, txFunc func(ctx context.Context, ds IDB) error) error
	NewTX() (ITX, error)

	getDbInstance() *gorm.DB
}

type ITX interface {
	WithinTX(ctx context.Context, txFunc func(ctx context.Context, ds IDB) error) error
	NewTX() (ITX, error)
	Commit() error
	Rollback() error

	getDbInstance() *gorm.DB
}

type db struct {
	gormDB *gorm.DB
}

func (d db) getDbInstance() *gorm.DB {
	return d.gormDB
}

func (d db) NewTX() (ITX, error) {
	gormTX := d.gormDB.Begin()
	if gormTX.Error != nil {
		return nil, errutil.NewInternal(gormTX.Error)
	}
	return &tx{gormTX: gormTX}, nil
}

func (d db) WithinTX(ctx context.Context, txFunc func(ctx context.Context, tx IDB) error) error {
	newTX, err := d.NewTX()
	if err != nil {
		return err
	}

	defer func() {
		if r := recover(); r != nil {
			newTX.Rollback()
			panic(r)
		}
	}()

	if err := txFunc(ctx, newTX); err != nil {
		newTX.Rollback()
		return err
	}

	return newTX.Commit()
}

type tx struct {
	gormTX *gorm.DB
}

func (t tx) NewTX() (ITX, error) {
	gormTX := t.gormTX.Begin()
	if gormTX.Error != nil {
		return nil, errutil.NewInternal(gormTX.Error)
	}
	return &tx{gormTX: gormTX}, nil
}

func (t tx) Commit() error {
	return errutil.NewInternal(t.gormTX.Commit().Error)
}

func (t tx) Rollback() error {
	return errutil.NewInternal(t.gormTX.Rollback().Error)
}

func (t tx) WithinTX(ctx context.Context, txFunc func(ctx context.Context, tx IDB) error) error {
	newTX, err := t.NewTX()
	if err != nil {
		return err
	}

	defer func() {
		if r := recover(); r != nil {
			newTX.Rollback()
			panic(r)
		}
	}()

	if err := txFunc(ctx, newTX); err != nil {
		newTX.Rollback()
		return err
	}

	return newTX.Commit()
}

func (t tx) getDbInstance() *gorm.DB {
	return t.gormTX
}

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

type genericCommand[DomainModel, GormModel any] struct {
	db  *gorm.DB
	dto func(DomainModel) GormModel
}

func NewGenericCommand[DomainModel, GormModel any](
	db *gorm.DB,
	dto func(DomainModel) GormModel,
) genericCommand[DomainModel, GormModel] {
	return genericCommand[DomainModel, GormModel]{
		db:  db,
		dto: dto,
	}
}

func (command genericCommand[DomainModel, GormModel]) Create(ctx context.Context, req DomainModel) error {
	gormModel := command.dto(req)
	res := command.db.WithContext(ctx).Create(&gormModel)
	return dbDTO.CommandError(res)
}

func (command genericCommand[DomainModel, GormModel]) Update(ctx context.Context, req DomainModel) error {
	gormModel := command.dto(req)
	res := command.db.WithContext(ctx).Save(&gormModel)
	return dbDTO.CommandError(res)
}

func (command genericCommand[DomainModel, GormModel]) Delete(ctx context.Context, req DomainModel) error {
	gormModel := command.dto(req)
	res := command.db.WithContext(ctx).Delete(&gormModel)
	return dbDTO.CommandError(res)
}
