package dbDTO

import (
	"time"

	"gorm.io/gorm"
)

func NewDeletedAt(val *time.Time) gorm.DeletedAt {
	if val == nil {
		return gorm.DeletedAt{}
	}
	return gorm.DeletedAt{Valid: true, Time: *val}
}

func NewDomainDeletedAt(val gorm.DeletedAt) *time.Time {
	if !val.Valid {
		return nil
	}
	return &val.Time
}
