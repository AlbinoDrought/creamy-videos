package files

import (
	"bytes"
	"io"
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
)

// dummy func to fulfill FileTransformer
var passthrough = func(b []byte) {}

func TestMkdirAll(t *testing.T) {
	root := "test-mkdirall"
	directories := []string{
		"some-dir/foo/bar/foobar/baz",
		"other/hello",
		"single",
		"trailing/slash/",
	}

	tfs := TransformFileSystem(
		LocalFileSystem(root),
		passthrough,
	)
	defer os.RemoveAll(root)

	for _, directory := range directories {
		err := tfs.MkdirAll(directory, os.ModePerm)
		assert.Nil(t, err)

		assert.DirExists(t, path.Join(root, directory))
	}
}

func TestStat(t *testing.T) {
	root := "test-stat"

	tfs := TransformFileSystem(
		LocalFileSystem("."),
		passthrough,
	)
	defer os.RemoveAll(root)

	// assert dir doesn't exist
	fileInfo, err := tfs.Stat(root)
	assert.Nil(t, fileInfo)
	assert.True(t, tfs.IsNotExist(err))

	// create dir
	err = os.Mkdir(root, os.ModePerm)
	assert.Nil(t, err)

	// assert exists
	fileInfo, err = tfs.Stat(root)
	assert.NotNil(t, fileInfo)
	assert.False(t, tfs.IsNotExist(err))
}

func TestOpen(t *testing.T) {
	root := "test-open"

	tfs := TransformFileSystem(
		LocalFileSystem(root),
		passthrough,
	)
	defer os.RemoveAll(root)

	filepath := "foo.bin"
	b := []byte{
		0xDE,
		0xAD,
		0xBE,
		0xEF,
	}
	reader := bytes.NewReader(b)

	// create dummy file
	file, err := os.Create(path.Join(root, filepath))
	assert.Nil(t, err)
	io.Copy(file, reader)
	file.Close()

	// read dummy file
	transformedFile, err := tfs.Open(filepath)
	assert.Nil(t, err)
	readBytes := make([]byte, len(b))
	bytesRead, err := transformedFile.Read(readBytes)
	transformedFile.Close()

	assert.Equal(t, len(b), bytesRead)
	assert.Nil(t, err)

	assert.Equal(t, b, readBytes)
}
