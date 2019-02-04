package videostore

import (
	"errors"
	"io"
)

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
	Upload(video Video, reader io.Reader) (Video, error)
	Save(video Video) (Video, error)
	FindById(id uint) (Video, error)
	All(limit uint, offset uint) ([]Video, error)
}

var ErrorVideoNotFound = errors.New("video not found")
