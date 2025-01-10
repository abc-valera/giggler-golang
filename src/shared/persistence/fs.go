package persistence

import (
	"context"
	"fmt"
	"os"

	"github.com/abc-valera/giggler-golang/src/components/dependency"
	"github.com/abc-valera/giggler-golang/src/components/env"
)

var (
	FS             IDS
	FsVal          string
	FsVariantLocal = "local"
)

func init() {
	switch FsVal = env.Load("FS"); FsVal {
	case FsVariantLocal:
		fsLocalPath := env.Load("FS_LOCAL_PATH")

		if err := os.MkdirAll(fsLocalPath, 0o755); err != nil {
			if !os.IsExist(err) {
				panic(fmt.Errorf("failed to create local file saver directory: %w", err))
			}
		}

		FS = fsLocal{dependency: fsLocalPath}
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
