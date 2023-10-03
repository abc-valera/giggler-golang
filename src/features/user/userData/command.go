package userData

import (
	"context"
	"time"

	"github.com/abc-valera/giggler-golang/src/components/env"
	"github.com/abc-valera/giggler-golang/src/components/passworder"
	"github.com/abc-valera/giggler-golang/src/features/user/userModel"
	"github.com/abc-valera/giggler-golang/src/shared/data"
	"github.com/abc-valera/giggler-golang/src/shared/otel"
	"github.com/abc-valera/giggler-golang/src/shared/validate"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel/trace"
)

var command = func() func(data.IDS) iCommand {
	switch data.DbVal {
	case data.DbVariantGorm:
		return func(dataStore data.IDS) iCommand {
			return data.NewGormGenericCommand(data.GormDS(dataStore), userModel.NewGormDTO)
		}
	default:
		panic(env.ErrInvalidEnvValue)
	}
}()

type iCommand data.IGenericCommandCreateUpdateDelete[userModel.User]

type CreateReq struct {
	Username string  `validate:"required,min=2,max=32"`
	Email    string  `validate:"required,email"`
	Password string  `validate:"required,min=2,max=32"`
	Fullname *string `validate:"omitempty,max=64"`
	Status   *string `validate:"omitempty,max=128"`
}

func Create(ctx context.Context, ds data.IDS, req CreateReq) (userModel.User, error) {
	var span trace.Span
	ctx, span = otel.Trace(ctx)
	defer span.End()

	if err := validate.Struct(req); err != nil {
		return userModel.User{}, err
	}

	hashedPassword, err := passworder.Hash(ctx, req.Password)
	if err != nil {
		return userModel.User{}, err
	}

	user := userModel.User{
		ID:             uuid.New().String(),
		Username:       req.Username,
		Email:          req.Email,
		HashedPassword: hashedPassword,
		Fullname:       req.Fullname,
		Status:         req.Status,
		CreatedAt:      time.Now().Truncate(time.Millisecond).Local(),
	}

	if err := command(ds).Create(ctx, user); err != nil {
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

func Update(ctx context.Context, ds data.IDS, req UpdateReq) (userModel.User, error) {
	var span trace.Span
	ctx, span = otel.Trace(ctx)
	defer span.End()

	if err := validate.Struct(req); err != nil {
		return userModel.User{}, err
	}

	user, err := Query(ds).GetByID(ctx, req.ID)
	if err != nil {
		return userModel.User{}, err
	}

	timeNow := time.Now().Truncate(time.Millisecond).Local()
	user.UpdatedAt = &timeNow

	if req.Password != nil {
		span.AddEvent("Hashing Password Start")

		hashedPassword, err := passworder.Hash(ctx, *req.Password)
		if err != nil {
			return userModel.User{}, err
		}
		user.HashedPassword = hashedPassword

		span.AddEvent("Hashing Password End")
	}

	if req.Fullname != nil {
		user.Fullname = req.Fullname
	}

	if req.Status != nil {
		user.Status = req.Status
	}

	if err := command(ds).Update(ctx, user); err != nil {
		return userModel.User{}, err
	}

	return user, nil
}

type DeleteReq struct {
	ID       string
	Password string `validate:"required"`
}

func Delete(ctx context.Context, ds data.IDS, req DeleteReq) error {
	var span trace.Span
	ctx, span = otel.Trace(ctx)
	defer span.End()

	if err := validate.Struct(req); err != nil {
		return err
	}

	user, err := Query(ds).GetByID(ctx, req.ID)
	if err != nil {
		return err
	}

	span.AddEvent("Checking Password Start")

	if err := passworder.Check(ctx, req.Password, user.HashedPassword); err != nil {
		return err
	}

	span.AddEvent("Checking Password End")

	return command(ds).Delete(ctx, userModel.User{ID: req.ID})
}
