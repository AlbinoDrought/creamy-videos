package videostore

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/AlbinoDrought/creamy-videos/files"
	"github.com/pkg/errors"
)

// dummyVideoRepo stores models to a local JSON file
type dummyVideoRepo struct {
	VideoRepo
	fs        files.FileSystem
	videos    []Video
	id        uint
	idLock    sync.Mutex
	videoLock sync.Mutex
}

func NewDummyVideoRepo(fs files.FileSystem) *dummyVideoRepo {
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
		repo.fs.MkdirAll(rootDir, os.ModePerm)
	}

	videoPath := path.Join(rootDir, "video"+path.Ext(video.OriginalFileName))

	files.PipeTo(repo.fs, videoPath, reader)

	video.Source = videoPath
	go eventuallyMakeThumbnail(video, repo, repo.fs)

	return repo.Save(video)
}

func (repo *dummyVideoRepo) dumpToDisk() {
	videoJSON, _ := json.Marshal(&repo.videos)
	files.PipeTo(repo.fs, "dummy.json", bytes.NewReader(videoJSON))
}

func (repo *dummyVideoRepo) Save(video Video) (Video, error) {
	repo.videoLock.Lock()
	defer repo.videoLock.Unlock()

	if !video.Exists() {
		// create
		video.ID = repo.makeID()
		repo.videos = append(repo.videos, video)
		video.TimeCreated = time.Now().Format(time.RFC3339)
		video.TimeUpdated = time.Now().Format(time.RFC3339)

		return video, nil
	}

	if len(repo.videos) < int(video.ID) {
		return Video{}, ErrorVideoNotFound
	}

	video.TimeUpdated = time.Now().Format(time.RFC3339)
	repo.videos[video.ID-1] = video
	repo.dumpToDisk()

	return video, nil
}

func (repo *dummyVideoRepo) Delete(video Video) error {
	// we can't actually delete videos because of
	// the way we store them :'(
	repo.videoLock.Lock()
	defer repo.videoLock.Unlock()

	index := video.ID - 1
	// soft delete
	repo.videos[index] = Video{}
	repo.dumpToDisk()

	_, err := repo.fs.Stat(video.Source)
	if !repo.fs.IsNotExist(err) {
		// video exists, attempt to delete
		err = repo.fs.Remove(video.Source)
		if err != nil {
			log.Print(errors.Wrap(err, "failed to remove video from disk"))
		}
	}

	_, err = repo.fs.Stat(video.Thumbnail)
	if !repo.fs.IsNotExist(err) {
		// thumbnail exists, attempt to delete
		err = repo.fs.Remove(video.Thumbnail)
		if err != nil {
			log.Print(errors.Wrap(err, "failed to remove thumbnail from disk"))
		}
	}

	return nil
}

func (repo *dummyVideoRepo) FindById(video uint) (Video, error) {
	if len(repo.videos) < int(video) {
		return Video{}, ErrorVideoNotFound
	}

	// ignore soft deleted videos
	videoInstance := repo.videos[int(video)-1]
	if !videoInstance.Exists() {
		return Video{}, ErrorVideoNotFound
	}

	return videoInstance, nil
}

func (repo *dummyVideoRepo) limitVideoSlice(videos []Video, limit uint, offset uint) []Video {
	max := uint(len(videos))

	start := offset
	if start > max {
		start = max
	}

	end := start + limit
	if end > max {
		end = max
	}

	return videos[start:end]
}

func videoHasAllTags(video Video, tags []string) bool {
	if len(video.Tags) == 0 {
		return false
	}

	for _, tag := range tags {
		hasTag := false
		for _, videoTag := range video.Tags {
			if videoTag == tag {
				hasTag = true
				break
			}
		}

		if !hasTag {
			return false
		}
	}

	return true
}

func (repo *dummyVideoRepo) All(filter VideoFilter, limit uint, offset uint) ([]Video, error) {
	var videos []Video

	if filter.Empty() {
		videos = repo.videos
	} else {
		videos = make([]Video, 0)
		// a very inefficient filter
		// accepting PRs ;)
		for _, video := range repo.videos {
			if len(filter.Title) > 0 && strings.Contains(video.Title, filter.Title) {
				videos = append(videos, video)
				continue
			}

			if len(filter.Tags) > 0 && videoHasAllTags(video, filter.Tags) {
				videos = append(videos, video)
				continue
			}

			if len(filter.Any) > 0 {
				if strings.Contains(video.Title, filter.Any) || videoHasAllTags(video, []string{filter.Any}) {
					videos = append(videos, video)
					continue
				}
			}
		}
	}

	// filter soft-deleted videos
	existingVideos := make([]Video, 0)
	for _, video := range videos {
		if video.Exists() {
			existingVideos = append(existingVideos, video)
		}
	}

	if filter.Sort() {
		var sortFunction func(i, j int) bool

		if filter.SortField == SortFieldTitle {
			sortFunction = func(i, j int) bool {
				return strings.Compare(existingVideos[i].Title, existingVideos[j].Title) < 0
			}
		} else if filter.SortField == SortFieldTimeCreated {
			sortFunction = func(i, j int) bool {
				iTime, _ := time.Parse(time.RFC3339, existingVideos[i].TimeCreated)
				jTime, _ := time.Parse(time.RFC3339, existingVideos[i].TimeCreated)

				return iTime.Before(jTime)
			}
		} else if filter.SortField == SortFieldTimeUpdated {
			sortFunction = func(i, j int) bool {
				iTime, _ := time.Parse(time.RFC3339, existingVideos[i].TimeUpdated)
				jTime, _ := time.Parse(time.RFC3339, existingVideos[i].TimeUpdated)

				return iTime.Before(jTime)
			}
		} else {
			return []Video{}, fmt.Errorf("unsupported sort field %v", filter.SortField)
		}

		if filter.SortDirection == SortDirectionDescending {
			oldSortFunction := sortFunction
			sortFunction = func(i, j int) bool {
				return !oldSortFunction(i, j)
			}
		} else if filter.SortDirection != SortDirectionAscending {
			return []Video{}, fmt.Errorf("unsupported sort direction %v", filter.SortDirection)
		}

		sort.Slice(existingVideos, sortFunction)
	}

	return repo.limitVideoSlice(existingVideos, limit, offset), nil
}
