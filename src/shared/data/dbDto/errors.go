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
		return errutil.NewCode(errutil.CodeInvalidArgument, err)
	}

	if errors.Is(err, gorm.ErrDuplicatedKey) {
		return errutil.NewCode(errutil.CodeAlreadyExists, err)
	}

	return errutil.NewInternal(err)
}

func QueryError(err error) error {
	if err == nil {
		return nil
	}

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return errutil.NewCode(errutil.CodeNotFound, err)
	}

	return errutil.NewInternal(err)
}
