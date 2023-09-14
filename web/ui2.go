package web

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
	"runtime/debug"
	"strconv"
	"strings"

	"github.com/AlbinoDrought/creamy-videos/files"
	"github.com/AlbinoDrought/creamy-videos/ui2/static"
	"github.com/AlbinoDrought/creamy-videos/ui2/tmpl"
	"github.com/AlbinoDrought/creamy-videos/videostore"
	"github.com/gorilla/mux"
)

type CreamyVideosUI2 interface {
	Home(w http.ResponseWriter, r *http.Request)
	Search(w http.ResponseWriter, r *http.Request)
	Watch(w http.ResponseWriter, r *http.Request)

	UploadForm(w http.ResponseWriter, r *http.Request)
	Upload(w http.ResponseWriter, r *http.Request)

	// todo: Upload, Show, Edit, Delete UI & Handler routes
}

type sortDir map[string]string

func (s sortDir) Get(key string) string {
	v, ok := s[key]
	if ok {
		return v
	}
	return ""
}

var sortDirs = map[string]sortDir{
	"newest": sortDir(map[string]string{
		"sort_field":     "time_created",
		"sort_direction": videostore.SortDirectionDescending,
	}),
	"oldest": sortDir(map[string]string{
		"sort_field":     "time_created",
		"sort_direction": videostore.SortDirectionAscending,
	}),
	"az": sortDir(map[string]string{
		"sort_field":     "title",
		"sort_direction": videostore.SortDirectionAscending,
	}),
	"za": sortDir(map[string]string{
		"sort_field":     "title",
		"sort_direction": videostore.SortDirectionDescending,
	}),
}

var defaultSortDir = "newest"

type cUI2 struct {
	ReadOnly  bool
	PublicURL tmpl.PublicURLGenerator
	FS        files.FileSystem
	Repo      videostore.VideoRepo
}

func (u *cUI2) WriteErrorPage(w http.ResponseWriter, r *http.Request, statusCode int, err error, msg string) {
	w.WriteHeader(statusCode)
	w.Write([]byte("todo"))
	log.Printf("%v error: %v", msg, err)
}

func (u *cUI2) Home(w http.ResponseWriter, r *http.Request) {
	pageInt, err := page(r)
	if err != nil {
		u.WriteErrorPage(w, r, http.StatusBadRequest, err, "bad page number")
		return
	}

	limit := videosPerPage
	offset := videosPerPage * (pageInt - 1)
	if offset < 0 {
		offset = 0
	}

	filter := videoFilterFromDict(sortDir(map[string]string{
		"tags":           "home",
		"sort_field":     "time_created",
		"sort_direction": videostore.SortDirectionDescending,
	}))

	videos, err := u.Repo.All(filter, uint(limit), uint(offset))
	if err != nil {
		u.WriteErrorPage(w, r, http.StatusInternalServerError, err, "failed listing videos")
		return
	}

	count, err := u.Repo.Count(filter)
	if err != nil {
		u.WriteErrorPage(w, r, http.StatusInternalServerError, err, "failed counting videos")
		return
	}

	w.Header().Add("Content-Type", "text/html")
	tmpl.Home(tmpl.AppState{
		ReadOnly:      u.ReadOnly,
		SortDirection: "",
		Sortable:      false,
		SearchText:    "",
		PUG:           u.PublicURL,
	}, tmpl.Paging{
		URL: func(p int) string {
			return fmt.Sprintf("/?page=%v", p)
		},
		CurrentPage: pageInt,
		Pages:       int(pages(count, videosPerPage)),
	}, videos).Render(r.Context(), w)
}

func (u *cUI2) Search(w http.ResponseWriter, r *http.Request) {
	pageInt, err := page(r)
	if err != nil {
		u.WriteErrorPage(w, r, http.StatusBadRequest, err, "bad page number")
		return
	}

	limit := videosPerPage
	offset := videosPerPage * (pageInt - 1)
	if offset < 0 {
		offset = 0
	}

	sort := r.URL.Query().Get("sort")
	_, ok := sortDirs[sort]
	if !ok {
		sort = defaultSortDir
	}

	filterArgs := map[string]string{
		"tags":   r.URL.Query().Get("tags"),
		"title":  r.URL.Query().Get("title"),
		"filter": r.URL.Query().Get("text"),
	}
	for k, v := range sortDirs[sort] {
		filterArgs[k] = v
	}

	filter := videoFilterFromDict(sortDir(filterArgs))

	videos, err := u.Repo.All(filter, uint(limit), uint(offset))
	if err != nil {
		u.WriteErrorPage(w, r, http.StatusInternalServerError, err, "failed listing videos")
		return
	}

	count, err := u.Repo.Count(filter)
	if err != nil {
		u.WriteErrorPage(w, r, http.StatusInternalServerError, err, "failed counting videos")
		return
	}

	baseURL := fmt.Sprintf("/search?sort=%v&", url.QueryEscape(sort))
	for _, key := range []string{"tags", "title", "filter"} {
		v := filterArgs[key]
		if v == "" {
			continue
		}
		baseURL = fmt.Sprintf("%v%v=%v&", baseURL, key, url.QueryEscape(v))
	}

	w.Header().Add("Content-Type", "text/html")
	tmpl.Search(tmpl.AppState{
		ReadOnly:      u.ReadOnly,
		SortDirection: sort,
		Sortable:      true,
		SearchText:    r.URL.Query().Get("text"),
		PUG:           u.PublicURL,
	}, tmpl.Paging{
		URL: func(p int) string {
			return fmt.Sprintf("%vpage=%v", baseURL, p)
		},
		CurrentPage: pageInt,
		Pages:       int(pages(count, videosPerPage)),
	}, videos).Render(r.Context(), w)
}

