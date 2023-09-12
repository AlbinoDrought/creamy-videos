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
	ReadOnly bool
}

func (u *cUI2) Home(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/html")
	tmpl.Home(tmpl.AppState{
		ReadOnly:      u.ReadOnly,
		SortDirection: "",
		Sortable:      false,
		SearchText:    "",
	}).Render(r.Context(), w)
}

func (u *cUI2) Search(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/html")
	sort := r.URL.Query().Get("sort")
	sortDir, ok := sortDirs[sort]
	if !ok {
		sort = defaultSortDir
		sortDir = sortDirs[sort]
	}

	log.Print("todo", sortDir)

	tmpl.Search(tmpl.AppState{
		ReadOnly:      u.ReadOnly,
		SortDirection: sort,
		Sortable:      true,
		SearchText:    r.URL.Query().Get("text"),
	}).Render(r.Context(), w)
}

func NewWriteableCUI2() http.Handler {
	u := &cUI2{
		ReadOnly: false,
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

func NewReadOnlyCUI2() http.Handler {
	u := &cUI2{
		ReadOnly: true,
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
