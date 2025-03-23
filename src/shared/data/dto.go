package data

import (
	"errors"
	"time"

	"github.com/abc-valera/giggler-golang/src/shared/errutil"
	"gorm.io/gorm"
)

func NewDtoGormDeletedAt(val *time.Time) gorm.DeletedAt {
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

func GormCommandError(res *gorm.DB) error {
	err := res.Error

	if err == nil {
		return nil
	}

	if errors.Is(err, gorm.ErrCheckConstraintViolated) {
		return errutil.NewCodeError(errutil.CodeInvalidArgument, err)
	}

	if errors.Is(err, gorm.ErrDuplicatedKey) {
		return errutil.NewCodeError(errutil.CodeAlreadyExists, err)
	}

	return errutil.NewInternalErr(err)
}

func GormQueryError(err error) error {
	if err == nil {
		return nil
	}

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return errutil.NewCodeError(errutil.CodeNotFound, err)
	}

	return errutil.NewInternalErr(err)
}
