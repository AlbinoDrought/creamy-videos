package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/AlbinoDrought/creamy-videos/files"
	"github.com/AlbinoDrought/creamy-videos/streamers"
	packr "github.com/gobuffalo/packr/v2"

	_ "net/http/pprof"
)

const maxMultipartFormSize = 1024 * 1024 // 1MB
const appUrl = "http://localhost:3000"

var config Config
var videoRepo VideoRepo
var transformedFileSystem files.TransformedFileSystem

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
}

func uploadFileHandler() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		enableCors(&w)
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
		enableCors(&w)
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

		transformedVideos := make([]Video, len(videos))
		for i, video := range videos {
			video.Source = config.AppUrl + config.HttpVideoDirectory + video.Source
			if len(video.Thumbnail) > 0 {
				video.Thumbnail = config.AppUrl + config.HttpVideoDirectory + video.Thumbnail
			}
			transformedVideos[i] = video
		}

		json.NewEncoder(w).Encode(transformedVideos)
	})
}

func main() {
	config = FillFromEnv()

	transformedFileSystem = files.TransformFileSystem(
		config.LocalVideoDirectory,
		func(reader io.Reader) io.Reader {
			return streamers.XorifyReader(reader, 0x69)
		},
	)

	videoRepo = NewDummyVideoRepo()

	// ghetto thumbnail regen
	/*
		videos, _ := videoRepo.All(1000, 0)
		for _, video := range videos {
			go eventuallyMakeThumbnail(video)
		}
	*/

	fileServer := http.FileServer(transformedFileSystem)

	box := packr.New("spa", "./ui/dist")

	http.Handle("/", http.FileServer(box))

	http.Handle(
		config.HttpVideoDirectory,
		http.StripPrefix(
			strings.TrimRight(config.HttpVideoDirectory, "/"),
			fileServer,
		),
	)
	http.HandleFunc(
		"/api/video",
		listVideosHandler(),
	)
	http.HandleFunc(
		"/api/upload",
		uploadFileHandler(),
	)

	log.Printf("Remote URL: %s\n", config.AppUrl)
	log.Printf("Serving videos from %s on %s\n", config.LocalVideoDirectory, config.HttpVideoDirectory)
	log.Printf("Listening on %s\n", config.Port)
	log.Fatal(http.ListenAndServe(":"+config.Port, nil))
}
