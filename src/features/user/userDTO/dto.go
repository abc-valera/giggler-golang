package userDTO

import (
	"giggler-golang/src/features/user/userModel"
	"giggler-golang/src/shared/data/dbDTO"
	"giggler-golang/src/shared/data/dbgen/gormModel"
)

func NewGorm(user userModel.User) *gormModel.User {
	return &gormModel.User{
		ID:             user.ID,
		Username:       user.Username,
		Email:          user.Email,
		HashedPassword: user.HashedPassword,
		Fullname:       user.Fullname,
		Status:         user.Status,
		CreatedAt:      user.CreatedAt,
		UpdatedAt:      user.UpdatedAt,
		DeletedAt:      dbDTO.NewDeletedAt(user.DeletedAt),
	}
}

func FromGorm(dto *gormModel.User) userModel.User {
	return userModel.User{
		ID:             dto.ID,
		Username:       dto.Username,
		Email:          dto.Email,
		HashedPassword: dto.HashedPassword,
		Fullname:       dto.Fullname,
		Status:         dto.Status,
		CreatedAt:      dto.CreatedAt,
		UpdatedAt:      dto.UpdatedAt,
		DeletedAt:      dbDTO.NewDomainDeletedAt(dto.DeletedAt),
	}
}

func NewGorms(users []userModel.User) []*gormModel.User {
	var dtos []*gormModel.User
	for _, user := range users {
		dtos = append(dtos, NewGorm(user))
	}
	return dtos
}

func FromGorms(dtos []*gormModel.User) []userModel.User {
	var users []userModel.User
	for _, dto := range dtos {
		users = append(users, FromGorm(dto))
	}
	return users
}
