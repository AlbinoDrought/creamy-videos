package files

import (
	"io"
	"net/http"
	"os"
	"path"
	"strings"
)

type ByteTransformer = func(original []byte)

// TransformedFile implements http.File but supports
// transforming the read stream
// If there's a better way to do this, please let me know
type TransformedFile struct {
	http.File
	innerFile   http.File
	transformer ByteTransformer
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
	n, err = tf.innerFile.Read(p)
	tf.transformer(p)
	return n, err
}

type WritingTransformedFile struct {
	file        *os.File
	transformer ByteTransformer
}

func (wtf WritingTransformedFile) Write(b []byte) (n int, err error) {
	wtf.transformer(b)
	return wtf.file.Write(b)
}

func (wtf WritingTransformedFile) Close() error {
	return wtf.file.Close()
}

// TransformedFileSystem custom file system handler
type TransformedFileSystem struct {
	fs          http.FileSystem
	dir         string
	transformer ByteTransformer
}

// TransformFileSystem using given FileTransformer
func TransformFileSystem(dir string, transformer ByteTransformer) TransformedFileSystem {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		os.MkdirAll(dir, os.ModePerm)
	}

	return TransformedFileSystem{
		fs:          http.Dir(dir),
		dir:         dir,
		transformer: transformer,
	}
}

func (fs TransformedFileSystem) MkdirAll(dirPath string, perm os.FileMode) error {
	return os.MkdirAll(path.Join(fs.dir, dirPath), perm)
}

func (fs TransformedFileSystem) Create(name string) (WritingTransformedFile, error) {
	file, err := os.Create(path.Join(fs.dir, name))
	if err != nil {
		return WritingTransformedFile{}, err
	}

	return WritingTransformedFile{
		file:        file,
		transformer: fs.transformer,
	}, nil
}

func (fs TransformedFileSystem) PipeTo(filePath string, reader io.Reader) error {
	file, err := fs.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, reader)

	return err
}

func (fs TransformedFileSystem) Stat(name string) (os.FileInfo, error) {
	return os.Stat(path.Join(fs.dir, name))
}

func (fs TransformedFileSystem) IsNotExist(err error) bool {
	return os.IsNotExist(err)
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
		if _, err := fs.fs.Open(index); err != nil { // possible mem leak here?
			return nil, err
		}
		// above fs.fs.Open file is never closed
	}

	return TransformedFile{
		innerFile:   f,
		transformer: fs.transformer,
	}, nil
}
