package view

import (
	"context"
	"fmt"
)

type HelloInput struct {
	Name string `path:"name" maxLength:"30" example:"John" doc:"Name to greet"`
}

type HelloOutput struct {
	Body struct {
		Greeting string `json:"message" example:"Hola! John" doc:"Greeting message"`
	}
}

func helloHandler(ctx context.Context, input *HelloInput) (*HelloOutput, error) {
	resp := &HelloOutput{}
	resp.Body.Greeting = fmt.Sprintf("Hola! %s", input.Name)
	return resp, nil
}
