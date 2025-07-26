package internal

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"giggler-golang/src/features/user/userData"
	"giggler-golang/src/features/user/userPassword"
	"giggler-golang/src/shared/data"
	"giggler-golang/src/shared/data/dbDto"
	"giggler-golang/src/shared/email"
	"giggler-golang/src/shared/otel"
)

type RegisterIn struct {
	Username string
	Email    string
	Password string
}

func Register(ctx context.Context, req RegisterIn) error {
	ctx, span := otel.Trace(ctx)
	defer span.End()

	txFunc := func(tx *gorm.DB) error {
		hash, err := userPassword.Hash(ctx, req.Password)
		if err != nil {
			return err
		}

		user := &userData.User{
			ID:             uuid.New(),
			Username:       req.Username,
			Email:          req.Email,
			HashedPassword: hash,
		}

		if res := data.DB().WithContext(ctx).Create(user); res != nil {
			return dbDto.CommandError(res)
		}

		return email.Send(email.EmailSendIn{
			Subject: "Verification Email for Giggler!",
			Content: fmt.Sprintf("%s, congrats with joining the Giggler community!", req.Username),
			To:      []string{req.Email},
		})
	}

	if err := data.DB().WithContext(ctx).Transaction(txFunc); err != nil {
		// TODO: make sure process this error properly
		// TODO: maybe make gorm errors to be domain errors
		return fmt.Errorf("failed to register user: %w", err)
	}

	return nil
}
