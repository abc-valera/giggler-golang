package jokeView

import (
	"context"

	"go.opentelemetry.io/otel/trace"

	"giggler-golang/src/features/joke/jokeData"
	"giggler-golang/src/shared/contexts"
	"giggler-golang/src/shared/data"
	"giggler-golang/src/shared/otel"
	"giggler-golang/src/shared/view/viewDTO"
	"giggler-golang/src/shared/view/viewgen"
)

type Handler struct{}

func (Handler) JokesGet(ctx context.Context, params viewgen.JokesGetParams) (viewgen.JokesSchema, error) {
	panic("not implemented")
}

func (Handler) JokesPost(ctx context.Context, req *viewgen.JokesPostReq) (*viewgen.JokeSchema, error) {
	var span trace.Span
	ctx, span = otel.Trace(ctx)
	defer span.End()

	userID, err := contexts.GetUserID(ctx)
	if err != nil {
		return nil, err
	}

	createdJoke, err := jokeData.NewCommand(data.DB()).Create(ctx, jokeData.CreateReq{
		Title:       req.Title,
		Text:        req.Text,
		Explanation: viewDTO.NewDomainPointer(req.Explanation),

		UserID: userID,
	})
	if err != nil {
		return nil, err
	}

	return NewJokeModel(createdJoke), err
}

func (Handler) JokesPut(ctx context.Context, req *viewgen.JokesPutReq) (*viewgen.JokeSchema, error) {
	var span trace.Span
	ctx, span = otel.Trace(ctx)
	defer span.End()

	updatedJoke, err := jokeData.NewCommand(data.DB()).Update(ctx, jokeData.UpdateReq{
		ID:          req.JokeID,
		Title:       viewDTO.NewDomainPointer(req.Title),
		Text:        viewDTO.NewDomainPointer(req.Text),
		Explanation: viewDTO.NewDomainPointer(req.Explanation),
	})
	if err != nil {
		return nil, err
	}

	return NewJokeModel(updatedJoke), err
}

func (Handler) JokesDel(ctx context.Context, req *viewgen.JokesDelReq) error {
	panic("not implemented")
}
