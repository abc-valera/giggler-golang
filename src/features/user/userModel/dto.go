package userModel

import (
	"github.com/abc-valera/giggler-golang/src/shared/data"
	"github.com/abc-valera/giggler-golang/src/shared/data/gormgen/gormModel"
	"github.com/abc-valera/giggler-golang/src/shared/view"
	"github.com/abc-valera/giggler-golang/src/shared/view/viewgen"
)

func NewGormDTO(user User) *gormModel.User {
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

func FromGormDTO(dto *gormModel.User) User {
	return User{
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

func NewGormDTOs(users []User) []*gormModel.User {
	var dtos []*gormModel.User
	for _, user := range users {
		dtos = append(dtos, NewGormDTO(user))
	}
	return dtos
}

func FromGormDTOs(dtos []*gormModel.User) []User {
	var users []User
	for _, dto := range dtos {
		users = append(users, FromGormDTO(dto))
	}
	return users
}

func NewViewDTO(user User) *viewgen.UserSchema {
	return &viewgen.UserSchema{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		Fullname:  view.NewOptString(user.Fullname),
		Status:    view.NewOptString(user.Fullname),
		CreatedAt: user.CreatedAt,
	}
}
