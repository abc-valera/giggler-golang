// Package joke implements joke authoring: creating, browsing, editing and
// deleting jokes owned by a user.
package joke

import (
	"context"

	"giggler-golang/src/features/identity"
	"giggler-golang/src/shared/errutil"

	"gorm.io/cli/gorm/typed"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Usecase struct {
	db *gorm.DB
}

func New(db *gorm.DB) Usecase {
	return Usecase{
		db: db,
	}
}

type CreateInput struct {
	Title       string
	Text        string
	Explanation *string
}

func (u Usecase) Create(ctx context.Context, req CreateInput) (Joke, error) {
	joke := &Joke{
		Title:       req.Title,
		Text:        req.Text,
		Explanation: req.Explanation,
		UserID:      identity.GetAuthToken(ctx).UserID,
	}

	if err := typed.G[Joke](u.db).Create(ctx, joke); err != nil {
		return Joke{}, err
	}

	return *joke, nil
}

func (u Usecase) List(ctx context.Context) ([]Joke, error) {
	userID := identity.GetAuthToken(ctx).UserID

	jokes, err := typed.G[Joke](u.db).
		Where(JokeHelper.UserID.Eq(userID)).
		Order(clause.OrderBy{Columns: []clause.OrderByColumn{JokeHelper.CreatedAt.Desc()}}).
		Find(ctx)
	if err != nil {
		return nil, err
	}

	return jokes, nil
}

func (u Usecase) Get(ctx context.Context, title string) (Joke, error) {
	userID := identity.GetAuthToken(ctx).UserID

	joke, err := typed.G[Joke](u.db).
		Where(JokeHelper.Title.Eq(title), JokeHelper.UserID.Eq(userID)).
		Take(ctx)
	if err != nil {
		return Joke{}, err
	}

	return joke, nil
}

type UpdateInput struct {
	Title       string
	Text        *string
	Explanation *string
}

func (u Usecase) Update(ctx context.Context, req UpdateInput) (Joke, error) {
	joke, err := u.Get(ctx, req.Title)
	if err != nil {
		return Joke{}, err
	}

	if req.Text != nil {
		joke.Text = *req.Text
	}
	if req.Explanation != nil {
		joke.Explanation = req.Explanation
	}

	if _, err := typed.G[Joke](u.db).
		Where(JokeHelper.Title.Eq(joke.Title), JokeHelper.UserID.Eq(joke.UserID)).
		Updates(ctx, joke); err != nil {
		return Joke{}, err
	}

	return joke, nil
}

func (u Usecase) Delete(ctx context.Context, title string) error {
	userID := identity.GetAuthToken(ctx).UserID

	rowsAffected, err := typed.G[Joke](u.db).
		Where(JokeHelper.Title.Eq(title), JokeHelper.UserID.Eq(userID)).
		Delete(ctx)
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errutil.ErrNotFound
	}

	return nil
}
