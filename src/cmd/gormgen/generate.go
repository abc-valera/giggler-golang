package main

import (
	"os"

	"gorm.io/driver/sqlite"
	"gorm.io/gen"
	"gorm.io/gorm"
)

func main() {
	g := gen.NewGenerator(gen.Config{
		OutPath:      "src/shared/data/gormgen/gormQuery",
		ModelPkgPath: "src/shared/data/gormgen/gormModel",

		FieldNullable: true,
	})

	dbPath, ok := os.LookupEnv("SQLITE_PATH")
	if !ok {
		panic("SQLITE_PATH env not set")
	}

	db, err := gorm.Open(sqlite.Open(dbPath))
	if err != nil {
		panic(err)
	}

	g.UseDB(db)

	g.ApplyBasic(g.GenerateAllTable()...)

	g.Execute()
}
