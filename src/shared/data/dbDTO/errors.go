package dbDTO

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
		return errutil.NewCodeError(errutil.CodeInvalidArgument, err)
	}

	if errors.Is(err, gorm.ErrDuplicatedKey) {
		return errutil.NewCodeError(errutil.CodeAlreadyExists, err)
	}

	return errutil.NewInternalErr(err)
}

func QueryError(err error) error {
	if err == nil {
		return nil
	}

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return errutil.NewCodeError(errutil.CodeNotFound, err)
	}

	return errutil.NewInternalErr(err)
}
