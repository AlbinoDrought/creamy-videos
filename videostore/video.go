package videostore

import (
	"errors"
)

const SortDirectionAscending = "asc"
const SortDirectionDescending = "desc"

const SortFieldTitle = "title"
const SortFieldTimeCreated = "time_created"
const SortFieldTimeUpdated = "time_updated"

var SortFields = []string{
	SortFieldTitle,
	SortFieldTimeCreated,
	SortFieldTimeUpdated,
}

// VideoFilter represents a "filter" used for
// querying stored videos.
// It's eventually shoved into the `WHERE` part of a query,
// or equivalent.
type VideoFilter struct {
	Title string
	Tags  []string
	Any   string

	SortDirection string
	SortField     string
}

func (filter VideoFilter) Sort() bool {
	return len(filter.SortDirection)+len(filter.SortField) > 0
}

func (filter VideoFilter) Empty() bool {
	return len(filter.Title)+len(filter.Tags)+len(filter.Any) == 0
}

func (filter VideoFilter) ValidSortDirection() bool {
	return filter.SortDirection == SortDirectionAscending || filter.SortDirection == SortDirectionDescending
}

func (filter VideoFilter) ValidSortField() bool {
	for _, field := range SortFields {
		if filter.SortField == field {
			return true
		}
	}

	return false
}

type Video struct {
	ID               uint     `json:"id"`
	Title            string   `json:"title"`
	Description      string   `json:"description"`
	Thumbnail        string   `json:"thumbnail"`
	Source           string   `json:"source"`
	OriginalFileName string   `json:"original_file_name"`
	TimeCreated      string   `json:"time_created"`
	TimeUpdated      string   `json:"time_updated"`
	Tags             []string `json:"tags"`
}

func (video Video) Exists() bool {
	return video.ID > 0
}

type VideoRepo interface {
	Save(video Video) (Video, error)
	FindById(id uint) (Video, error)
	All(filter VideoFilter, limit uint, offset uint) ([]Video, error)
	Delete(video Video) error
}

var ErrorVideoNotFound = errors.New("video not found")
