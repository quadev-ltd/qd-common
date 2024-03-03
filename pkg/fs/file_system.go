package fs

import (
	"os"
)

type FileSystem interface {
	Create(name string) (*os.File, error)
	ReadFile(filename string) ([]byte, error)
	WriteFile(filename string, data []byte, perm os.FileMode) error
	Stat(name string) (os.FileInfo, error)
	Mkdir(name string, perm os.FileMode) error
	IsNotExist(err error) bool
}

var _ FileSystem = OSFileSystem{}

type OSFileSystem struct{}

func (OSFileSystem) ReadFile(filename string) ([]byte, error) {
	return os.ReadFile(filename)
}

func (OSFileSystem) WriteFile(filename string, data []byte, perm os.FileMode) error {
	return os.WriteFile(filename, data, perm)
}

func (OSFileSystem) Stat(name string) (os.FileInfo, error) {
	return os.Stat(name)
}

func (OSFileSystem) Mkdir(name string, perm os.FileMode) error {
	return os.Mkdir(name, perm)
}

func (OSFileSystem) IsNotExist(err error) bool {
	return os.IsNotExist(err)
}

func (OSFileSystem) Create(name string) (*os.File, error) {
	return os.Create(name)
}
