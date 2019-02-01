package cmd

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"runtime/debug"
	"strconv"
	"strings"

	"github.com/AlbinoDrought/creamy-videos/files"
	"github.com/AlbinoDrought/creamy-videos/videostore"
	packr "github.com/gobuffalo/packr/v2"
	"github.com/spf13/cobra"
)

const maxMultipartFormSize = 1024 * 1024 // 1MB

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
}

func uploadFileHandler(instance application) http.HandlerFunc {
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

const videosPerPage = 30

func listVideosHandler(instance application) http.HandlerFunc {
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

		videos, err := instance.repo.All(uint(limit), uint(offset))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		transformedVideos := make([]videostore.Video, len(videos))
		for i, video := range videos {
			video.Source = instance.config.AppURL + instance.config.HTTPVideoDirectory + video.Source
			if len(video.Thumbnail) > 0 {
				video.Thumbnail = instance.config.AppURL + instance.config.HTTPVideoDirectory + video.Thumbnail
			}
			transformedVideos[i] = video
		}

		json.NewEncoder(w).Encode(transformedVideos)
	})
}

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Provide videos, UI, and API over HTTP",
	Run: func(cmd *cobra.Command, args []string) {
		fileServer := http.FileServer(app.fs)

		box := packr.New("spa", "./../ui/dist")

		http.Handle("/", http.FileServer(files.CreateSPAFileSystem(box, "/index.html")))

		http.Handle(
			app.config.HTTPVideoDirectory,
			http.StripPrefix(
				strings.TrimRight(app.config.HTTPVideoDirectory, "/"),
				fileServer,
			),
		)
		http.HandleFunc(
			"/api/video",
			listVideosHandler(app),
		)
		http.HandleFunc(
			"/api/upload",
			uploadFileHandler(app),
		)

		log.Printf("Remote URL: %s\n", app.config.AppURL)
		log.Printf("Serving videos from %s on %s\n", app.config.LocalVideoDirectory, app.config.HTTPVideoDirectory)
		log.Printf("Listening on %s\n", app.config.Port)
		log.Fatal(http.ListenAndServe(":"+app.config.Port, nil))
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
}
