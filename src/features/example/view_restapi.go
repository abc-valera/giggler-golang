package example

import (
	"context"
	"fmt"

	"github.com/danielgtaylor/huma/v2"
)

// NOTE: this exists for demonstration purposes only

func ApplyRoutes(api huma.API) {
	huma.Get(api, "/hello/{name}", helloHandler)
}

type HelloIn struct {
	Name string `path:"name" maxLength:"30" example:"John" doc:"Name to greet"`
}

type HelloOut struct {
	Body struct {
		Greeting string `json:"message" example:"Hola! John" doc:"Greeting message"`
	}
}

func helloHandler(ctx context.Context, input *HelloIn) (*HelloOut, error) {
	resp := &HelloOut{}
	resp.Body.Greeting = fmt.Sprintf("Hola! %s", input.Name)
	return resp, nil
}
