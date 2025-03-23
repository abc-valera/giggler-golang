package data

import (
	"context"

	"github.com/abc-valera/giggler-golang/src/shared/dependency"
	"github.com/abc-valera/giggler-golang/src/shared/env"
	"github.com/abc-valera/giggler-golang/src/shared/errutil"
	"github.com/abc-valera/giggler-golang/src/shared/log"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

var (
	DbVariantGorm = "gorm"
	DbVal         = env.Load("DB")
	DB            = initDB()
)

func initDB() IDS {
	switch DbVal {
	case DbVariantGorm:
		// Data source name (DSN) for the database connection.
		//
		// Consider the following:
		// - Set WAL mode, so readers and writers can access the database concurrently
		// - Set busy timeout, so concurrent writers wait on each other instead of erroring immediately
		// - Enable foreign key checks
		dsn := env.Load("SQLITE_PATH") +
			"?_pragma=journal_mode(WAL)" +
			"&_pragma=busy_timeout(5000)" +
			"&_pragma=foreign_keys(1)"

		return dbGorm{errutil.Must(gorm.Open(sqlite.Open(dsn), &gorm.Config{TranslateError: true}))}
	default:
		panic(env.ErrInvalidEnvValue)
	}
}

type dbGorm struct{ dependency *gorm.DB }

var _ dependency.Interface[*gorm.DB] = dbGorm{}

func (db dbGorm) WithinTX(
	ctx context.Context,
	fn func(context.Context, IDS) error,
) error {
	tx, err := newDbGormTX(ctx, db)
	if err != nil {
		return err
	}

	if err := fn(ctx, tx); err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()

	return nil
}

func (db dbGorm) BeginTX(ctx context.Context) (ITX, error) {
	return newDbGormTX(ctx, db)
}

func (db dbGorm) Dependency() *gorm.DB {
	return db.dependency
}

type dbGormTX struct{ dependency *gorm.DB }

var _ dependency.Interface[*gorm.DB] = dbGormTX{}

func newDbGormTX(ctx context.Context, ds dbGorm) (dbGormTX, error) {
	tx := ds.dependency.WithContext(ctx).Begin()
	if err := tx.Error; err != nil {
		return dbGormTX{}, errutil.NewInternalErr(err)
	}
	return dbGormTX{dependency: tx}, nil
}

func newDbGormNestedTX(ctx context.Context, tx dbGormTX) (dbGormTX, error) {
	nestedTX := tx.dependency.WithContext(ctx).Begin()
	if err := nestedTX.Error; err != nil {
		return dbGormTX{}, errutil.NewInternalErr(err)
	}
	return dbGormTX{dependency: nestedTX}, nil
}

func (tx dbGormTX) Commit() error {
	if err := tx.dependency.Commit().Error; err != nil {
		return errutil.NewInternalErr(err)
	}
	return nil
}

func (tx dbGormTX) Rollback() {
	if err := tx.dependency.Rollback().Error; err != nil {
		log.Error("gorm rollback error", "err", err)
	}
}

func (tx dbGormTX) WithinTX(
	ctx context.Context,
	fn func(context.Context, IDS) error,
) error {
	nestedTX, err := newDbGormNestedTX(ctx, tx)
	if err != nil {
		return err
	}

	if err := fn(ctx, nestedTX); err != nil {
		nestedTX.Rollback()
		return err
	}
	nestedTX.Commit()

	return nil
}

func (tx dbGormTX) BeginTX(ctx context.Context) (ITX, error) {
	return newDbGormNestedTX(ctx, tx)
}

func (tx dbGormTX) Dependency() *gorm.DB {
	return tx.dependency
}
