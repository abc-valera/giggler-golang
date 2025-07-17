package jokeData

import (
	"context"
	"time"

	"giggler-golang/src/shared/data/dataModel"
	"giggler-golang/src/shared/data/dbgen/gormModel"
	"giggler-golang/src/shared/otel"
	"giggler-golang/src/shared/validate"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type command struct {
	dbCommand dataModel.IGenericCommandCreateUpdateDelete[gormModel.Joke]
	query     query
}

func NewCommand(db *gorm.DB) command {
	return command{
		dbCommand: dataModel.NewGenericCommand[gormModel.Joke](db),
		query:     NewQuery(db),
	}
}

type CreateReq struct {
	Title       string  `validate:"required,min=4,max=64"`
	Text        string  `validate:"required,min=4,max=4096"`
	Explanation *string `validate:"omitempty,max=4096"`

	UserID string `validate:"required,uuid"`
}

func (c command) Create(ctx context.Context, req CreateReq) (*gormModel.Joke, error) {
	ctx, span := otel.Trace(ctx)
	defer span.End()

	if err := validate.Struct(req); err != nil {
		return nil, err
	}

	joke := &gormModel.Joke{
		ID:          uuid.New().String(),
		Title:       req.Title,
		Text:        req.Text,
		Explanation: req.Explanation,
		CreatedAt:   time.Now().Truncate(time.Millisecond).Local(),
		UserID:      req.UserID,
	}

	if err := c.dbCommand.Create(ctx, joke); err != nil {
		return nil, err
	}

	return joke, nil
}

type UpdateReq struct {
	ID          string
	Title       *string `validate:"omitempty,min=4,max=64"`
	Text        *string `validate:"omitempty,min=4,max=4096"`
	Explanation *string `validate:"omitempty,max=4096"`
}

func (c command) Update(ctx context.Context, req UpdateReq) (*gormModel.Joke, error) {
	ctx, span := otel.Trace(ctx)
	defer span.End()

	if err := validate.Struct(req); err != nil {
		return nil, err
	}

	joke, err := c.query.GetByID(ctx, req.ID)
	if err != nil {
		return nil, err
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
		return nil, err
	}

	return joke, nil
}

func (c command) Delete(ctx context.Context, id string) error {
	ctx, span := otel.Trace(ctx)
	defer span.End()

	return c.dbCommand.Delete(ctx, &gormModel.Joke{ID: id})
}
