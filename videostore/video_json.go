package videostore

import (
	"bytes"
	"encoding/json"
	"io"
	"path"
	"strconv"
	"sync"

	"github.com/AlbinoDrought/creamy-videos/files"
)

// dummyVideoRepo stores models to a local JSON file
type dummyVideoRepo struct {
	VideoRepo
	fs        files.TransformedFileSystem
	videos    []Video
	id        uint
	idLock    sync.Mutex
	videoLock sync.Mutex
}

func NewDummyVideoRepo(fs files.TransformedFileSystem) *dummyVideoRepo {
	var videos []Video

	storedDatabase, err := fs.Open("dummy.json")
	if err == nil {
		defer storedDatabase.Close()
		err = json.NewDecoder(storedDatabase).Decode(&videos)
		// continue loading
	}

	if err != nil {
		// create new video repo if:
		// - dummy.json not found
		// - failed to load dummy.json
		videos = make([]Video, 0)
	}

	return &dummyVideoRepo{
		fs:     fs,
		videos: videos,
		id:     uint(len(videos)),
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

	rootDir := strconv.Itoa(int(video.ID))
	if _, err := repo.fs.Stat(rootDir); repo.fs.IsNotExist(err) {
		repo.fs.MkdirAll(rootDir, 0600)
	}

	videoPath := path.Join(rootDir, "video"+path.Ext(video.OriginalFileName))

	repo.fs.PipeTo(videoPath, reader)

	video.Source = videoPath
	go eventuallyMakeThumbnail(video, repo, repo.fs)

	return repo.Save(video)
}

func (repo *dummyVideoRepo) Save(video Video) (Video, error) {
	repo.videoLock.Lock()
	defer repo.videoLock.Unlock()

	if !video.Exists() {
		// create
		video.ID = repo.makeID()
		repo.videos = append(repo.videos, video)

		return video, nil
	}

	if len(repo.videos) < int(video.ID) {
		return Video{}, ErrorVideoNotFound
	}

	repo.videos[video.ID-1] = video
	videoJSON, _ := json.Marshal(&repo.videos)
	repo.fs.PipeTo("dummy.json", bytes.NewReader(videoJSON))

	return video, nil
}

func (repo *dummyVideoRepo) FindById(video uint) (Video, error) {
	if len(repo.videos) < int(video) {
		return Video{}, ErrorVideoNotFound
	}

	return repo.videos[int(video)-1], nil
}

func (repo *dummyVideoRepo) All(limit uint, offset uint) ([]Video, error) {
	max := uint(len(repo.videos))

	start := offset
	if start > max {
		start = max
	}

	end := start + limit
	if end > max {
		end = max
	}

	return repo.videos[start:end], nil
}
