package userDto

import (
	"giggler-golang/src/features/user/userModel"
	"giggler-golang/src/shared/data"
	"giggler-golang/src/shared/data/gormgen/gormModel"
	"giggler-golang/src/shared/view/viewUtil"
	"giggler-golang/src/shared/view/viewgen"
)

func NewGormDTO(user userModel.User) *gormModel.User {
	return &gormModel.User{
		ID:             user.ID,
		Username:       user.Username,
		Email:          user.Email,
		HashedPassword: user.HashedPassword,
		Fullname:       user.Fullname,
		Status:         user.Status,
		CreatedAt:      user.CreatedAt,
		UpdatedAt:      user.UpdatedAt,
		DeletedAt:      data.NewDtoGormDeletedAt(user.DeletedAt),
	}
}

func FromGormDTO(dto *gormModel.User) userModel.User {
	return userModel.User{
		ID:             dto.ID,
		Username:       dto.Username,
		Email:          dto.Email,
		HashedPassword: dto.HashedPassword,
		Fullname:       dto.Fullname,
		Status:         dto.Status,
		CreatedAt:      dto.CreatedAt,
		UpdatedAt:      dto.UpdatedAt,
		DeletedAt:      data.NewDomainDeletedAt(dto.DeletedAt),
	}
}

func NewGormDTOs(users []userModel.User) []*gormModel.User {
	var dtos []*gormModel.User
	for _, user := range users {
		dtos = append(dtos, NewGormDTO(user))
	}
	return dtos
}

func FromGormDTOs(dtos []*gormModel.User) []userModel.User {
	var users []userModel.User
	for _, dto := range dtos {
		users = append(users, FromGormDTO(dto))
	}
	return users
}

func NewViewDTO(user userModel.User) *viewgen.UserSchema {
	return &viewgen.UserSchema{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		Fullname:  viewUtil.NewOptString(user.Fullname),
		Status:    viewUtil.NewOptString(user.Fullname),
		CreatedAt: user.CreatedAt,
	}
}
