package userView

import (
	"giggler-golang/src/shared/data/dbgen/gormModel"
	"giggler-golang/src/shared/view/viewDTO"
	"giggler-golang/src/shared/view/viewgen"
)

func NewUserModel(user *gormModel.User) *viewgen.UserSchema {
	return &viewgen.UserSchema{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		Fullname:  viewDTO.NewOptString(user.Fullname),
		Status:    viewDTO.NewOptString(user.Fullname),
		CreatedAt: user.CreatedAt,
	}
}
