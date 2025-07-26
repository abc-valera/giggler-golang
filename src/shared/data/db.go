package data

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"

	"giggler-golang/src/shared/env"
	"giggler-golang/src/shared/errutil/must"
)

var DB = func() func() *gorm.DB {
	dsn := fmt.Sprintf(
		"user=%s password=%s dbname=%s host=%s port=%s sslmode=%s",
		env.Load("POSTGRES_USER"),
		env.Load("POSTGRES_PASSWORD"),
		env.Load("POSTGRES_DB"),
		env.Load("POSTGRES_HOST"),
		env.Load("POSTGRES_PORT"),
		env.Load("POSTGRES_SSLMODE"),
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
