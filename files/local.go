package files

import (
	"os"
	"path"
)

type localFileSystem struct {
	dir string
}

// LocalFileSystem provides access to things in a
// local directory
func LocalFileSystem(dir string) FileSystem {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		os.MkdirAll(dir, os.ModePerm)
	}

	return localFileSystem{
		dir,
	}
}

func (fs localFileSystem) localize(part string) string {
	return path.Join(fs.dir, path.Clean(part))
}

func (fs localFileSystem) MkdirAll(dirPath string, perm os.FileMode) error {
	return os.MkdirAll(fs.localize(dirPath), perm)
}

func (fs localFileSystem) Create(name string) (WriteableFile, error) {
	return os.Create(fs.localize(name))
}

func (fs localFileSystem) Remove(name string) error {
	return os.Remove(fs.localize(name))
}

func (fs localFileSystem) Stat(name string) (os.FileInfo, error) {
	return os.Stat(fs.localize(name))
}

func (fs localFileSystem) IsNotExist(err error) bool {
	return os.IsNotExist(err)
}

func (fs localFileSystem) Open(path string) (ReadableFile, error) {
	return os.Open(fs.localize(path))
}
