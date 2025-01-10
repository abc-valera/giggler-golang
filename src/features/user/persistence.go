package user

import (
	"context"
	"time"

	"github.com/abc-valera/giggler-golang/src/shared/otel"
	"github.com/abc-valera/giggler-golang/src/shared/persistence"
	"github.com/abc-valera/giggler-golang/src/shared/validate"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel/trace"
)

type CreateRequest struct {
	Username string `validate:"required,min=2,max=32"`
	Email    string `validate:"required,email"`
	Password string `validate:"required,min=2,max=32"`
	Fullname string `validate:"max=64"`
	Status   string `validate:"max=128"`
}

func Create(ctx context.Context, ds persistence.IDS, req CreateRequest) (User, error) {
	var span trace.Span
	ctx, span = otel.Trace(ctx)
	defer span.End()

	if err := validate.Struct(req); err != nil {
		return User{}, err
	}

	hashedPassword, err := PasswordHash(ctx, req.Password)
	if err != nil {
		return User{}, err
	}

	user := User{
		ID:             uuid.New().String(),
		Username:       req.Username,
		Email:          req.Email,
		HashedPassword: hashedPassword,
		Fullname:       req.Fullname,
		Status:         req.Status,
		CreatedAt:      time.Now().Truncate(time.Millisecond).Local(),
	}

	if err := command(ds).Create(ctx, user); err != nil {
		return User{}, err
	}

	return user, nil
}

type UpdateRequest struct {
	ID       string
	Password *string `validate:"min=2,max=32"`
	Fullname *string `validate:"max=64"`
	Status   *string `validate:"max=128"`
}

func Update(ctx context.Context, ds persistence.IDS, req UpdateRequest) (User, error) {
	var span trace.Span
	ctx, span = otel.Trace(ctx)
	defer span.End()

	if err := validate.Struct(req); err != nil {
		return User{}, err
	}

	user, err := Query(ds).GetByID(ctx, req.ID)
	if err != nil {
		return User{}, err
	}

	user.UpdatedAt = time.Now().Truncate(time.Millisecond).Local()

	if req.Password != nil {
		span.AddEvent("Hashing Password Start")

		hashedPassword, err := PasswordHash(ctx, *req.Password)
		if err != nil {
			return User{}, err
		}
		user.HashedPassword = hashedPassword

		span.AddEvent("Hashing Password End")
	}

	if req.Fullname != nil {
		user.Fullname = *req.Fullname
	}

	if req.Status != nil {
		user.Status = *req.Status
	}

	if err := command(ds).Update(ctx, user); err != nil {
		return User{}, err
	}

	return user, nil
}

type DeleteRequest struct {
	ID       string
	Password string `validate:"required"`
}

func Delete(ctx context.Context, ds persistence.IDS, req DeleteRequest) error {
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

	if err := PasswordCheck(ctx, req.Password, user.HashedPassword); err != nil {
		return err
	}

	span.AddEvent("Checking Password End")

	return command(ds).Delete(ctx, User{ID: req.ID})
}
