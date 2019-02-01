package videostore

import (
	"log"
	"os/exec"
	"path"

	"github.com/AlbinoDrought/creamy-videos/files"
)

func eventuallyMakeThumbnail(video Video, repo VideoRepo, fs files.TransformedFileSystem) {
	thumbnailPath := path.Join(path.Dir(video.Source), "thumbnail.jpg")

	videoStream, err := fs.Open(video.Source)
	if err != nil {
		log.Printf("failed to open video: %+v", err)
		return
	}
	defer videoStream.Close()

	createdThumbnailStream, err := fs.Create(thumbnailPath)
	if err != nil {
		log.Printf("failed to create thumbnail: %+v", err)
		return
	}
	defer createdThumbnailStream.Close()

	cmd := exec.Command("ffmpeg", "-i", "-", "-vf", "thumbnail,scale=640:-1", "-frames:v", "1", "-f", "singlejpeg", "-")
	cmd.Stdin = videoStream
	cmd.Stdout = createdThumbnailStream

	err = cmd.Run()
	if err != nil {
		log.Printf("failed to run ffmpeg: %+v", err)
		return
	}

	video.Thumbnail = thumbnailPath
	_, err = repo.Save(video)

	if err != nil {
		log.Printf("failed to save video thumbnail: %+v", err)
		return
	}
}
