package cmd

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"runtime/debug"
	"strconv"
	"strings"

	"github.com/gorilla/mux"

	"github.com/AlbinoDrought/creamy-videos/files"
	"github.com/AlbinoDrought/creamy-videos/videostore"
	packr "github.com/gobuffalo/packr/v2"
	"github.com/spf13/cobra"
)

const maxMultipartFormSize = 1024 * 1024 // 1MB

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		next.ServeHTTP(w, r)
	})
}

func uploadFileHandler(instance application) http.HandlerFunc {
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

		video := videostore.Video{
			Title:            r.FormValue("title"),
			Description:      r.FormValue("description"),
			OriginalFileName: header.Filename,
			Tags:             tags,
		}

		video, err = instance.repo.Upload(video, file)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("error creating video: %+v", err)))
			return
		}

		go debug.FreeOSMemory() // hack to request our memory back :'(

		json.NewEncoder(w).Encode(video)
	})
}

const videosPerPage = 32

type stringDict interface {
	Get(key string) string
}

func videoFilterFromDict(dict stringDict) videostore.VideoFilter {
	// don't do strings.Split("", ",")
	// that would give us a slice with length=1,
	// containing an empty string
	rawTags := dict.Get("tags")
	var tags []string
	if len(rawTags) > 0 {
		tags = strings.Split(rawTags, ",")
	} else {
		tags = make([]string, 0)
	}

	return videostore.VideoFilter{
		Title: dict.Get("title"),
		Tags:  tags,
		Any:   dict.Get("filter"),
	}
}

func transformVideo(instance application, video videostore.Video) videostore.Video {
	video.Source = instance.config.AppURL + instance.config.HTTPVideoDirectory + video.Source
	if len(video.Thumbnail) > 0 {
		video.Thumbnail = instance.config.AppURL + instance.config.HTTPVideoDirectory + video.Thumbnail
	}

	return video
}

func deleteVideoHandler(instance application) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		vars := mux.Vars(r)
		rawID := vars["id"]
		id, err := strconv.Atoi(rawID)

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		video, err := instance.repo.FindById(uint(id))
		if err == videostore.ErrorVideoNotFound {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Printf("error while retrieving video: %+v", err)
			return
		}

		err = instance.repo.Delete(video)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Printf("error deleting video: %+v", err)
			return
		}

		json.NewEncoder(w).Encode(transformVideo(instance, video))
	})
}

func editVideoHandler(instance application) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		vars := mux.Vars(r)
		rawID := vars["id"]
		id, err := strconv.Atoi(rawID)

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		video, err := instance.repo.FindById(uint(id))
		if err == videostore.ErrorVideoNotFound {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		postedVideo := videostore.Video{}
		err = json.NewDecoder(r.Body).Decode(&postedVideo)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		video.Title = postedVideo.Title
		video.Description = postedVideo.Description
		video.Tags = postedVideo.Tags

		video, err = instance.repo.Save(video)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Printf("error updating video: %+v", err)
			return
		}

		json.NewEncoder(w).Encode(transformVideo(instance, video))
	})
}

func showVideoHandler(instance application) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		vars := mux.Vars(r)
		rawID := vars["id"]
		id, err := strconv.Atoi(rawID)

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		video, err := instance.repo.FindById(uint(id))
		if err == videostore.ErrorVideoNotFound {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Printf("error while retrieving video: %+v", err)
			return
		}

		json.NewEncoder(w).Encode(transformVideo(instance, video))
	})
}

func listVideosHandler(instance application) http.HandlerFunc {
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

		filter := videoFilterFromDict(r.URL.Query())

		videos, err := instance.repo.All(filter, uint(limit), uint(offset))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Printf("error listing videos: %+v", err)
			return
		}

		transformedVideos := make([]videostore.Video, len(videos))
		for i, video := range videos {
			transformedVideos[i] = transformVideo(instance, video)
		}

		json.NewEncoder(w).Encode(transformedVideos)
	})
}

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Provide videos, UI, and API over HTTP",
	Run: func(cmd *cobra.Command, args []string) {
		fileServer := http.FileServer(files.AdaptToHTTPFileSystem(app.fs, false))

		box := packr.New("spa", "./../ui/dist")

		r := mux.NewRouter()

		r.HandleFunc(
			"/api/video",
			listVideosHandler(app),
		)

		r.HandleFunc(
			"/api/video/{id:[0-9]+}",
			showVideoHandler(app),
		).Methods("GET")
		r.HandleFunc(
			"/api/video/{id:[0-9]+}",
			editVideoHandler(app),
		).Methods("POST")
		r.HandleFunc(
			"/api/video/{id:[0-9]+}",
			deleteVideoHandler(app),
		).Methods("DELETE")

		r.HandleFunc(
			"/api/upload",
			uploadFileHandler(app),
		)

		r.PathPrefix(app.config.HTTPVideoDirectory).Handler(
			http.StripPrefix(
				strings.TrimRight(app.config.HTTPVideoDirectory, "/"),
				fileServer,
			),
		)
		r.PathPrefix("/").Handler(http.FileServer(files.CreateSPAFileSystem(box, "/index.html")))

		r.Use(corsMiddleware)

		http.Handle("/", r)

		log.Printf("Remote URL: %s\n", app.config.AppURL)
		log.Printf("Serving videos from %s on %s\n", app.config.LocalVideoDirectory, app.config.HTTPVideoDirectory)
		log.Printf("Listening on %s\n", app.config.Port)
		log.Fatal(http.ListenAndServe(":"+app.config.Port, nil))
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
}
