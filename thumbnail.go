package main

import (
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
)

func makeThumbnailOfExistingFile(existingFile string, outputThumbnail string) error {
	cmd := exec.Command("ffmpeg", "-i", existingFile, "-vframes", "1", outputThumbnail)

	return cmd.Run()
}

func eventuallyMakeThumbnail(video Video) {
	// todo: queue
	dir, err := ioutil.TempDir("", "eventual-thumbnail")
	if err != nil {
		log.Printf("failed to make tempdir for thumbnail: %+v", err)
		return
	}

	defer os.RemoveAll(dir) // clean up

	newThumbnailPath := path.Join(dir, "thumbnail.jpg")

	newVideoPath := path.Join(dir, path.Base(video.Source))
	newVideo, err := os.Create(newVideoPath)
	if err != nil {
		log.Printf("failed to make new video file: %+v", err)
		return
	}
	defer newVideo.Close()

	oldVideo, err := transformedFileSystem.Open(video.Source)
	if err != nil {
		log.Printf("failed to open old video: %+v", err)
		return
	}
	defer oldVideo.Close()

	_, err = io.Copy(newVideo, oldVideo)
	if err != nil {
		log.Printf("failed to copy file: %+v", err)
		return
	}

	err = makeThumbnailOfExistingFile(newVideoPath, newThumbnailPath)
	if err != nil {
		log.Printf("failed to make thumbnail: %+v", err)
		return
	}

	createdThumbnail, err := os.Open(newThumbnailPath)
	if err != nil {
		log.Printf("failed to open created thumbnail: %+v", err)
		return
	}
	defer createdThumbnail.Close()

	// copy thumbnail to be beside video, save video
	oldThumbnailPath := path.Join(path.Dir(video.Source), path.Base(newThumbnailPath))

	err = transformedFileSystem.PipeTo(oldThumbnailPath, createdThumbnail)
	if err != nil {
		log.Printf("failed to make old thumbnail: %+v", err)
		return
	}

	video.Thumbnail = oldThumbnailPath
	_, err = videoRepo.Save(video)

	if err != nil {
		log.Printf("faild to save video thumbnail: %+v", err)
		return
	}
}
