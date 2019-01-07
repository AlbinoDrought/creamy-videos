package main

import (
	"errors"
	"io"
	"os"
	"path"
	"strconv"
	"sync"
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

type VideoRepo interface {
	Upload(video Video, reader io.Reader) (Video, error)
	Save(video Video) (Video, error)
	FindById(id uint) (Video, error)
	All(limit uint, offset uint) ([]Video, error)
}

var ErrorVideoNotFound = errors.New("video not found")

type dummyVideoRepo struct {
	VideoRepo
	videos    []Video
	id        uint
	idLock    sync.Mutex
	videoLock sync.Mutex
}

func NewDummyVideoRepo() *dummyVideoRepo {
	return &dummyVideoRepo{
		videos: make([]Video, 32),
		id:     0,
	}
}

func (repo *dummyVideoRepo) makeID() uint {
	repo.idLock.Lock()
	defer repo.idLock.Unlock()
	repo.id = repo.id + 1
	return repo.id
}

func (repo *dummyVideoRepo) Upload(video Video, reader io.Reader) (Video, error) {
	video.Thumbnail = ""
	video.Source = ""
	video, err := repo.Save(video)

	if err != nil {
		return video, err
	}

	rootDir := path.Join(".", "dummyvideos", strconv.Itoa(int(video.ID)))
	if _, err := os.Stat(rootDir); os.IsNotExist(err) {
		os.Mkdir(rootDir, os.ModePerm)
	}

	videoPath := path.Join(rootDir, "video"+path.Ext(video.OriginalFileName))

	file, err := os.Create(videoPath)
	if err != nil {
		return video, err
	}
	defer file.Close()

	_, err = io.Copy(file, reader)
	if err != nil {
		return video, err
	}

	video.Source = videoPath

	return repo.Save(video)
}

func (repo *dummyVideoRepo) Save(video Video) (Video, error) {
	repo.videoLock.Lock()
	defer repo.videoLock.Unlock()

	if video.ID == 0 {
		// create
		video.ID = repo.makeID()
		repo.videos = append(repo.videos, video)

		return video, nil
	}

	if len(repo.videos) < int(video.ID) {
		return Video{}, ErrorVideoNotFound
	}

	repo.videos[video.ID] = video
	return video, nil
}

func (repo *dummyVideoRepo) FindById(video uint) (Video, error) {
	if len(repo.videos) < int(video) {
		return Video{}, ErrorVideoNotFound
	}

	return repo.videos[int(video)], nil
}

func (repo *dummyVideoRepo) All(limit uint, offset uint) ([]Video, error) {
	return repo.videos[offset:(limit + offset)], nil
}
