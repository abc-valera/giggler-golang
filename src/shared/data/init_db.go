package data

import (
	"gorm.io/gorm"

	"gorm.io/driver/sqlite"

	"giggler-golang/src/features/joke/jokeData"
	"giggler-golang/src/features/user/userData"
	"giggler-golang/src/shared/env"
	"giggler-golang/src/shared/errutil"
)

var DB = func() func() *gorm.DB {
	var gormDB *gorm.DB

	switch env.Load("DB") {
	case "sqlite":
		gormDB = errutil.Must(gorm.Open(
			sqlite.Open(env.Load("SQLITE_DSN")),
			&gorm.Config{TranslateError: true},
		))
	default:
		panic(env.ErrInvalidEnvValue)
	}

	gormDB.AutoMigrate(
		userData.User{},
		jokeData.Joke{},
	)

	return func() *gorm.DB { return gormDB }
}()
