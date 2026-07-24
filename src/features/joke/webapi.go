package joke

import (
	"context"
	"net/http"
	"time"

	"giggler-golang/src/features/identity"

	"github.com/danielgtaylor/huma/v2"
	"github.com/google/uuid"
)

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

type Handler struct {
	usecase Usecase
}

func InitHandler(api huma.API, usecase Usecase, identityHandler identity.Handler) Handler {
	h := Handler{usecase: usecase}

	huma.Register(api, huma.Operation{
		OperationID: "joke-create",
		Method:      http.MethodPost,
		Path:        "/jokes",
		Middlewares: huma.Middlewares{identityHandler.AuthMiddleware},
	}, h.create)
	huma.Register(api, huma.Operation{
		OperationID: "joke-list",
		Method:      http.MethodGet,
		Path:        "/jokes",
		Middlewares: huma.Middlewares{identityHandler.AuthMiddleware},
	}, h.list)
	huma.Register(api, huma.Operation{
		OperationID: "joke-get",
		Method:      http.MethodGet,
		Path:        "/jokes/{title}",
		Middlewares: huma.Middlewares{identityHandler.AuthMiddleware},
	}, h.get)
	huma.Register(api, huma.Operation{
		OperationID: "joke-update",
		Method:      http.MethodPatch,
		Path:        "/jokes/{title}",
		Middlewares: huma.Middlewares{identityHandler.AuthMiddleware},
	}, h.update)
	huma.Register(api, huma.Operation{
		OperationID: "joke-delete",
		Method:      http.MethodDelete,
		Path:        "/jokes/{title}",
		Middlewares: huma.Middlewares{identityHandler.AuthMiddleware},
	}, h.delete)

	return h
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

func (h Handler) create(ctx context.Context, input *createInput) (*createOutput, error) {
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

func (h Handler) list(ctx context.Context, _ *listInput) (*listOutput, error) {
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

func (h Handler) get(ctx context.Context, input *getInput) (*getOutput, error) {
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

func (h Handler) update(ctx context.Context, input *updateInput) (*updateOutput, error) {
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

func (h Handler) delete(ctx context.Context, input *deleteInput) (*deleteOutput, error) {
	if err := h.usecase.Delete(ctx, input.Title); err != nil {
		return nil, err
	}
	return &deleteOutput{}, nil
}
