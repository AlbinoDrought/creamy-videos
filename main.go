package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/AlbinoDrought/creamy-videos/files"
	"github.com/AlbinoDrought/creamy-videos/streamers"

	_ "net/http/pprof"
)

const maxMultipartFormSize = 1024 * 1024 // 1MB

var videoRepo VideoRepo
var transformedFileSystem files.TransformedFileSystem

func uploadFileHandler() http.HandlerFunc {
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

		video := Video{
			Title:            r.FormValue("title"),
			Description:      r.FormValue("description"),
			OriginalFileName: header.Filename,
			Tags:             tags,
		}

		video, err = videoRepo.Upload(video, file)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("error creating video: %+v", err)))
			return
		}

		json.NewEncoder(w).Encode(video)
	})
}

const videosPerPage = 30

func listVideosHandler() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		page := r.URL.Query().Get("page")
		if len(page) <= 0 {
			page = "1"
		}
		pageInt, err := strconv.Atoi(page)

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		limit := videosPerPage
		offset := videosPerPage * (pageInt - 1)
		if offset < 0 {
			offset = 0
		}

		videos, err := videoRepo.All(uint(limit), uint(offset))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(videos)
	})
}

func main() {
	port := flag.String("p", "3000", "port to serve on")
	directory := flag.String("d", ".", "the directory of static file to host")
	flag.Parse()

	transformedFileSystem = files.TransformFileSystem(
		http.Dir(*directory),
		func(reader io.Reader) io.Reader {
			return streamers.XorifyReader(reader, 0x69)
		},
	)

	videoRepo = NewDummyVideoRepo()

	fileServer := http.FileServer(transformedFileSystem)

	http.Handle("/static/", http.StripPrefix(strings.TrimRight("/static/", "/"), fileServer))
	http.HandleFunc(
		"/api/video",
		listVideosHandler(),
	)
	http.HandleFunc(
		"/api/upload",
		uploadFileHandler(),
	)

	log.Printf("Serving %s on HTTP port: %s\n", *directory, *port)
	log.Fatal(http.ListenAndServe(":"+*port, nil))
}
