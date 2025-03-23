package data

import (
	"context"
	"fmt"
	"os"

	"giggler-golang/src/shared/dependency"
	"giggler-golang/src/shared/env"
)

var (
	FsVariantLocal = "local"

	FsVal = env.Load("FS")
	FS    = initFS()
)

func initFS() IDS {
	switch FsVal {
	case FsVariantLocal:
		fsLocalPath := env.Load("LOCAL_DSN")

		if err := os.MkdirAll(fsLocalPath, 0o755); err != nil {
			if !os.IsExist(err) {
				panic(fmt.Errorf("failed to create local file saver directory: %w", err))
			}
		}

		return fsLocal{dependency: fsLocalPath}
	default:
		panic(env.ErrInvalidEnvValue)
	}
}

type fsLocal struct{ dependency string }

var _ dependency.Interface[string] = fsLocal{}

func (fsLocal) BeginTX(ctx context.Context) (ITX, error) {
	panic("unimplemented")
}

func (fsLocal) WithinTX(ctx context.Context, fn func(context.Context, IDS) error) error {
	panic("unimplemented")
}

func (fs fsLocal) Dependency() string {
	return fs.dependency
}

type fsLocalTX struct{ dependency string }

var _ dependency.Interface[string] = fsLocalTX{}

func (tx fsLocalTX) Dependency() string {
	return tx.dependency
}
