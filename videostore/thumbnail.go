package videostore

import (
	"log"
	"os/exec"
	"path"

	"github.com/AlbinoDrought/creamy-videos/files"
	"github.com/pkg/errors"
)

func GenerateThumbnail(video Video, repo VideoRepo, fs files.TransformedFileSystem) (Video, error) {
	thumbnailPath := path.Join(path.Dir(video.Source), "thumbnail.jpg")

	videoStream, err := fs.Open(video.Source)
	if err != nil {
		return video, errors.Wrap(err, "failed to open video")
	}
	defer videoStream.Close()

	createdThumbnailStream, err := fs.Create(thumbnailPath)
	if err != nil {
		return video, errors.Wrap(err, "failed to create thumbnail")
	}
	defer createdThumbnailStream.Close()

	cmd := exec.Command("ffmpeg", "-i", "-", "-vf", "thumbnail,scale=640:-1", "-frames:v", "1", "-f", "singlejpeg", "-")
	cmd.Stdin = videoStream
	cmd.Stdout = createdThumbnailStream

	err = cmd.Run()
	if err != nil {
		return video, errors.Wrap(err, "failed to run ffmpeg")
	}

	video.Thumbnail = thumbnailPath
	_, err = repo.Save(video)

	if err != nil {
		return video, errors.Wrap(err, "failed to save video thumbnail to disk")
	}

	return video, nil
}

func eventuallyMakeThumbnail(video Video, repo VideoRepo, fs files.TransformedFileSystem) {
	_, err := GenerateThumbnail(video, repo, fs)
	if err != nil {
		log.Printf("failed to make thumbnail thumbnail: %+v", err)
	}
}
