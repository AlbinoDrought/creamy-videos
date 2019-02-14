package files

import (
	"net/http"
	"path"
)

type httpCompatibleFileSystem struct {
	fs                FileSystem
	directoryListings bool
}

func (fs httpCompatibleFileSystem) Open(name string) (http.File, error) {
	file, err := fs.fs.Open(name)

	if err == nil && fs.directoryListings {
		stat, err := file.Stat()
		if err != nil {
			return nil, err
		}

		if stat.IsDir() {
			// attempt to load index.html instead
			index := path.Join(name, "index.html")
			indexFile, err := fs.fs.Open(index)
			if err != nil {
				// index.html doesn't exist, do not show directory listing
				return nil, err
			}

			// index.html exists, show it instead
			file.Close()
			file = indexFile
		}
	}

	return file, err
}

// AdaptToHTTPFileSystem converts any FileSystem implementation
// into one that can be used with http.FileServer
func AdaptToHTTPFileSystem(fs FileSystem, directoryListings bool) http.FileSystem {
	return httpCompatibleFileSystem{
		fs,
		directoryListings,
	}
}
