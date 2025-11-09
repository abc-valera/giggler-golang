package jokeData

import (
	"reflect"
	"time"

	"giggler-golang/src/shared/data"
	"giggler-golang/src/shared/validate"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func init() {
	data.GetDb().AutoMigrate(&Joke{})
}

type Joke struct {
	ID          uuid.UUID
	Title       string  `validate:"required,min=4,max=64" gorm:"not null;uniqueIndex:idx_user_id_joke_title"`
	Text        string  `validate:"required,min=4,max=4096" gorm:"not null"`
	Explanation *string `validate:"omitempty,max=4096" gorm:""`
	CreatedAt   time.Time
	UpdatedAt   *time.Time
	DeletedAt   gorm.DeletedAt

	UserID uuid.UUID `validate:"required" gorm:"uniqueIndex:idx_user_id_joke_title"`
}

// TODO: create a generic func from this,
// allow to pass additional model-related validation logic
func (j *Joke) BeforeCreate(tx *gorm.DB) error {
	j.ID = uuid.New()
	j.CreatedAt = time.Now()

	if err := validate.Struct(j); err != nil {
		return err
	}

	return nil
}

// TODO: create a generic func from this,
// allow to pass additional model-related validation logic
func (j *Joke) BeforeUpdate(tx *gorm.DB) error {
	if !tx.Statement.Changed() {
		return nil
	}

	now := time.Now()
	j.UpdatedAt = &now

	t := reflect.TypeOf(*j)

	var changedFields []string
	for i := range t.NumField() {
		if fieldName := t.Field(i).Name; tx.Statement.Changed(fieldName) {
			changedFields = append(changedFields, fieldName)
		}
	}

	if err := validate.StructPartial(j, changedFields...); err != nil {
		return err
	}

	return nil
}
