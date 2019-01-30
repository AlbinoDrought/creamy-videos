package main

import (
	"log"
	"os/exec"
	"path"
)

func eventuallyMakeThumbnail(video Video) {
	thumbnailPath := path.Join(path.Dir(video.Source), "thumbnail.jpg")

	videoStream, err := transformedFileSystem.Open(video.Source)
	if err != nil {
		log.Printf("failed to open video: %+v", err)
		return
	}
	defer videoStream.Close()

	createdThumbnailStream, err := transformedFileSystem.Create(thumbnailPath)
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
	_, err = videoRepo.Save(video)

	if err != nil {
		log.Printf("faild to save video thumbnail: %+v", err)
		return
	}
}
