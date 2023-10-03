package data

import (
	"errors"

	"github.com/abc-valera/giggler-golang/src/components/errutil"
	"gorm.io/gorm"
)

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
