package userData

import (
	"context"
	"time"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel/trace"

	"giggler-golang/src/features/auth/passworder"
	"giggler-golang/src/features/user/userDTO"
	"giggler-golang/src/features/user/userModel"
	"giggler-golang/src/shared/data"
	"giggler-golang/src/shared/otel"
	"giggler-golang/src/shared/validate"
)

type command struct {
	dbCommand data.IGenericCommandCreateUpdateDelete[userModel.User]
	query     query
}

func NewCommand(ds data.IDB) command {
	return command{
		dbCommand: data.NewGenericCommand(getDB(ds), userDTO.NewGorm),
		query:     NewQuery(ds),
	}
}

type CreateReq struct {
	Username string  `validate:"required,min=2,max=32"`
	Email    string  `validate:"required,email"`
	Password string  `validate:"required,min=2,max=32"`
	Fullname *string `validate:"omitempty,max=64"`
	Status   *string `validate:"omitempty,max=128"`
}

func (c command) Create(ctx context.Context, req CreateReq) (userModel.User, error) {
	var span trace.Span
	ctx, span = otel.Trace(ctx)
	defer span.End()

	if err := validate.Struct(req); err != nil {
		return userModel.User{}, err
	}

	user := userModel.User{
		ID:             uuid.New().String(),
		Username:       req.Username,
		Email:          req.Email,
		HashedPassword: passworder.Hash(ctx, req.Password),
		Fullname:       req.Fullname,
		Status:         req.Status,
		CreatedAt:      time.Now().Truncate(time.Millisecond).Local(),
	}

	if err := c.dbCommand.Create(ctx, user); err != nil {
		return userModel.User{}, err
	}

	return user, nil
}

type UpdateReq struct {
	ID       string
	Password *string `validate:"omitempty,min=2,max=32"`
	Fullname *string `validate:"omitempty,max=64"`
	Status   *string `validate:"omitempty,max=128"`
}

func (c command) Update(ctx context.Context, req UpdateReq) (userModel.User, error) {
	var span trace.Span
	ctx, span = otel.Trace(ctx)
	defer span.End()

	if err := validate.Struct(req); err != nil {
		return userModel.User{}, err
	}

	user, err := c.query.GetByID(ctx, req.ID)
	if err != nil {
		return userModel.User{}, err
	}

	timeNow := time.Now().Truncate(time.Millisecond).Local()
	user.UpdatedAt = &timeNow

	if req.Password != nil {
		span.AddEvent("Hashing Password Start")
		user.HashedPassword = passworder.Hash(ctx, *req.Password)
		span.AddEvent("Hashing Password End")
	}

	if req.Fullname != nil {
		user.Fullname = req.Fullname
	}

	if req.Status != nil {
		user.Status = req.Status
	}

	if err := c.dbCommand.Update(ctx, user); err != nil {
		return userModel.User{}, err
	}

	return user, nil
}

type DeleteReq struct {
	ID       string
	Password string `validate:"required"`
}

func (c command) Delete(ctx context.Context, req DeleteReq) error {
	var span trace.Span
	ctx, span = otel.Trace(ctx)
	defer span.End()

	if err := validate.Struct(req); err != nil {
		return err
	}

	user, err := c.query.GetByID(ctx, req.ID)
	if err != nil {
		return err
	}

	span.AddEvent("Checking Password Start")

	if !passworder.IsReal(ctx, req.Password, user.HashedPassword) {
		return passworder.ErrInvalidPass
	}

	span.AddEvent("Checking Password End")

	return c.dbCommand.Delete(ctx, userModel.User{ID: req.ID})
}
