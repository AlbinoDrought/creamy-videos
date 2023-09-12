package web

import (
	"log"
	"net/http"

	"github.com/AlbinoDrought/creamy-videos/ui2/static"
	"github.com/AlbinoDrought/creamy-videos/ui2/tmpl"
	"github.com/AlbinoDrought/creamy-videos/videostore"
	"github.com/gorilla/mux"
)

type CreamyVideosUI2 interface {
	Home(w http.ResponseWriter, r *http.Request)
	Search(w http.ResponseWriter, r *http.Request)

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

	w.Header().Add("Content-Type", "text/html")
	tmpl.Home(tmpl.AppState{
		ReadOnly:      u.ReadOnly,
		SortDirection: "",
		Sortable:      false,
		SearchText:    "",
		PUG:           u.PublicURL,
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

	w.Header().Add("Content-Type", "text/html")

	tmpl.Search(tmpl.AppState{
		ReadOnly:      u.ReadOnly,
		SortDirection: sort,
		Sortable:      true,
		SearchText:    r.URL.Query().Get("text"),
		PUG:           u.PublicURL,
	}, videos).Render(r.Context(), w)
}

func NewWriteableCUI2(publicURL tmpl.PublicURLGenerator, repo videostore.VideoRepo) http.Handler {
	u := &cUI2{
		ReadOnly:  false,
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

	r.HandleFunc(
		"/search",
		u.Search,
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
