// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package gormModel

import (
	"time"

	"gorm.io/gorm"
)

const TableNameUser = "users"

// User mapped from table <users>
type User struct {
	ID             string         `gorm:"column:id;primaryKey" json:"id"`
	Username       string         `gorm:"column:username;not null" json:"username"`
	Email          string         `gorm:"column:email;not null" json:"email"`
	HashedPassword string         `gorm:"column:hashed_password;not null" json:"hashed_password"`
	Fullname       *string        `gorm:"column:fullname" json:"fullname"`
	Status         *string        `gorm:"column:status" json:"status"`
	CreatedAt      time.Time      `gorm:"column:created_at;not null" json:"created_at"`
	UpdatedAt      *time.Time     `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"column:deleted_at" json:"deleted_at"`
}

// TableName User's table name
func (*User) TableName() string {
	return TableNameUser
}
