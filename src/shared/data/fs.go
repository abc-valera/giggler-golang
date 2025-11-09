package data

import (
	"fmt"
	"os"

	"giggler-golang/src/shared/data/internal/localFS"
	"giggler-golang/src/shared/errutil/must"
	"giggler-golang/src/shared/singleton"
)

var GetFs = singleton.New(func() fileSystem {
	switch must.GetEnv("FS") {
	case "local":
		folderPath := must.GetEnv("LOCAL_DSN")

		if err := os.MkdirAll(folderPath, 0o755); err != nil {
			if !os.IsExist(err) {
				panic(fmt.Errorf("failed to create local file saver directory: %w", err))
			}
		}

		return localFS.New(folderPath)
	default:
		panic(must.ErrInvalidEnvValue)
	}
})

type fileSystem interface {
	Create(filename string) error
	Read(filename string) ([]byte, error)
	Delete(filename string) error
}
