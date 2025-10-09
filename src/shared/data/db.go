package data

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"

	"giggler-golang/src/shared/errutil/must"
)

var DB = func() func() *gorm.DB {
	dsn := fmt.Sprintf(
		"user=%s password=%s dbname=%s host=%s port=%s sslmode=%s",
		must.Env("POSTGRES_USER"),
		must.Env("POSTGRES_PASSWORD"),
		must.Env("POSTGRES_DB"),
		must.Env("POSTGRES_HOST"),
		must.Env("POSTGRES_PORT"),
		must.Env("POSTGRES_SSLMODE"),
	)
	postgresConfig := postgres.Config{
		DSN:                  dsn,
		PreferSimpleProtocol: true, // disables implicit prepared statement usage
	}
	gormConfig := &gorm.Config{
		TranslateError: true,
		NamingStrategy: schema.NamingStrategy{
			NoLowerCase: true, // Disables snake_case fields conversion
		},
	}

	db := must.Do(gorm.Open(postgres.New(postgresConfig), gormConfig))

	return func() *gorm.DB { return db }
}()
