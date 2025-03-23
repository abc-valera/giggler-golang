package data

type iFS interface {
	Create(filename string) error
	Read(filename string) ([]byte, error)
	Delete(filename string) error
}
