package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/AlbinoDrought/creamy-videos/files"
	"github.com/AlbinoDrought/creamy-videos/streamers"
)
import _ "net/http/pprof"

type UploadTransformer = func(io.Reader) io.Reader

const maxMultipartFormSize = 1024 * 1024 // 1KB

var videoRepo VideoRepo = NewDummyVideoRepo()

func uploadFileHandler(uploadTransformer UploadTransformer) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		if err := r.ParseMultipartForm(maxMultipartFormSize); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		defer r.MultipartForm.RemoveAll()

		// convert "foo, bar" and "foo,bar" into
		// ["foo", "bar"]
		tags := strings.Split(r.FormValue("tags"), ",")
		for i, tag := range tags {
			tags[i] = strings.Trim(tag, " ")
		}

		file, header, err := r.FormFile("file")
		defer file.Close()

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		stream := uploadTransformer(file)

		video := Video{
			Title:            r.FormValue("title"),
			Description:      r.FormValue("description"),
			OriginalFileName: header.Filename,
			Tags:             tags,
		}

		video, err = videoRepo.Upload(video, stream)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("error creating video: %+v", err)))
			return
		}

		json.NewEncoder(w).Encode(video)
	})
}

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

	http.Handle("/static/", http.StripPrefix(strings.TrimRight("/static/", "/"), fileServer))
	http.HandleFunc(
		"/api/upload",
		uploadFileHandler(func(reader io.Reader) io.Reader {
			return streamers.XorifyReader(reader, 0x69)
		}),
	)

	log.Printf("Serving %s on HTTP port: %s\n", *directory, *port)
	log.Fatal(http.ListenAndServe(":"+*port, nil))
}
