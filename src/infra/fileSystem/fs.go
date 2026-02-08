package fileSystem

import (
	"fmt"
	"os"

	"giggler-golang/src/core/must"
)

type Interface interface {
	Create(filename string) error
	Read(filename string) ([]byte, error)
	Delete(filename string) error
}

type localFS struct {
	folderPath string
}

func InitLocal() Interface {
	folderPath := must.GetEnv("LOCAL_DSN")

	if err := os.MkdirAll(folderPath, 0o755); err != nil {
		if !os.IsExist(err) {
			panic(fmt.Errorf("failed to create local file saver directory: %w", err))
		}
	}

	return localFS{folderPath: folderPath}
}

func (localFS) Create(filename string) error {
	panic("unimplemented")
}

func (localFS) Read(filename string) ([]byte, error) {
	panic("unimplemented")
}

func (localFS) Delete(filename string) error {
	panic("unimplemented")
}
