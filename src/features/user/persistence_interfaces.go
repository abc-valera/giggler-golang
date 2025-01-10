package user

import (
	"context"
	"time"

	"github.com/abc-valera/giggler-golang/gen/gorm/gormModel"
	"github.com/abc-valera/giggler-golang/gen/gorm/gormQuery"
	"github.com/abc-valera/giggler-golang/src/components/env"
	"github.com/abc-valera/giggler-golang/src/shared/persistence"
	"gorm.io/gorm"
)

var (
	command func(persistence.IDS) iCommand
	Query   func(persistence.IDS) iQuery
)

func init() {
	switch persistence.DbVal {
	case persistence.DbVariantGorm:
		command = func(dataStore persistence.IDS) iCommand {
			return persistence.NewGormGenericCommand(persistence.GormDS(dataStore), InfraGormNewDTO)
		}
		Query = func(dataStore persistence.IDS) iQuery {
			return &gormQ{gormQuery.Use(persistence.GormDS(dataStore))}
		}
	default:
		panic(env.ErrInvalidEnvValue)
	}
}

type (
	User struct {
		ID             string
		Username       string
		Email          string
		HashedPassword string
		Fullname       string
		Status         string
		CreatedAt      time.Time
		UpdatedAt      time.Time
		DeletedAt      time.Time
	}

	iCommand persistence.IGenericCommandCreateUpdateDelete[User]

	iQuery interface {
		GetByID(ctx context.Context, id string) (User, error)
		GetByEmail(ctx context.Context, email string) (User, error)
		GetAll(ctx context.Context, s persistence.Selector) ([]User, error)
	}
)

type gormQ struct{ *gormQuery.Query }

func (q gormQ) GetByID(ctx context.Context, id string) (User, error) {
	dto, err := q.WithContext(ctx).User.Where(q.User.ID.Eq(id)).First()
	if err != nil {
		return User{}, err
	}
	return InfraGormNewUser(dto), nil
}

func (q gormQ) GetByEmail(ctx context.Context, email string) (User, error) {
	dto, err := q.WithContext(ctx).User.Where(q.User.Email.Eq(email)).First()
	if err != nil {
		return User{}, err
	}
	return InfraGormNewUser(dto), nil
}

func (q gormQ) GetAll(ctx context.Context, s persistence.Selector) ([]User, error) {
	dtos, err := q.WithContext(ctx).User.Limit(int(s.PagingLimit)).Offset(int(s.PagingOffset)).Find()
	if err != nil {
		return nil, err
	}
	return InfraGormNewUsers(dtos), nil
}

func InfraGormNewDTO(user User) *gormModel.User {
	return &gormModel.User{
		ID:             user.ID,
		Username:       user.Username,
		Email:          user.Email,
		HashedPassword: user.HashedPassword,
		Fullname:       user.Fullname,
		Status:         user.Status,
		CreatedAt:      user.CreatedAt,
		UpdatedAt:      user.UpdatedAt,
		DeletedAt:      gorm.DeletedAt{Time: user.DeletedAt, Valid: true},
	}
}

func InfraGormNewUser(dto *gormModel.User) User {
	return User{
		ID:             dto.ID,
		Username:       dto.Username,
		Email:          dto.Email,
		HashedPassword: dto.HashedPassword,
		Fullname:       dto.Fullname,
		Status:         dto.Status,
		CreatedAt:      dto.CreatedAt,
		UpdatedAt:      dto.UpdatedAt,
		DeletedAt:      dto.DeletedAt.Time,
	}
}

func InfaGormNewDTOs(users []User) []*gormModel.User {
	var dtos []*gormModel.User
	for _, user := range users {
		dtos = append(dtos, InfraGormNewDTO(user))
	}
	return dtos
}

func InfraGormNewUsers(dtos []*gormModel.User) []User {
	var users []User
	for _, dto := range dtos {
		users = append(users, InfraGormNewUser(dto))
	}
	return users
}