func (u *cUI2) Watch(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	rawID := vars["id"]
	id, err := strconv.Atoi(rawID)
	if err != nil {
		u.WriteErrorPage(w, r, http.StatusBadRequest, err, "bad ID")
		return
	}

	video, err := u.Repo.FindById(uint(id))
	if err == videostore.ErrorVideoNotFound {
		u.WriteErrorPage(w, r, http.StatusNotFound, err, "video not found")
		return
	}
	if err != nil {
		u.WriteErrorPage(w, r, http.StatusInternalServerError, err, "failed finding video")
		return
	}

	w.Header().Add("Content-Type", "text/html")
	tmpl.Watch(tmpl.AppState{
		ReadOnly:      u.ReadOnly,
		SortDirection: "",
		Sortable:      false,
		SearchText:    "",
		PUG:           u.PublicURL,
	}, video).Render(r.Context(), w)
}

func (u *cUI2) UploadForm(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/html")
	tmpl.UploadForm(tmpl.AppState{
		ReadOnly:      u.ReadOnly,
		SortDirection: "",
		Sortable:      false,
		SearchText:    "",
		PUG:           u.PublicURL,
	}, tmpl.UploadFormState{}).Render(r.Context(), w)
}

func (u *cUI2) Upload(w http.ResponseWriter, r *http.Request) {
	writeErrorPage := func(statusCode int, err error, msg string) {
		log.Printf("%v error: %v", msg, err)
		w.Header().Add("Content-Type", "text/html")
		w.WriteHeader(statusCode)
		tmpl.UploadForm(tmpl.AppState{
			ReadOnly:      u.ReadOnly,
			SortDirection: "",
			Sortable:      false,
			SearchText:    "",
			PUG:           u.PublicURL,
		}, tmpl.UploadFormState{
			Error: msg,

			Title:       r.FormValue("title"),
			Tags:        r.FormValue("tags"),
			Description: r.FormValue("description"),
		}).Render(r.Context(), w)
	}

	defer r.Body.Close()

	if err := r.ParseMultipartForm(maxMultipartFormSize); err != nil {
		writeErrorPage(http.StatusBadRequest, err, "Bad multipart/form-data request")
		return
	}
	defer r.MultipartForm.RemoveAll()

	// convert "foo, bar" and "foo,bar" into
	// ["foo", "bar"]
	tags := strings.Split(r.FormValue("tags"), ",")
	for i, tag := range tags {
		tags[i] = strings.Trim(tag, " ")
	}
	if len(tags) == 1 && tags[0] == "" {
		tags = []string{}
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		writeErrorPage(http.StatusBadRequest, err, "Bad file")
		return
	}
	defer file.Close()

	video := videostore.Video{
		Title:            r.FormValue("title"),
		Description:      r.FormValue("description"),
		OriginalFileName: header.Filename,
		Tags:             tags,
	}

	video, err = u.Repo.Save(video)
	if err != nil {
		writeErrorPage(http.StatusInternalServerError, err, "Internal error creating video resource")
		return
	}

	rootDir := strconv.Itoa(int(video.ID))
	if _, err := u.FS.Stat(rootDir); u.FS.IsNotExist(err) {
		u.FS.MkdirAll(rootDir, os.ModePerm)
	}

	videoPath := path.Join(rootDir, "video"+path.Ext(video.OriginalFileName))

	err = files.PipeTo(u.FS, videoPath, file)
	if err != nil {
		writeErrorPage(http.StatusInternalServerError, err, "Internal error saving video stream")
		return
	}

	video.Source = videoPath
	video, err = u.Repo.Save(video)

	if err != nil {
		writeErrorPage(http.StatusInternalServerError, err, "Internal error setting video source")
		return
	}

	go func() {
		_, err := videostore.GenerateThumbnail(video, u.Repo, u.FS)
		if err != nil {
			log.Printf("failed to make thumbnail: %+v", err)
		}
	}()

	go debug.FreeOSMemory() // hack to request our memory back :'(

	http.Redirect(w, r, fmt.Sprintf("/watch/%v", video.ID), http.StatusFound)
}

func NewWriteableCUI2(publicURL tmpl.PublicURLGenerator, fs files.FileSystem, repo videostore.VideoRepo) http.Handler {
	u := &cUI2{
		ReadOnly:  false,
		PublicURL: publicURL,
		FS:        fs,
		Repo:      repo,
	}

	r := mux.NewRouter()

	fileServer := http.FileServer(http.FS(static.FS))
	r.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Cache-Control", "public, max-age=86400, stale-while-revalidate")
		fileServer.ServeHTTP(w, r)
	})

	r.HandleFunc(
		"/",
		u.Home,
	).Methods("GET")

	r.HandleFunc(
		"/search",
		u.Search,
	).Methods("GET")

	r.HandleFunc(
		"/upload",
		u.UploadForm,
	).Methods("GET")
	r.HandleFunc(
		"/upload",
		u.Upload,
	).Methods("POST")

	r.HandleFunc(
		"/watch/{id:[0-9]+}",
		u.Watch,
	).Methods("GET")

	return r
}

func NewReadOnlyCUI2(publicURL tmpl.PublicURLGenerator, repo videostore.VideoRepo) http.Handler {
	u := &cUI2{
		ReadOnly:  true,
		PublicURL: publicURL,
		Repo:      repo,
	}

	r := mux.NewRouter()

	fileServer := http.FileServer(http.FS(static.FS))
	r.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Cache-Control", "public, max-age=86400, stale-while-revalidate")
		fileServer.ServeHTTP(w, r)
	})

	r.HandleFunc(
		"/",
		u.Home,
	).Methods("GET")

	return r
}
