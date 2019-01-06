package files

import (
	"io"
	"net/http"
	"os"
	"strings"
)

type FileTransformer = func(http.File) io.Reader

// TransformedFile implements http.File but supports
// transforming the read stream
// If there's a better way to do this, please let me know
type TransformedFile struct {
	http.File
	innerFile   http.File
	transformer FileTransformer
}

// Close calls the parent file Close()
func (tf TransformedFile) Close() error {
	return tf.innerFile.Close()
}

// Seek calls the parent file Seek
func (tf TransformedFile) Seek(offset int64, whence int) (int64, error) {
	return tf.innerFile.Seek(offset, whence)
}

// Readdir calls the parent file Readdir
func (tf TransformedFile) Readdir(count int) ([]os.FileInfo, error) {
	return tf.innerFile.Readdir(count)
}

// Stat calls the parent file Stat
func (tf TransformedFile) Stat() (os.FileInfo, error) {
	return tf.innerFile.Stat()
}

// Read transforms the parent implementation of Read
func (tf TransformedFile) Read(p []byte) (n int, err error) {
	return tf.transformer(tf.innerFile).Read(p)
}

// TransformedFileSystem custom file system handler
type TransformedFileSystem struct {
	fs          http.FileSystem
	transformer FileTransformer
}

// TransformFileSystem using given FileTransformer
func TransformFileSystem(fs http.FileSystem, transformer FileTransformer) TransformedFileSystem {
	return TransformedFileSystem{
		fs:          fs,
		transformer: transformer,
	}
}

// Open opens file, prevents directory listing
// https://gist.github.com/hauxe/f2ea1901216177ccf9550a1b8bd59178#file-http_static_correct-go
func (fs TransformedFileSystem) Open(path string) (http.File, error) {
	f, err := fs.fs.Open(path)
	if err != nil {
		return nil, err
	}

	s, err := f.Stat()
	if s.IsDir() {
		index := strings.TrimSuffix(path, "/") + "/index.html"
		if _, err := fs.fs.Open(index); err != nil {
			return nil, err
		}
	}

	return TransformedFile{
		innerFile:   f,
		transformer: fs.transformer,
	}, nil
}
