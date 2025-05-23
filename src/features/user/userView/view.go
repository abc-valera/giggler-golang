package userView

import (
	"context"

	"github.com/abc-valera/giggler-golang/src/features/joke/jokeData"
	"github.com/abc-valera/giggler-golang/src/features/joke/jokeModel"
	"github.com/abc-valera/giggler-golang/src/features/user/userData"
	"github.com/abc-valera/giggler-golang/src/features/user/userModel"
	"github.com/abc-valera/giggler-golang/src/shared/contexts"
	"github.com/abc-valera/giggler-golang/src/shared/data"
	"github.com/abc-valera/giggler-golang/src/shared/otel"
	"github.com/abc-valera/giggler-golang/src/shared/view"
	"github.com/abc-valera/giggler-golang/src/shared/view/viewgen"
	"go.opentelemetry.io/otel/trace"
)

type View struct{}

func (View) UserGet(ctx context.Context) (*viewgen.UserSchema, error) {
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
	return userModel.NewViewDTO(user), nil
}

func (View) UserPut(ctx context.Context, req *viewgen.UserPutReq) (*viewgen.UserSchema, error) {
	var span trace.Span
	ctx, span = otel.Trace(ctx)
	defer span.End()

	userID, err := contexts.GetUserID(ctx)
	if err != nil {
		return nil, err
	}

	user, err := userData.Update(ctx, data.DB, userData.UpdateReq{
		ID:       userID,
		Password: view.NewDomainPointer(req.Password),
		Fullname: view.NewDomainPointer(req.Fullname),
		Status:   view.NewDomainPointer(req.Status),
	})
	if err != nil {
		return nil, err
	}
	return userModel.NewViewDTO(user), nil
}

func (View) UserDel(ctx context.Context, req *viewgen.UserDelReq) error {
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

func (View) UserJokesGet(ctx context.Context, params viewgen.UserJokesGetParams) (viewgen.JokesSchema, error) {
	panic("not implemented")
}

func (View) UserJokesPost(ctx context.Context, req *viewgen.UserJokesPostReq) (*viewgen.JokeSchema, error) {
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
		Explanation: view.NewDomainPointer(req.Explanation),

		UserID: userID,
	})
	if err != nil {
		return nil, err
	}

	return jokeModel.NewViewDTO(createdJoke), err
}

func (View) UserJokesPut(ctx context.Context, req *viewgen.UserJokesPutReq) (*viewgen.JokeSchema, error) {
	var span trace.Span
	ctx, span = otel.Trace(ctx)
	defer span.End()

	updatedJoke, err := jokeData.Update(ctx, data.DB, jokeData.UpdateReq{
		ID:          req.JokeID,
		Title:       view.NewDomainPointer(req.Title),
		Text:        view.NewDomainPointer(req.Text),
		Explanation: view.NewDomainPointer(req.Explanation),
	})
	if err != nil {
		return nil, err
	}

	return jokeModel.NewViewDTO(updatedJoke), err
}

func (View) UserJokesDel(ctx context.Context, req *viewgen.UserJokesDelReq) error {
	panic("not implemented")
}
