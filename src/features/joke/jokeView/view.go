package jokeView

import (
	"context"

	"github.com/abc-valera/giggler-golang/src/components/viewgen"
	"github.com/abc-valera/giggler-golang/src/components/viewutil"
	"github.com/abc-valera/giggler-golang/src/features/joke/jokeData"
	"github.com/abc-valera/giggler-golang/src/features/joke/jokeModel"
	"github.com/abc-valera/giggler-golang/src/shared/data"
	"github.com/abc-valera/giggler-golang/src/shared/otel"
	"go.opentelemetry.io/otel/trace"
)

type View struct{}

func (View) JokesGet(ctx context.Context, params viewgen.JokesGetParams) (viewgen.JokesSchema, error) {
	var span trace.Span
	ctx, span = otel.Trace(ctx)
	defer span.End()

	userID, err := viewutil.GetUserID(ctx)
	if err != nil {
		return nil, err
	}

	jokes, err := jokeData.Query(data.DB).GetAllByUserID(
		ctx,
		userID,
		viewutil.NewSelector(params.Limit, params.Offset),
	)

	return jokeModel.NewViewDTOs(jokes), err
}
