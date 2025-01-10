package user

import (
	"context"

	"github.com/abc-valera/giggler-golang/gen/view"
	"github.com/abc-valera/giggler-golang/src/components/viewutil"
	"github.com/abc-valera/giggler-golang/src/shared/otel"
	"github.com/abc-valera/giggler-golang/src/shared/persistence"
	"go.opentelemetry.io/otel/trace"
)

type View struct{}

func (View) MeGet(ctx context.Context) (*view.User, error) {
	var span trace.Span
	ctx, span = otel.Trace(ctx)
	defer span.End()

	userID, err := viewutil.GetUserID(ctx)
	if err != nil {
		return nil, err
	}

	user, err := Query(persistence.DB).GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	return ViewNewUserDTO(user), nil
}

func (View) MePut(ctx context.Context, req *view.MePutReq) (*view.User, error) {
	var span trace.Span
	ctx, span = otel.Trace(ctx)
	defer span.End()

	userID, err := viewutil.GetUserID(ctx)
	if err != nil {
		return nil, err
	}

	user, err := Update(ctx, persistence.DB, UpdateRequest{
		ID:       userID,
		Password: viewutil.NewString(req.Password),
		Fullname: viewutil.NewString(req.Fullname),
		Status:   viewutil.NewString(req.Status),
	})
	if err != nil {
		return nil, err
	}
	return ViewNewUserDTO(user), nil
}

func (View) MeDel(ctx context.Context, req *view.MeDelReq) error {
	var span trace.Span
	ctx, span = otel.Trace(ctx)
	defer span.End()

	userID, err := viewutil.GetUserID(ctx)
	if err != nil {
		return err
	}

	return Delete(ctx, persistence.DB, DeleteRequest{
		ID:       userID,
		Password: req.Password,
	})
}

func ViewNewUserDTO(user User) *view.User {
	return &view.User{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		Fullname:  view.OptString{Value: user.Fullname, Set: user.Fullname != ""},
		Status:    view.OptString{Value: user.Status, Set: user.Status != ""},
		CreatedAt: user.CreatedAt,
	}
}
