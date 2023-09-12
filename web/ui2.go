package web

import (
	"net/http"

	"github.com/AlbinoDrought/creamy-videos/ui2/static"
	"github.com/AlbinoDrought/creamy-videos/ui2/tmpl"
	"github.com/gorilla/mux"
)

type CreamyVideosUI2 interface {
	Home(w http.ResponseWriter, r *http.Request)
	Search(w http.ResponseWriter, r *http.Request)

	// todo: Upload, Show, Edit, Delete UI & Handler routes
}

type cUI2 struct {
	ReadOnly bool
}

func (u *cUI2) Home(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/html")
	tmpl.Home(tmpl.AppState{
		ReadOnly:      false,
		SortDirection: "",
		Sortable:      false,
		SearchText:    "",
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
