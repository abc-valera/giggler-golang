package jokeData

import (
	"context"
	"time"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel/trace"

	"giggler-golang/src/features/joke/jokeDTO"
	"giggler-golang/src/features/joke/jokeModel"
	"giggler-golang/src/shared/data"
	"giggler-golang/src/shared/otel"
	"giggler-golang/src/shared/validate"
)

type command struct {
	dbCommand data.IGenericCommandCreateUpdateDelete[jokeModel.Joke]
	query     query
}

func NewCommand(ds data.IDB) command {
	return command{
		dbCommand: data.NewGenericCommand(getDB(ds), jokeDTO.NewGorm),
		query:     NewQuery(ds),
	}
}

type CreateReq struct {
	Title       string  `validate:"required,min=4,max=64"`
	Text        string  `validate:"required,min=4,max=4096"`
	Explanation *string `validate:"omitempty,max=4096"`

	UserID string `validate:"required,uuid"`
}

func (c command) Create(ctx context.Context, req CreateReq) (jokeModel.Joke, error) {
	var span trace.Span
	ctx, span = otel.Trace(ctx)
	defer span.End()

	if err := validate.Struct(req); err != nil {
		return jokeModel.Joke{}, err
	}

	joke := jokeModel.Joke{
		ID:          uuid.New().String(),
		Title:       req.Title,
		Text:        req.Text,
		Explanation: req.Explanation,
		CreatedAt:   time.Now().Truncate(time.Millisecond).Local(),
		UserID:      req.UserID,
	}

	if err := c.dbCommand.Create(ctx, joke); err != nil {
		return jokeModel.Joke{}, err
	}

	return joke, nil
}

type UpdateReq struct {
	ID          string
	Title       *string `validate:"omitempty,min=4,max=64"`
	Text        *string `validate:"omitempty,min=4,max=4096"`
	Explanation *string `validate:"omitempty,max=4096"`
}

func (c command) Update(ctx context.Context, req UpdateReq) (jokeModel.Joke, error) {
	var span trace.Span
	ctx, span = otel.Trace(ctx)
	defer span.End()

	if err := validate.Struct(req); err != nil {
		return jokeModel.Joke{}, err
	}

	joke, err := c.query.GetByID(ctx, req.ID)
	if err != nil {
		return jokeModel.Joke{}, err
	}

	timeNow := time.Now().Truncate(time.Millisecond).Local()
	joke.UpdatedAt = &timeNow

	if req.Title != nil {
		joke.Title = *req.Title
	}

	if req.Text != nil {
		joke.Text = *req.Text
	}

	if req.Explanation != nil {
		joke.Explanation = req.Explanation
	}

	if err := c.dbCommand.Update(ctx, joke); err != nil {
		return jokeModel.Joke{}, err
	}

	return joke, nil
}

func (c command) Delete(ctx context.Context, id string) error {
	var span trace.Span
	ctx, span = otel.Trace(ctx)
	defer span.End()

	return c.dbCommand.Delete(ctx, jokeModel.Joke{ID: id})
}
