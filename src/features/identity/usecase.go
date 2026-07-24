package identity

import (
	"context"
	"fmt"

	"giggler-golang/src/features/identity/internal"
	"giggler-golang/src/shared/db/dbDto"
	"giggler-golang/src/shared/emailer"

	"github.com/google/uuid"
	"gorm.io/cli/gorm/typed"
	"gorm.io/gorm"
)

type Usecase struct {
	db               *gorm.DB
	emailer          emailer.Interface
	authTokenService AuthTokenService
}

func New(db *gorm.DB, emailer emailer.Interface) Usecase {
	return Usecase{
		db:               db,
		emailer:          emailer,
		authTokenService: initAuthTokenService(),
	}
}

type RegisterInput struct {
	Username string
	Email    string
	Password string
}

func (u Usecase) Register(ctx context.Context, req RegisterInput) error {
	txFunc := func(tx *gorm.DB) error {
		hash, err := internal.NewPasswordHash(ctx, req.Password)
		if err != nil {
			return err
		}

		user := &User{
			ID:             uuid.New(),
			Username:       req.Username,
			Email:          req.Email,
			HashedPassword: hash,
		}

		if res := tx.WithContext(ctx).Create(user); res != nil {
			return dbDto.CommandError(res)
		}

		return u.emailer.Send(emailer.EmailData{
			Subject: "Verification Email for Giggler!",
			Content: fmt.Sprintf("%s, congrats with joining the Giggler community!", req.Username),
			To:      []string{req.Email},
		})
	}

	if err := u.db.WithContext(ctx).Transaction(txFunc); err != nil {
		// TODO: make sure process this error properly
		// TODO: maybe make gorm errors to be domain errors
		// TODO: check if the txFunc panics, then everything is rolled back automatically
		return fmt.Errorf("failed to register user: %w", err)
	}

	return nil
}

type LoginInput struct {
	Email    string
	Password string
}

type LoginOutput struct {
	User            User
	AccessTokenJWT  string
	RefreshTokenJWT string
}

func (u Usecase) Login(ctx context.Context, req LoginInput) (LoginOutput, error) {
	foundUser, err := typed.G[User](u.db).
		Where(UserHelper.Email.Eq(req.Email)).
		Take(ctx)
	if err != nil {
		return LoginOutput{}, err
	}

	if err := internal.PasswordVerify(ctx, req.Password, foundUser.HashedPassword); err != nil {
		return LoginOutput{}, err
	}

	return LoginOutput{
		User:            foundUser,
		AccessTokenJWT:  u.authTokenService.NewAccessTokenJWT(foundUser.ID),
		RefreshTokenJWT: u.authTokenService.NewRefreshTokenJWT(foundUser.ID),
	}, nil
}

func (u Usecase) Refresh(ctx context.Context, refreshToken string) (string, error) {
	payload, err := u.authTokenService.VerifyRefreshTokenJWT(refreshToken)
	if err != nil {
		return "", err
	}

	return u.authTokenService.NewAccessTokenJWT(payload.UserID), nil
}

func (u Usecase) Profile(ctx context.Context) (User, error) {
	token := GetAuthToken(ctx)

	user, err := typed.G[User](u.db).Where(UserHelper.ID.Eq(token.UserID)).First(ctx)
	if err != nil {
		return User{}, err
	}

	return user, nil
}
