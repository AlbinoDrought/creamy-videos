package files

import (
	"io"
	"os"
)

// ReadableFile is a file that we can read data from
type ReadableFile interface {
	io.ReadCloser
	io.Seeker
	Readdir(count int) ([]os.FileInfo, error)
	Stat() (os.FileInfo, error)
}

// WriteableFile is a file that we can write data to
type WriteableFile interface {
	io.WriteCloser
}

// FileSystem is an abstract file read/write interface
type FileSystem interface {
	MkdirAll(dirPath string, perm os.FileMode) error
	Create(name string) (WriteableFile, error)
	PipeTo(filePath string, reader io.Reader) error
	Remove(name string) error
	Stat(name string) (os.FileInfo, error)
	IsNotExist(err error) bool
	Open(path string) (ReadableFile, error)
}
