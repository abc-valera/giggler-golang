package main

import (
	"gorm.io/gen"

	"giggler-golang/src/features/joke/jokeData"
	"giggler-golang/src/features/user/userData"
)

func main() {
	g := gen.NewGenerator(gen.Config{
		OutPath: "src/shared/data/gormgenQuery",

		FieldNullable: true,
	})

	g.ApplyBasic(
		userData.User{},
		jokeData.Joke{},
	)

	g.Execute()
}
