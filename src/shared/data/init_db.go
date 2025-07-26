package data

import (
	"gorm.io/gorm"

	"gorm.io/driver/sqlite"

	"giggler-golang/src/shared/env"
	"giggler-golang/src/shared/errutil"
	"giggler-golang/src/shared/serviceLocator"
)

var DB = func() func() *gorm.DB {
	var gormDB *gorm.DB

	switch env.Load("DB") {
	case "sqlite":
		gormDB = errutil.Must(gorm.Open(
			sqlite.Open(env.Load("SQLITE_DSN")),
			&gorm.Config{TranslateError: true},
		))

		serviceLocator.Set(gormDB)
	default:
		panic(env.ErrInvalidEnvValue)
	}

	return func() *gorm.DB { return gormDB }
}()
