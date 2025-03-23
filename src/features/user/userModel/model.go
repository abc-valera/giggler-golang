package userModel

import "time"

type User struct {
	ID             string
	Username       string
	Email          string
	HashedPassword string
	Fullname       *string
	Status         *string
	CreatedAt      time.Time
	UpdatedAt      *time.Time
	DeletedAt      *time.Time
}
