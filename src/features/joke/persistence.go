package joke

import (
	"context"
	"time"

	"github.com/abc-valera/giggler-golang/src/shared/otel"
	"github.com/abc-valera/giggler-golang/src/shared/persistence"
	"github.com/abc-valera/giggler-golang/src/shared/validate"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel/trace"
)

type CreateRequest struct {
	Title       string `validate:"required,min=4,max=64"`
	Text        string `validate:"required,min=4,max=4096"`
	Explanation string `validate:"max=4096"`

	UserID string `validate:"required,uuid"`
}

func Create(ctx context.Context, ds persistence.IDS, req CreateRequest) (Joke, error) {
	var span trace.Span
	ctx, span = otel.Trace(ctx)
	defer span.End()

	if err := validate.Struct(req); err != nil {
		return Joke{}, err
	}

	joke := Joke{
		ID:          uuid.New().String(),
		Title:       req.Title,
		Text:        req.Text,
		Explanation: req.Explanation,
		CreatedAt:   time.Now().Truncate(time.Millisecond).Local(),
		UserID:      req.UserID,
	}

	if err := command(ds).Create(ctx, joke); err != nil {
		return Joke{}, err
	}

	return joke, nil
}

type UpdateRequest struct {
	ID          string
	Title       *string `validate:"min=4,max=64"`
	Text        *string `validate:"min=4,max=4096"`
	Explanation *string `validate:"max=4096"`
}

func Update(ctx context.Context, ds persistence.IDS, req UpdateRequest) (Joke, error) {
	var span trace.Span
	ctx, span = otel.Trace(ctx)
	defer span.End()

	if err := validate.Struct(req); err != nil {
		return Joke{}, err
	}

	joke, err := Query(ds).GetByID(ctx, req.ID)
	if err != nil {
		return Joke{}, err
	}

	joke.UpdatedAt = time.Now().Truncate(time.Millisecond).Local()

	if req.Title != nil {
		joke.Title = *req.Title
	}

	if req.Text != nil {
		joke.Text = *req.Text
	}

	if req.Explanation != nil {
		joke.Explanation = *req.Explanation
	}

	if err := command(ds).Update(ctx, joke); err != nil {
		return Joke{}, err
	}

	return joke, nil
}

func Delete(ctx context.Context, ds persistence.IDS, id string) error {
	var span trace.Span
	ctx, span = otel.Trace(ctx)
	defer span.End()

	return command(ds).Delete(ctx, Joke{ID: id})
}
