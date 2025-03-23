package main

import (
	"gorm.io/driver/postgres"
	"gorm.io/gen"
	"gorm.io/gorm"

	"giggler-golang/src/shared/env"
	"giggler-golang/src/shared/errutil"
)

func main() {
	g := gen.NewGenerator(gen.Config{
		OutPath:      "src/shared/data/dbgen/gormQuery",
		ModelPkgPath: "src/shared/data/dbgen/gormModel",

		FieldNullable: true,
	})

	g.UseDB(errutil.Must(gorm.Open(postgres.Open(env.Load("POSTGRES_DSN")))))

	g.ApplyBasic(g.GenerateAllTable()...)

	g.Execute()
}
