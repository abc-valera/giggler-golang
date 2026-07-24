package dbDto

import (
	"errors"

	"gorm.io/gorm"

	"giggler-golang/src/shared/errutil"
)

func CommandError(res *gorm.DB) error {
	err := res.Error

	if err == nil {
		return nil
	}

	if errors.Is(err, gorm.ErrCheckConstraintViolated) {
		return errutil.Wrap(err, errutil.ErrValidation)
	}

	if errors.Is(err, gorm.ErrDuplicatedKey) {
		return errutil.Wrap(err, errutil.ErrConflict)
	}

	return errutil.Wrap(err, errutil.ErrInternalServer)
}

func QueryError(err error) error {
	if err == nil {
		return nil
	}

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return errutil.Wrap(err, errutil.ErrNotFound)
	}

	return errutil.Wrap(err, errutil.ErrInternalServer)
}
