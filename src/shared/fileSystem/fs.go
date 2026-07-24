package fileSystem

type Interface interface {
	Create(filename string) error
	Read(filename string) ([]byte, error)
	Delete(filename string) error
}
