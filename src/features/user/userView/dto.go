package userView

import (
	"time"

	"giggler-golang/src/features/user/userData"

	"github.com/google/uuid"
)

type UserModel struct {
	ID        uuid.UUID
	Username  string
	Email     string
	Fullname  *string
	Status    *string
	CreatedAt time.Time
}

func NewUserDTO(user *userData.User) *UserModel {
	return &UserModel{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		Fullname:  user.Fullname,
		Status:    user.Status,
		CreatedAt: user.CreatedAt,
	}
}
