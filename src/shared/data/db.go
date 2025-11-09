package data

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"

	"giggler-golang/src/shared/errutil/must"
	"giggler-golang/src/shared/singleton"
)

var GetDb = singleton.New(func() *gorm.DB {
	dsn := fmt.Sprintf(
		"user=%s password=%s dbname=%s host=%s port=%s sslmode=%s",
		must.GetEnv("POSTGRES_USER"),
		must.GetEnv("POSTGRES_PASSWORD"),
		must.GetEnv("POSTGRES_DB"),
		must.GetEnv("POSTGRES_HOST"),
		must.GetEnv("POSTGRES_PORT"),
		must.GetEnv("POSTGRES_SSLMODE"),
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

	return must.Do(gorm.Open(postgres.New(postgresConfig), gormConfig))
})
