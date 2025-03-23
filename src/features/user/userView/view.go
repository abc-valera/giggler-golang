package userView

import (
	"context"

	"giggler-golang/src/features/joke/jokeData"
	"giggler-golang/src/features/joke/jokeDto"
	"giggler-golang/src/features/user/userData"
	"giggler-golang/src/features/user/userDto"
	"giggler-golang/src/shared/contexts"
	"giggler-golang/src/shared/data"
	"giggler-golang/src/shared/otel"
	"giggler-golang/src/shared/view/viewUtil"
	"giggler-golang/src/shared/view/viewgen"

	"go.opentelemetry.io/otel/trace"
)

type Handler struct{}

func (Handler) UserGet(ctx context.Context) (*viewgen.UserSchema, error) {
	var span trace.Span
	ctx, span = otel.Trace(ctx)
	defer span.End()

	userID, err := contexts.GetUserID(ctx)
	if err != nil {
		return nil, err
	}

	user, err := userData.Query(data.DB).GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	return userDto.NewViewDTO(user), nil
}

func (Handler) UserPut(ctx context.Context, req *viewgen.UserPutReq) (*viewgen.UserSchema, error) {
	var span trace.Span
	ctx, span = otel.Trace(ctx)
	defer span.End()

	userID, err := contexts.GetUserID(ctx)
	if err != nil {
		return nil, err
	}

	user, err := userData.Update(ctx, data.DB, userData.UpdateReq{
		ID:       userID,
		Password: viewUtil.NewDomainPointer(req.Password),
		Fullname: viewUtil.NewDomainPointer(req.Fullname),
		Status:   viewUtil.NewDomainPointer(req.Status),
	})
	if err != nil {
		return nil, err
	}
	return userDto.NewViewDTO(user), nil
}

func (Handler) UserDel(ctx context.Context, req *viewgen.UserDelReq) error {
	var span trace.Span
	ctx, span = otel.Trace(ctx)
	defer span.End()

	userID, err := contexts.GetUserID(ctx)
	if err != nil {
		return err
	}

	return userData.Delete(ctx, data.DB, userData.DeleteReq{
		ID:       userID,
		Password: req.Password,
	})
}

func (Handler) UserJokesGet(ctx context.Context, params viewgen.UserJokesGetParams) (viewgen.JokesSchema, error) {
	panic("not implemented")
}

func (Handler) UserJokesPost(ctx context.Context, req *viewgen.UserJokesPostReq) (*viewgen.JokeSchema, error) {
	var span trace.Span
	ctx, span = otel.Trace(ctx)
	defer span.End()

	userID, err := contexts.GetUserID(ctx)
	if err != nil {
		return nil, err
	}

	createdJoke, err := jokeData.Create(ctx, data.DB, jokeData.CreateReq{
		Title:       req.Title,
		Text:        req.Text,
		Explanation: viewUtil.NewDomainPointer(req.Explanation),

		UserID: userID,
	})
	if err != nil {
		return nil, err
	}

	return jokeDto.NewViewDTO(createdJoke), err
}

func (Handler) UserJokesPut(ctx context.Context, req *viewgen.UserJokesPutReq) (*viewgen.JokeSchema, error) {
	var span trace.Span
	ctx, span = otel.Trace(ctx)
	defer span.End()

	updatedJoke, err := jokeData.Update(ctx, data.DB, jokeData.UpdateReq{
		ID:          req.JokeID,
		Title:       viewUtil.NewDomainPointer(req.Title),
		Text:        viewUtil.NewDomainPointer(req.Text),
		Explanation: viewUtil.NewDomainPointer(req.Explanation),
	})
	if err != nil {
		return nil, err
	}

	return jokeDto.NewViewDTO(updatedJoke), err
}

func (Handler) UserJokesDel(ctx context.Context, req *viewgen.UserJokesDelReq) error {
	panic("not implemented")
}
