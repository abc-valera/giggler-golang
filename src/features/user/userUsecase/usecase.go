package userUsecase

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"giggler-golang/src/core/must"
	"giggler-golang/src/features/user/jwt"
	"giggler-golang/src/features/user/userData"
	"giggler-golang/src/features/user/userPassword"
	"giggler-golang/src/infra/db/dbDto"
	"giggler-golang/src/infra/db/gormgenQuery"
	"giggler-golang/src/infra/emailer"
)

type Usecase struct {
	db      *gorm.DB
	emailer emailer.Interface
}

func New(db *gorm.DB, emailer emailer.Interface) Usecase {
	return must.NoNilInterfaceFields(Usecase{
		db:      db,
		emailer: emailer,
	})
}

type RegisterIn struct {
	Username string
	Email    string
	Password string
}

func (u Usecase) Register(ctx context.Context, req RegisterIn) error {
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

		if res := tx.WithContext(ctx).Create(user); res != nil {
			return dbDto.CommandError(res)
		}

		return u.emailer.Send(emailer.SendIn{
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

type LoginIn struct {
	Email    string
	Password string
}

type LoginOut struct {
	User         *userData.User
	AccessToken  string
	RefreshToken string
}

func (u Usecase) Login(ctx context.Context, req LoginIn) (LoginOut, error) {
	q := gormgenQuery.Use(u.db)
	foundUser, err := q.WithContext(ctx).User.Where(q.User.Email.Eq(req.Email)).First()
	if err != nil {
		return LoginOut{}, err
	}

	if !userPassword.IsReal(ctx, req.Password, foundUser.HashedPassword) {
		return LoginOut{}, userPassword.ErrInvalidPass
	}

	payload := jwt.Payload{
		UserID: foundUser.ID,
	}

	return LoginOut{
		User:         foundUser,
		AccessToken:  jwt.NewAccessToken(payload),
		RefreshToken: jwt.NewRefreshToken(payload),
	}, nil
}

func (u Usecase) Refresh(ctx context.Context, refreshToken string) (string, error) {
	payload, err := jwt.VerifyRefreshToken(refreshToken)
	if err != nil {
		return "", err
	}

	return jwt.NewAccessToken(jwt.Payload{
		UserID: payload.UserID,
	}), nil
}

func (u Usecase) Profile(ctx context.Context) (*userData.User, error) {
	p := jwt.Get(ctx)

	q := gormgenQuery.Use(u.db)
	user, err := q.WithContext(ctx).User.Where(q.User.ID.Eq(p.UserID)).First()
	if err != nil {
		return nil, err
	}

	return user, nil
}
