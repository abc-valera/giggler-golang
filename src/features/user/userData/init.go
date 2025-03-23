package userData

import (
	"gorm.io/gorm"

	"giggler-golang/src/shared/data"
	"giggler-golang/src/shared/serviceLocator"
)

var getDB = serviceLocator.Get[func(ds data.IDB) *gorm.DB]()
