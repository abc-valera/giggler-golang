package auth

import (
	"context"
	"fmt"

	"github.com/abc-valera/giggler-golang/src/features/user"
	"github.com/abc-valera/giggler-golang/src/shared/otel"
	"github.com/abc-valera/giggler-golang/src/shared/persistence"
	"go.opentelemetry.io/otel/trace"
)

type authRegisterRequest struct {
	Username string
	Email    string
	Password string
	Fullname string
	Status   string
}

var welcomeEmailTemplateFunc = func(username, email string) emailSendRequest {
	return emailSendRequest{
		Subject: "Verification Email for Giggler!",
		Content: fmt.Sprintf("%s, congrats with joining the Giggler community!", username),
		To:      []string{email},
	}
}

// authRegister performs user sign-up:
//   - it creates new user entity with unique username and email,
//   - creates hash of the password provided by user,
//   - then it sends welcome email to the users's email address,
func authRegister(ctx context.Context, req authRegisterRequest) (user.User, error) {
	var span trace.Span
	ctx, span = otel.Trace(ctx)
	defer span.End()

	var userModel user.User
	txFunc := func(ctx context.Context, ds persistence.IDS) error {
		var err error
		userModel, err = user.Create(ctx, ds, user.CreateRequest{
			Username: req.Username,
			Email:    req.Email,
			Password: req.Password,
			Fullname: req.Fullname,
			Status:   req.Status,
		})
		if err != nil {
			return err
		}

		return emailer.SendEmail(welcomeEmailTemplateFunc(req.Username, req.Email))
	}

	if err := persistence.DB.WithinTX(ctx, txFunc); err != nil {
		return user.User{}, err
	}

	return userModel, nil
}

type authLoginRequest struct {
	Email    string
	Password string
}

// authLogin performs user sign-in: it checks if user with provided email exists,
// then creates hash of the provided password and compares it to the hash stored in database.
// The authLogin returns user, accessToken and refreshToken.
func authLogin(ctx context.Context, req authLoginRequest) (user.User, error) {
	var span trace.Span
	ctx, span = otel.Trace(ctx)
	defer span.End()

	foundUser, err := user.Query(persistence.DB).GetByEmail(ctx, req.Email)
	if err != nil {
		return user.User{}, err
	}

	if err := user.PasswordCheck(ctx, req.Password, foundUser.HashedPassword); err != nil {
		return user.User{}, err
	}

	return foundUser, nil
}
