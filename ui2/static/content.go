package static

import "embed"

//go:embed css/* icons/* img/* js/* source.tar.gz favicon.ico
var FS embed.FS
