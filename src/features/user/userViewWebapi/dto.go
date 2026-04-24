package userViewWebapi

import (
	"time"

	"giggler-golang/src/features/user/userData"

	"github.com/google/uuid"
)

type Model struct {
	ID        uuid.UUID
	Username  string
	Email     string
	Fullname  *string
	Status    *string
	CreatedAt time.Time
}

func NewModel(user userData.User) *Model {
	return &Model{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		Fullname:  user.Fullname,
		Status:    user.Status,
		CreatedAt: user.CreatedAt,
	}
}
