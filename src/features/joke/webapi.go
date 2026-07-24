package joke

import (
	"context"
	"net/http"
	"time"

	"github.com/danielgtaylor/huma/v2"
	"github.com/google/uuid"
)

func InitRoutes(api huma.API, usecase Usecase, authMiddleware func(huma.Context, func(huma.Context))) {
	h := handler{usecase: usecase}

	huma.Register(api, huma.Operation{
		OperationID: "joke-create",
		Method:      http.MethodPost,
		Path:        "/jokes",
		Middlewares: huma.Middlewares{authMiddleware},
	}, h.createHandler)
	huma.Register(api, huma.Operation{
		OperationID: "joke-list",
		Method:      http.MethodGet,
		Path:        "/jokes",
		Middlewares: huma.Middlewares{authMiddleware},
	}, h.listHandler)
	huma.Register(api, huma.Operation{
		OperationID: "joke-get",
		Method:      http.MethodGet,
		Path:        "/jokes/{title}",
		Middlewares: huma.Middlewares{authMiddleware},
	}, h.getHandler)
	huma.Register(api, huma.Operation{
		OperationID: "joke-update",
		Method:      http.MethodPatch,
		Path:        "/jokes/{title}",
		Middlewares: huma.Middlewares{authMiddleware},
	}, h.updateHandler)
	huma.Register(api, huma.Operation{
		OperationID: "joke-delete",
		Method:      http.MethodDelete,
		Path:        "/jokes/{title}",
		Middlewares: huma.Middlewares{authMiddleware},
	}, h.deleteHandler)
}

type WebapiDTO struct {
	UserID      uuid.UUID
	Title       string
	Text        string
	Explanation *string
	CreatedAt   time.Time
}

func NewModel(joke Joke) *WebapiDTO {
	return &WebapiDTO{
		UserID:      joke.UserID,
		Title:       joke.Title,
		Text:        joke.Text,
		Explanation: joke.Explanation,
		CreatedAt:   joke.CreatedAt,
	}
}

func NewModels(jokes []Joke) []*WebapiDTO {
	models := make([]*WebapiDTO, len(jokes))
	for i, joke := range jokes {
		models[i] = NewModel(joke)
	}
	return models
}

type handler struct {
	usecase Usecase
}

type createInput struct {
	Body struct {
		Title       string  `example:"Why did the chicken cross the road?"`
		Text        string  `example:"To get to the other side."`
		Explanation *string `example:"It's a classic anti-joke."`
	}
}

type createOutput struct {
	Body struct {
		Joke *WebapiDTO
	}
}

func (h handler) createHandler(ctx context.Context, input *createInput) (*createOutput, error) {
	joke, err := h.usecase.Create(ctx, CreateInput{
		Title:       input.Body.Title,
		Text:        input.Body.Text,
		Explanation: input.Body.Explanation,
	})
	if err != nil {
		return nil, err
	}

	var out createOutput
	out.Body.Joke = NewModel(joke)
	return &out, nil
}

type listInput struct{}

type listOutput struct {
	Body struct {
		Jokes []*WebapiDTO
	}
}

func (h handler) listHandler(ctx context.Context, _ *listInput) (*listOutput, error) {
	jokes, err := h.usecase.List(ctx)
	if err != nil {
		return nil, err
	}

	var out listOutput
	out.Body.Jokes = NewModels(jokes)
	return &out, nil
}

type getInput struct {
	Title string `path:"title"`
}

type getOutput struct {
	Body struct {
		Joke *WebapiDTO
	}
}

func (h handler) getHandler(ctx context.Context, input *getInput) (*getOutput, error) {
	joke, err := h.usecase.Get(ctx, input.Title)
	if err != nil {
		return nil, err
	}

	var out getOutput
	out.Body.Joke = NewModel(joke)
	return &out, nil
}

type updateInput struct {
	Title string `path:"title"`
	Body  struct {
		Text        *string `example:"To get to the other side."`
		Explanation *string `example:"It's a classic anti-joke."`
	}
}

type updateOutput struct {
	Body struct {
		Joke *WebapiDTO
	}
}

func (h handler) updateHandler(ctx context.Context, input *updateInput) (*updateOutput, error) {
	joke, err := h.usecase.Update(ctx, UpdateInput{
		Title:       input.Title,
		Text:        input.Body.Text,
		Explanation: input.Body.Explanation,
	})
	if err != nil {
		return nil, err
	}

	var out updateOutput
	out.Body.Joke = NewModel(joke)
	return &out, nil
}

type deleteInput struct {
	Title string `path:"title"`
}

type deleteOutput struct{}

func (h handler) deleteHandler(ctx context.Context, input *deleteInput) (*deleteOutput, error) {
	if err := h.usecase.Delete(ctx, input.Title); err != nil {
		return nil, err
	}
	return &deleteOutput{}, nil
}
