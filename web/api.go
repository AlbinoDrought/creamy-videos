package web

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path"
	"runtime/debug"
	"strconv"
	"strings"

	"github.com/AlbinoDrought/creamy-videos/files"
	"github.com/AlbinoDrought/creamy-videos/ui2/tmpl"
	"github.com/AlbinoDrought/creamy-videos/videostore"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
)

const (
	videosPerPage        = 32
	maxMultipartFormSize = 1024 * 1024 // 1MB
)

type CreamyVideosAPI interface {
	ListVideos(w http.ResponseWriter, r *http.Request)

	UploadVideo(w http.ResponseWriter, r *http.Request)

	ShowVideo(w http.ResponseWriter, r *http.Request)
	EditVideo(w http.ResponseWriter, r *http.Request)
	DeleteVideo(w http.ResponseWriter, r *http.Request)
}

func writeJSON(w http.ResponseWriter, thing any) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(thing)
}

type api struct {
	PublicURL tmpl.PublicURLGenerator
	FS        files.FileSystem
	Repo      videostore.VideoRepo
}

func (a *api) transformVideo(video videostore.Video) videostore.Video {
	video.Source = a.PublicURL(video.Source)
	if len(video.Thumbnail) > 0 {
		video.Thumbnail = a.PublicURL(video.Thumbnail)
	}
	return video
}

func (a *api) ListVideos(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	pageInt, err := page(r)
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

	videos, err := a.Repo.All(filter, uint(limit), uint(offset))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("error listing videos: %+v", err)
		return
	}

	transformedVideos := make([]videostore.Video, len(videos))
	for i, video := range videos {
		transformedVideos[i] = a.transformVideo(video)
	}

	writeJSON(w, transformedVideos)
}

func (a *api) UploadVideo(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	if err := r.ParseMultipartForm(maxMultipartFormSize); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Bad multipart/form-data request"))
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
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Bad file"))
		return
	}
	defer file.Close()

	video := videostore.Video{
		Title:            r.FormValue("title"),
		Description:      r.FormValue("description"),
		OriginalFileName: header.Filename,
		Tags:             tags,
	}

	video, err = a.Repo.Save(video)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("error creating video: %+v", err)))
		return
	}

	rootDir := strconv.Itoa(int(video.ID))
	if _, err := a.FS.Stat(rootDir); a.FS.IsNotExist(err) {
		a.FS.MkdirAll(rootDir, os.ModePerm)
	}

	videoPath := path.Join(rootDir, "video"+path.Ext(video.OriginalFileName))

	err = files.PipeTo(a.FS, videoPath, file)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("error saving video stream: %+v", err)))
		return
	}

	video.Source = videoPath
	video, err = a.Repo.Save(video)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("error setting video source: %+v", err)))
		return
	}

	go func() {
		_, err := videostore.GenerateThumbnail(video, a.Repo, a.FS)
		if err != nil {
			log.Printf("failed to make thumbnail: %+v", err)
		}
	}()

	go debug.FreeOSMemory() // hack to request our memory back :'(

	w.WriteHeader(http.StatusCreated)
	writeJSON(w, a.transformVideo(video))
}

func (a *api) ShowVideo(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	vars := mux.Vars(r)
	rawID := vars["id"]
	id, err := strconv.Atoi(rawID)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	video, err := a.Repo.FindById(uint(id))
	if err == videostore.ErrorVideoNotFound {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("error while retrieving video: %+v", err)
		return
	}

	writeJSON(w, a.transformVideo(video))
}

func (a *api) EditVideo(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	vars := mux.Vars(r)
	rawID := vars["id"]
	id, err := strconv.Atoi(rawID)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	video, err := a.Repo.FindById(uint(id))
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

	video, err = a.Repo.Save(video)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("error updating video: %+v", err)
		return
	}

	writeJSON(w, a.transformVideo(video))
}

func (a *api) DeleteVideo(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	vars := mux.Vars(r)
	rawID := vars["id"]
	id, err := strconv.Atoi(rawID)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	video, err := a.Repo.FindById(uint(id))
	if err == videostore.ErrorVideoNotFound {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("error while retrieving video: %+v", err)
		return
	}

	err = a.Repo.Delete(video)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("error deleting video: %+v", err)
		return
	}

	_, err = a.FS.Stat(video.Source)
	if !a.FS.IsNotExist(err) {
		// video exists, attempt to delete
		err = a.FS.Remove(video.Source)
		if err != nil {
			log.Print(errors.Wrap(err, "failed to remove video from disk"))
		}
	}

	_, err = a.FS.Stat(video.Thumbnail)
	if !a.FS.IsNotExist(err) {
		// thumbnail exists, attempt to delete
		err = a.FS.Remove(video.Thumbnail)
		if err != nil {
			log.Print(errors.Wrap(err, "failed to remove thumbnail from disk"))
		}
	}

	writeJSON(w, a.transformVideo(video))
}

func newAPI(PublicURL tmpl.PublicURLGenerator, FS files.FileSystem, Repo videostore.VideoRepo) CreamyVideosAPI {
	return &api{PublicURL, FS, Repo}
}

func NewWriteableAPI(PublicURL tmpl.PublicURLGenerator, FS files.FileSystem, Repo videostore.VideoRepo) http.Handler {
	api := newAPI(PublicURL, FS, Repo)
	r := mux.NewRouter()

	r.HandleFunc(
		"/api/video",
		api.ListVideos,
	)

	r.HandleFunc(
		"/api/video/{id:[0-9]+}",
		api.ShowVideo,
	).Methods("GET")

	r.HandleFunc(
		"/api/video/{id:[0-9]+}",
		api.EditVideo,
	).Methods("POST")

	r.HandleFunc(
		"/api/video/{id:[0-9]+}",
		api.DeleteVideo,
	).Methods("DELETE")

	r.HandleFunc(
		"/api/upload",
		api.UploadVideo,
	)

	return r
}

func NewReadOnlyAPI(PublicURL tmpl.PublicURLGenerator, FS files.FileSystem, Repo videostore.VideoRepo) http.Handler {
	api := newAPI(PublicURL, FS, Repo)
	r := mux.NewRouter()

	r.HandleFunc(
		"/api/video",
		api.ListVideos,
	)

	r.HandleFunc(
		"/api/video/{id:[0-9]+}",
		api.ShowVideo,
	).Methods("GET")

	return r
}
