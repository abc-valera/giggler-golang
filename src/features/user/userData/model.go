package userData

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID             string         `gorm:"column:primaryKey"`
	Username       string         `gorm:"column:not null"`
	Email          string         `gorm:"column:not null"`
	HashedPassword string         `gorm:"column:not null"`
	Fullname       *string        `gorm:"column:"`
	Status         *string        `gorm:"column:"`
	CreatedAt      time.Time      `gorm:"column:not null"`
	UpdatedAt      *time.Time     `gorm:"column:"`
	DeletedAt      gorm.DeletedAt `gorm:"column:"`
}
