package internal

import (
	"context"
	"fmt"

	"gorm.io/gorm"

	"giggler-golang/src/features/user/userData"
	"giggler-golang/src/features/user/userData/userRepo"
	"giggler-golang/src/shared/data"
	"giggler-golang/src/shared/email"
	"giggler-golang/src/shared/otel"
)

type RegisterReq struct {
	Username string
	Email    string
	Password string
}

func Register(ctx context.Context, req RegisterReq) (*userData.User, error) {
	ctx, span := otel.Trace(ctx)
	defer span.End()

	var u *userData.User
	txFunc := func(tx *gorm.DB) error {
		var err error
		u, err = userRepo.NewCommand(tx).Create(ctx, userRepo.CreateReq{
			Username: req.Username,
			Email:    req.Email,
			Password: req.Password,
		})
		if err != nil {
			return err
		}

		return email.Send(email.EmailSendReq{
			Subject: "Verification Email for Giggler!",
			Content: fmt.Sprintf("%s, congrats with joining the Giggler community!", req.Username),
			To:      []string{req.Email},
		})
	}

	if err := data.DB().WithContext(ctx).Transaction(txFunc); err != nil {
		// TODO: make sure process this error properly
		// TODO: maybe make gorm errors to be domain errors
		return nil, fmt.Errorf("failed to register user: %w", err)
	}

	return u, nil
}
