package authLogic

import (
	"context"
	"fmt"

	"github.com/abc-valera/giggler-golang/src/components/passworder"
	"github.com/abc-valera/giggler-golang/src/features/user/userData"
	"github.com/abc-valera/giggler-golang/src/features/user/userModel"
	"github.com/abc-valera/giggler-golang/src/shared/data"
	"github.com/abc-valera/giggler-golang/src/shared/otel"
	"go.opentelemetry.io/otel/trace"
)

type AuthRegisterReq struct {
	Username string
	Email    string
	Password string
}

var welcomeEmailTemplateFunc = func(username, email string) emailSendReq {
	return emailSendReq{
		Subject: "Verification Email for Giggler!",
		Content: fmt.Sprintf("%s, congrats with joining the Giggler community!", username),
		To:      []string{email},
	}
}

// AuthRegister performs user sign-up:
//   - it creates new user entity with unique username and email,
//   - creates hash of the password provided by user,
//   - then it sends welcome email to the users's email address,
func AuthRegister(ctx context.Context, req AuthRegisterReq) (userModel.User, error) {
	var span trace.Span
	ctx, span = otel.Trace(ctx)
	defer span.End()

	var user userModel.User
	txFunc := func(ctx context.Context, ds data.IDS) error {
		var err error
		user, err = userData.Create(ctx, ds, userData.CreateReq{
			Username: req.Username,
			Email:    req.Email,
			Password: req.Password,
		})
		if err != nil {
			return err
		}

		return emailer.SendEmail(welcomeEmailTemplateFunc(req.Username, req.Email))
	}

	if err := data.DB.WithinTX(ctx, txFunc); err != nil {
		return userModel.User{}, err
	}

	return user, nil
}

type AuthLoginReq struct {
	Email    string
	Password string
}

// AuthLogin performs user sign-in: it checks if user with provided email exists,
// then creates hash of the provided password and compares it to the hash stored in database.
// The AuthLogin returns user, accessToken and refreshToken.
func AuthLogin(ctx context.Context, req AuthLoginReq) (userModel.User, error) {
	var span trace.Span
	ctx, span = otel.Trace(ctx)
	defer span.End()

	foundUser, err := userData.Query(data.DB).GetByEmail(ctx, req.Email)
	if err != nil {
		return userModel.User{}, err
	}

	if err := passworder.Check(ctx, req.Password, foundUser.HashedPassword); err != nil {
		return userModel.User{}, err
	}

	return foundUser, nil
}
