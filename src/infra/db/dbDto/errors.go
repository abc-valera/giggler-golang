package dbDto

import (
	"errors"

	"gorm.io/gorm"

	"giggler-golang/src/core/errutil"
)

func CommandError(res *gorm.DB) error {
	err := res.Error

	if err == nil {
		return nil
	}

	if errors.Is(err, gorm.ErrCheckConstraintViolated) {
		return errutil.Wrap(err, errutil.ErrorValidation)
	}

	if errors.Is(err, gorm.ErrDuplicatedKey) {
		return errutil.Wrap(err, errutil.ErrorConflict)
	}

	return errutil.Wrap(err, errutil.ErrorInternalServer)
}

func QueryError(err error) error {
	if err == nil {
		return nil
	}

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return errutil.Wrap(err, errutil.ErrorNotFound)
	}

	return errutil.Wrap(err, errutil.ErrorInternalServer)
}
