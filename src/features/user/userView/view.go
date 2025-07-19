package userView

import (
	"context"

	"giggler-golang/src/features/user/userData/userRepo"
	"giggler-golang/src/shared/contexts"
	"giggler-golang/src/shared/data"
	"giggler-golang/src/shared/otel"
	"giggler-golang/src/shared/view/viewDTO"
	"giggler-golang/src/shared/view/viewgen"
)

type Handler struct{}

func (Handler) UserGet(ctx context.Context) (*viewgen.UserSchema, error) {
	ctx, span := otel.Trace(ctx)
	defer span.End()

	userID, err := contexts.GetUserID(ctx)
	if err != nil {
		return nil, err
	}

	user, err := userRepo.NewQuery(data.DB()).GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	return NewUserModel(user), nil
}

func (Handler) UserPut(ctx context.Context, req *viewgen.UserPutReq) (*viewgen.UserSchema, error) {
	ctx, span := otel.Trace(ctx)
	defer span.End()

	userID, err := contexts.GetUserID(ctx)
	if err != nil {
		return nil, err
	}

	user, err := userRepo.NewCommand(data.DB()).Update(ctx, userRepo.UpdateReq{
		ID:       userID,
		Password: viewDTO.NewDomainPointer(req.Password),
		Fullname: viewDTO.NewDomainPointer(req.Fullname),
		Status:   viewDTO.NewDomainPointer(req.Status),
	})
	if err != nil {
		return nil, err
	}
	return NewUserModel(user), nil
}

func (Handler) UserDel(ctx context.Context, req *viewgen.UserDelReq) error {
	ctx, span := otel.Trace(ctx)
	defer span.End()

	userID, err := contexts.GetUserID(ctx)
	if err != nil {
		return err
	}

	return userRepo.NewCommand(data.DB()).Delete(ctx, userRepo.DeleteReq{
		ID:       userID,
		Password: req.Password,
	})
}
