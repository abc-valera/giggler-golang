package jokeData

import (
	"context"
	"time"

	"github.com/abc-valera/giggler-golang/src/components/data"
	"github.com/abc-valera/giggler-golang/src/components/env"
	"github.com/abc-valera/giggler-golang/src/features/joke/jokeModel"
	"github.com/abc-valera/giggler-golang/src/shared/ds"
	"github.com/abc-valera/giggler-golang/src/shared/otel"
	"github.com/abc-valera/giggler-golang/src/shared/validate"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel/trace"
)

var command = initCommand()

type iCommand data.IGenericCommandCreateUpdateDelete[jokeModel.Joke]

func initCommand() func(dataStore ds.IDS) iCommand {
	switch ds.DbVal {
	case ds.DbVariantGorm:
		return func(dataStore ds.IDS) iCommand {
			return data.NewGormGenericCommand(ds.GormDS(dataStore), jokeModel.NewGormDTO)
		}
	default:
		panic(env.ErrInvalidEnvValue)
	}
}

type CreateReq struct {
	Title       string  `validate:"required,min=4,max=64"`
	Text        string  `validate:"required,min=4,max=4096"`
	Explanation *string `validate:"omitempty,max=4096"`

	UserID string `validate:"required,uuid"`
}

func Create(ctx context.Context, ds ds.IDS, req CreateReq) (jokeModel.Joke, error) {
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

	if err := command(ds).Create(ctx, joke); err != nil {
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

func Update(ctx context.Context, ds ds.IDS, req UpdateReq) (jokeModel.Joke, error) {
	var span trace.Span
	ctx, span = otel.Trace(ctx)
	defer span.End()

	if err := validate.Struct(req); err != nil {
		return jokeModel.Joke{}, err
	}

	joke, err := Query(ds).GetByID(ctx, req.ID)
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

	if err := command(ds).Update(ctx, joke); err != nil {
		return jokeModel.Joke{}, err
	}

	return joke, nil
}

func Delete(ctx context.Context, ds ds.IDS, id string) error {
	var span trace.Span
	ctx, span = otel.Trace(ctx)
	defer span.End()

	return command(ds).Delete(ctx, jokeModel.Joke{ID: id})
}
