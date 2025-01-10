package joke

import (
	"context"

	"github.com/abc-valera/giggler-golang/gen/view"
	"github.com/abc-valera/giggler-golang/src/components/viewutil"
	"github.com/abc-valera/giggler-golang/src/shared/otel"
	"github.com/abc-valera/giggler-golang/src/shared/persistence"
	"go.opentelemetry.io/otel/trace"
)

type View struct{}

func (View) JokesGet(ctx context.Context, params view.JokesGetParams) (view.Jokes, error) {
	var span trace.Span
	ctx, span = otel.Trace(ctx)
	defer span.End()

	userID, err := viewutil.GetUserID(ctx)
	if err != nil {
		return nil, err
	}

	jokes, err := Query(persistence.DB).GetAllByUserID(
		ctx,
		userID,
		viewutil.NewSelector(params.Limit, params.Offset),
	)

	return NewViewJokes(jokes), err
}

func (View) JokesPost(ctx context.Context, req *view.JokesPostReq) (*view.Joke, error) {
	var span trace.Span
	ctx, span = otel.Trace(ctx)
	defer span.End()

	userID, err := viewutil.GetUserID(ctx)
	if err != nil {
		return nil, err
	}

	joke, err := Create(ctx, persistence.DB, CreateRequest{
		UserID:      userID,
		Title:       req.Title,
		Text:        req.Text,
		Explanation: req.Explanation.Value,
	})
	if err != nil {
		return nil, err
	}
	return NewViewJoke(joke), err
}

func (View) JokesPut(ctx context.Context, req *view.JokesPutReq) (*view.Joke, error) {
	var span trace.Span
	ctx, span = otel.Trace(ctx)
	defer span.End()

	joke, err := Update(ctx, persistence.DB, UpdateRequest{
		ID:          req.JokeID,
		Title:       viewutil.NewString(req.Title),
		Text:        viewutil.NewString(req.Text),
		Explanation: viewutil.NewString(req.Explanation),
	})
	if err != nil {
		return nil, err
	}
	return NewViewJoke(joke), err
}

func (View) JokesDel(ctx context.Context, req *view.JokesDelReq) error {
	var span trace.Span
	ctx, span = otel.Trace(ctx)
	defer span.End()

	return Delete(ctx, persistence.DB, req.JokeID)
}

func NewViewJoke(joke Joke) *view.Joke {
	return &view.Joke{
		ID:          joke.ID,
		UserID:      joke.UserID,
		Title:       joke.Title,
		Text:        joke.Text,
		Explanation: view.OptString{Value: joke.Explanation, Set: joke.Explanation != ""},
		CreatedAt:   joke.CreatedAt,
	}
}

func NewViewJokes(jokes []Joke) view.Jokes {
	viewJokes := make(view.Jokes, len(jokes))
	for i, joke := range jokes {
		viewJokes[i] = *NewViewJoke(joke)
	}
	return viewJokes
}
