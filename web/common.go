package web

import (
	"net/http"
	"strconv"
)

func page(r *http.Request) (int, error) {
	page := r.URL.Query().Get("page")
	if len(page) <= 0 {
		page = "1"
	}
	return strconv.Atoi(page)
}
