package main

import (
	"fmt"

	"giggler-golang/src/shared/env"
	"giggler-golang/src/shared/errutil"

	"gorm.io/driver/postgres"
	"gorm.io/gen"
	"gorm.io/gorm"
)

func main() {
	g := gen.NewGenerator(gen.Config{
		OutPath:      "src/shared/data/gormgen/gormQuery",
		ModelPkgPath: "src/shared/data/gormgen/gormModel",

		FieldNullable: true,
	})

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		env.Load("POSTGRES_HOST"),
		env.Load("POSTGRES_PORT"),
		env.Load("POSTGRES_USER"),
		env.Load("POSTGRES_PASSWORD"),
		env.Load("POSTGRES_DB"),
		env.Load("POSTGRES_SSLMODE"),
	)

	db := errutil.Must(gorm.Open(postgres.Open(dsn)))

	g.UseDB(db)

	g.ApplyBasic(g.GenerateAllTable()...)

	g.Execute()
}
