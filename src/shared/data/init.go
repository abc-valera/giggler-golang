package data

import (
	"fmt"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"giggler-golang/src/shared/data/internal/localFS"
	"giggler-golang/src/shared/env"
	"giggler-golang/src/shared/errutil"
	"giggler-golang/src/shared/serviceLocator"
)

var DB = func() func() IDB {
	serviceLocator.Set(func(ds IDB) *gorm.DB {
		return ds.getDbInstance()
	})

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

	ds := db{gormDB: gormDB}

	return func() IDB { return ds }
}()

var FS = func() func() iFS {
	var fs iFS

	switch env.Load("FS") {
	case "local":
		folderPath := env.Load("LOCAL_DSN")

		if err := os.MkdirAll(folderPath, 0o755); err != nil {
			if !os.IsExist(err) {
				panic(fmt.Errorf("failed to create local file saver directory: %w", err))
			}
		}

		fs = localFS.New(folderPath)
	default:
		panic(env.ErrInvalidEnvValue)
	}

	return func() iFS { return fs }
}()
