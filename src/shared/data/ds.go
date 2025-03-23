package data

import (
	"context"

	"github.com/abc-valera/giggler-golang/src/shared/dependency"
	"gorm.io/gorm"
)

type (
	// IDS is an abstraction for datastore operations.
	IDS interface {
		BeginTX(ctx context.Context) (ITX, error)
		WithinTX(
			ctx context.Context,
			fn func(context.Context, IDS) error,
		) error
	}

	// ITX is an abstraction for datastore transaction.
	ITX interface {
		IDS

		Commit() error
		Rollback()
	}
)

func GormDS(datastore IDS) *gorm.DB {
	return datastore.(dependency.Interface[*gorm.DB]).Dependency()
}
