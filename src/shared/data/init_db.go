package data

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"giggler-golang/src/shared/env"
	"giggler-golang/src/shared/errutil"
)

var DB = func() func() *gorm.DB {
	var gormDB *gorm.DB

	switch env.Load("DB") {
	case "postgres":
		gormDB = errutil.Must(gorm.Open(
			postgres.Open(env.Load("POSTGRES_DSN")),
			&gorm.Config{TranslateError: true},
		))
	default:
		panic(env.ErrInvalidEnvValue)
	}

	return func() *gorm.DB { return gormDB }
}()
