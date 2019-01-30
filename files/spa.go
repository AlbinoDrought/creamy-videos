package files

import (
	"net/http"
	"os"
	"strings"
)

// SPAFileSystem responds with the given
// fallback file if the requested path
// does not exist.
type SPAFileSystem struct {
	fs       http.FileSystem
	fallback string
}

// CreateSPAFileSystem that falls back to `fallback`
func CreateSPAFileSystem(fs http.FileSystem, fallback string) SPAFileSystem {
	return SPAFileSystem{
		fs,
		fallback,
	}
}

func (fs SPAFileSystem) forceFallback() (http.File, error) {
	return fs.fs.Open(fs.fallback)
}

// Open a file, or respond with the fallback contents
func (fs SPAFileSystem) Open(name string) (http.File, error) {
	file, err := fs.fs.Open(name)

	// prevent redirect loops, ignore root /
	if strings.TrimLeft(name, "/") == "" {
		return file, err
	}

	// load SPA if file doesn't exist
	if err != nil && os.IsNotExist(err) {
		return fs.forceFallback()
	}

	// load SPA if this is a directory listing
	// (this always happens with packr2 for some reason)
	s, err := file.Stat()
	if s.IsDir() {
		_ = file.Close()
		return fs.forceFallback()
	}

	// file exists and is not a directory,
	// actually return it
	return file, err
}
