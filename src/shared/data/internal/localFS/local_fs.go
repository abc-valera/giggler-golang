package localFS

type localFS struct {
	folderPath string
}

func New(folderPath string) localFS {
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
