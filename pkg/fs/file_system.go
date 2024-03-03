package fs

import (
	"os"
)

// FileSystem is an interface for file system operations
type FileSystem interface {
	Create(name string) (*os.File, error)
	ReadFile(filename string) ([]byte, error)
	WriteFile(filename string, data []byte, perm os.FileMode) error
	Stat(name string) (os.FileInfo, error)
	Mkdir(name string, perm os.FileMode) error
	IsNotExist(err error) bool
}

var _ FileSystem = OSFileSystem{}

// OSFileSystem is an implementation of FileSystem
type OSFileSystem struct{}

// ReadFile reads a file
func (OSFileSystem) ReadFile(filename string) ([]byte, error) {
	return os.ReadFile(filename)
}

// WriteFile writes a file
func (OSFileSystem) WriteFile(filename string, data []byte, perm os.FileMode) error {
	return os.WriteFile(filename, data, perm)
}

// Stat stats a file
func (OSFileSystem) Stat(name string) (os.FileInfo, error) {
	return os.Stat(name)
}

// Mkdir makes a directory
func (OSFileSystem) Mkdir(name string, perm os.FileMode) error {
	return os.Mkdir(name, perm)
}

// IsNotExist checks if an error is a not exist error
func (OSFileSystem) IsNotExist(err error) bool {
	return os.IsNotExist(err)
}

// Create creates a file
func (OSFileSystem) Create(name string) (*os.File, error) {
	return os.Create(name)
}
