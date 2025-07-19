package userView

import (
	"giggler-golang/src/features/user/userData"
	"giggler-golang/src/shared/view/viewDTO"
	"giggler-golang/src/shared/view/viewgen"
)

func NewUserModel(user *userData.User) *viewgen.UserSchema {
	return &viewgen.UserSchema{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		Fullname:  viewDTO.NewOptString(user.Fullname),
		Status:    viewDTO.NewOptString(user.Fullname),
		CreatedAt: user.CreatedAt,
	}
}
