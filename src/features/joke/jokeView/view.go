package jokeView

import (
	"context"

	"giggler-golang/src/features/joke/jokeData"
	"giggler-golang/src/features/joke/jokeDto"
	"giggler-golang/src/shared/contexts"
	"giggler-golang/src/shared/data"
	"giggler-golang/src/shared/otel"
	"giggler-golang/src/shared/view/viewUtil"
	"giggler-golang/src/shared/view/viewgen"

	"go.opentelemetry.io/otel/trace"
)

type Handler struct{}

func (Handler) JokesGet(ctx context.Context, params viewgen.JokesGetParams) (viewgen.JokesSchema, error) {
	var span trace.Span
	ctx, span = otel.Trace(ctx)
	defer span.End()

	userID, err := contexts.GetUserID(ctx)
	if err != nil {
		return nil, err
	}

	jokes, err := jokeData.Query(data.DB).GetAllByUserID(
		ctx,
		userID,
		viewUtil.NewSelector(params.Limit, params.Offset),
	)

	return jokeDto.NewViewDTOs(jokes), err
}
