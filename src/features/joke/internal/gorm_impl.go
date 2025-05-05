package internal

import (
	"context"

	"giggler-golang/src/features/joke/jokeDto"
	"giggler-golang/src/features/joke/jokeModel"
	"giggler-golang/src/shared/data"
	"giggler-golang/src/shared/data/gormgen/gormQuery"
)

type GormQ struct{ *gormQuery.Query }

func (q GormQ) GetByID(ctx context.Context, id string) (jokeModel.Joke, error) {
	dto, err := q.WithContext(ctx).Joke.Where(q.Joke.ID.Eq(id)).First()
	return jokeDto.FromGorm(dto), data.GormQueryError(err)
}

func (q GormQ) GetAllByUserID(ctx context.Context, userID string, s data.Selector) ([]jokeModel.Joke, error) {
	dtos, err := q.WithContext(ctx).Joke.Where(q.Joke.UserID.Eq(userID)).Find()
	return jokeDto.FromGorms(dtos), data.GormQueryError(err)
}
