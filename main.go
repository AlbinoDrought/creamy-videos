package main

import (
	"flag"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/AlbinoDrought/creamy-videos/files"
	"github.com/AlbinoDrought/creamy-videos/streamers"
)

func main() {
	port := flag.String("p", "3000", "port to serve on")
	directory := flag.String("d", ".", "the directory of static file to host")
	flag.Parse()

	fileServer := http.FileServer(files.TransformFileSystem(
		http.Dir(*directory),
		func(file http.File) io.Reader {
			return streamers.XorifyReader(file, 0x69)
		},
	))
	http.Handle("/statics/", http.StripPrefix(strings.TrimRight("/statics/", "/"), fileServer))

	log.Printf("Serving %s on HTTP port: %s\n", *directory, *port)
	log.Fatal(http.ListenAndServe(":"+*port, nil))
}
