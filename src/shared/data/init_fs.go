package data

import (
	"fmt"
	"os"

	"giggler-golang/src/shared/data/internal/localFS"
	"giggler-golang/src/shared/env"
)

var FS = func() func() iFS {
	var fs iFS

	switch env.Load("FS") {
	case "local":
		folderPath := env.Load("LOCAL_DSN")

		if err := os.MkdirAll(folderPath, 0o755); err != nil {
			if !os.IsExist(err) {
				panic(fmt.Errorf("failed to create local file saver directory: %w", err))
			}
		}

		fs = localFS.New(folderPath)
	default:
		panic(env.ErrInvalidEnvValue)
	}

	return func() iFS { return fs }
}()

type iFS interface {
	Create(filename string) error
	Read(filename string) ([]byte, error)
	Delete(filename string) error
}
