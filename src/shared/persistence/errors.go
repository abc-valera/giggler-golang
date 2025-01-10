package persistence

import (
	"github.com/abc-valera/giggler-golang/src/components/errutil"
	"gorm.io/gorm"
)

func GormCommandError(res *gorm.DB) error {
	return errutil.NewInternalErr(res.Error)
}
