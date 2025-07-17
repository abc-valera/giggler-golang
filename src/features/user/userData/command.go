package userData

import (
	"context"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"giggler-golang/src/features/user/userPassword"
	"giggler-golang/src/shared/data/dataModel"
	"giggler-golang/src/shared/data/dbgen/gormModel"
	"giggler-golang/src/shared/otel"
	"giggler-golang/src/shared/validate"
)

type command struct {
	dbCommand dataModel.IGenericCommandCreateUpdateDelete[gormModel.User]
	query     query
}

func NewCommand(db *gorm.DB) command {
	return command{
		dbCommand: dataModel.NewGenericCommand[gormModel.User](db),
		query:     NewQuery(db),
	}
}

type CreateReq struct {
	Username string  `validate:"required,min=2,max=32"`
	Email    string  `validate:"required,email"`
	Password string  `validate:"required,min=2,max=32"`
	Fullname *string `validate:"omitempty,max=64"`
	Status   *string `validate:"omitempty,max=128"`
}

func (c command) Create(ctx context.Context, req CreateReq) (*gormModel.User, error) {
	ctx, span := otel.Trace(ctx)
	defer span.End()

	if err := validate.Struct(req); err != nil {
		return nil, err
	}

	user := &gormModel.User{
		ID:             uuid.New().String(),
		Username:       req.Username,
		Email:          req.Email,
		HashedPassword: userPassword.Hash(ctx, req.Password),
		Fullname:       req.Fullname,
		Status:         req.Status,
		CreatedAt:      time.Now().Truncate(time.Millisecond).Local(),
	}

	if err := c.dbCommand.Create(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

type UpdateReq struct {
	ID       string
	Password *string `validate:"omitempty,min=2,max=32"`
	Fullname *string `validate:"omitempty,max=64"`
	Status   *string `validate:"omitempty,max=128"`
}

func (c command) Update(ctx context.Context, req UpdateReq) (*gormModel.User, error) {
	ctx, span := otel.Trace(ctx)
	defer span.End()

	if err := validate.Struct(req); err != nil {
		return nil, err
	}

	user, err := c.query.GetByID(ctx, req.ID)
	if err != nil {
		return nil, err
	}

	timeNow := time.Now().Truncate(time.Millisecond).Local()
	user.UpdatedAt = &timeNow

	if req.Password != nil {
		span.AddEvent("Hashing Password Start")
		user.HashedPassword = userPassword.Hash(ctx, *req.Password)
		span.AddEvent("Hashing Password End")
	}

	if req.Fullname != nil {
		user.Fullname = req.Fullname
	}

	if req.Status != nil {
		user.Status = req.Status
	}

	if err := c.dbCommand.Update(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

type DeleteReq struct {
	ID       string
	Password string `validate:"required"`
}

func (c command) Delete(ctx context.Context, req DeleteReq) error {
	ctx, span := otel.Trace(ctx)
	defer span.End()

	if err := validate.Struct(req); err != nil {
		return err
	}

	user, err := c.query.GetByID(ctx, req.ID)
	if err != nil {
		return err
	}

	span.AddEvent("Checking Password Start")

	if !userPassword.IsReal(ctx, req.Password, user.HashedPassword) {
		return userPassword.ErrInvalidPass
	}

	span.AddEvent("Checking Password End")

	return c.dbCommand.Delete(ctx, &gormModel.User{ID: req.ID})
}
