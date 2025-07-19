package jokeData

import (
	"time"

	"gorm.io/gorm"
)

type Joke struct {
	ID          string         `gorm:"column:primaryKey"`
	Title       string         `gorm:"column:not null;uniqueIndex:idx_user_joke_title"`
	Text        string         `gorm:"column:not null"`
	Explanation *string        `gorm:"column:"`
	CreatedAt   time.Time      `gorm:"column:not null"`
	UpdatedAt   *time.Time     `gorm:"column:"`
	DeletedAt   gorm.DeletedAt `gorm:"column:"`

	UserID string `gorm:"column:not null;uniqueIndex:idx_user_joke_title"`
}
