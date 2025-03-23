package userDTO

import (
	"giggler-golang/src/features/user/userModel"
	"giggler-golang/src/shared/view/viewDto"
	"giggler-golang/src/shared/view/viewgen"
)

func NewView(user userModel.User) *viewgen.UserSchema {
	return &viewgen.UserSchema{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		Fullname:  viewDto.NewOptString(user.Fullname),
		Status:    viewDto.NewOptString(user.Fullname),
		CreatedAt: user.CreatedAt,
	}
}
