package web

import (
	"strings"

	"github.com/AlbinoDrought/creamy-videos/videostore"
)

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

	sortDirection := dict.Get("sort_direction")
	sortField := dict.Get("sort_field")

	hasSortDirection := sortDirection != ""
	hasSortField := sortField != ""

	if hasSortField && !hasSortDirection {
		sortDirection = videostore.SortDirectionAscending
	}

	return videostore.VideoFilter{
		Title: dict.Get("title"),
		Tags:  tags,
		Any:   dict.Get("filter"),

		SortDirection: sortDirection,
		SortField:     sortField,
	}
}
