package userData

import (
	"reflect"
	"time"

	"giggler-golang/src/shared/data"
	"giggler-golang/src/shared/validate"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func init() {
	data.GetDb().AutoMigrate(&User{})
}

// TODO: check if 'required' keyword is needed here
// (I suppose it can be enabled by default)
type User struct {
	ID             uuid.UUID
	Username       string  `validate:"required,min=3,max=32" gorm:"unique;not null"`
	Email          string  `validate:"required,email" gorm:"unique;not null"`
	HashedPassword string  `validate:"required" gorm:"not null"`
	Fullname       *string `validate:"omitempty,max=64" gorm:""`
	Status         *string `validate:"omitempty,max=128" gorm:""`
	CreatedAt      time.Time
	UpdatedAt      *time.Time
	DeletedAt      gorm.DeletedAt
}

// TODO: create a generic func from this,
// allow to pass additional model-related validation logic
func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.ID == uuid.Nil {
		u.ID = uuid.New()
	}
	if u.CreatedAt.IsZero() {
		u.CreatedAt = time.Now()
	}

	if err := validate.Struct(u); err != nil {
		return err
	}

	return nil
}

// TODO: create a generic func from this,
// allow to pass additional model-related validation logic
func (u *User) BeforeUpdate(tx *gorm.DB) error {
	if !tx.Statement.Changed() {
		return nil
	}

	now := time.Now()
	u.UpdatedAt = &now

	t := reflect.TypeOf(*u)

	var changedFields []string
	for i := range t.NumField() {
		if fieldName := t.Field(i).Name; tx.Statement.Changed(fieldName) {
			changedFields = append(changedFields, fieldName)
		}
	}

	if err := validate.StructPartial(u, changedFields...); err != nil {
		return err
	}

	return nil
}
