package videostore

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
	"strconv"

	"github.com/AlbinoDrought/creamy-videos/files"
	"github.com/pkg/errors"
)

func GenerateThumbnail(video Video, repo VideoRepo, fs files.FileSystem) (Video, error) {
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
	if err == nil {
		video.Thumbnail = thumbnailPath
	} else {
		// some formats cannot be processed when being piped (some MOV files)
		// attempt to generate them by saving to disk
		var tempFileError error
		video, tempFileError = generateThumbnailUsingTemporaryFile(video, fs)
		if tempFileError != nil {
			return video, fmt.Errorf("failed to generate thumbnail using any method.\npipe: %+v\ntemp: %+v", err, tempFileError)
		}
		// if tempFileError is nil,
		// we successfully generated a thumbnail
		// with this method :)
	}

	_, err = repo.Save(video)

	if err != nil {
		return video, errors.Wrap(err, "failed to save video thumbnail to disk")
	}

	return video, nil
}

func generateThumbnailUsingTemporaryFile(video Video, fs files.FileSystem) (Video, error) {
	// create a temporary directory to store our junk
	tempDir, err := ioutil.TempDir("", "eventual-thumbnail-"+strconv.Itoa(int(video.ID)))
	if err != nil {
		return video, errors.Wrap(err, "failed to make tempdir for thumbnail generation")
	}
	defer os.RemoveAll(tempDir) // clean up

	// create the file in our temporary directory
	temporaryVideoPath := path.Join(tempDir, path.Base(video.Source))
	temporaryVideoStream, err := os.Create(temporaryVideoPath)
	if err != nil {
		return video, errors.Wrap(err, "failed to create temporary video file")
	}
	defer temporaryVideoStream.Close()

	// open the existing (and possibly remote) video file
	realVideoStream, err := fs.Open(video.Source)
	if err != nil {
		return video, errors.Wrap(err, "failed to open real video file")
	}
	defer realVideoStream.Close()

	// download it to our temporary directory
	_, err = io.Copy(temporaryVideoStream, realVideoStream)
	if err != nil {
		return video, errors.Wrap(err, "failed to download video to temporary path")
	}

	// close both of our streams, the file has been downloaded
	realVideoStream.Close()
	temporaryVideoStream.Close()

	temporaryThumbnailPath := path.Join(tempDir, "thumbnail.jpg")

	// actually generate the thumbnail using ffmpeg
	cmd := exec.Command("ffmpeg", "-i", temporaryVideoPath, "-vf", "thumbnail,scale=640:-1", "-frames:v", "1", "-f", "singlejpeg", temporaryThumbnailPath)
	err = cmd.Run()
	if err != nil {
		return video, errors.Wrap(err, "failed to run ffmpeg")
	}

	// open our saved thumbnail
	temporaryThumbnailStream, err := os.Open(temporaryThumbnailPath)
	if err != nil {
		return video, errors.Wrap(err, "failed to open created temporary thumbnail")
	}
	defer temporaryThumbnailStream.Close()

	// copy thumbnail to be beside video
	finalThumbnailPath := path.Join(path.Dir(video.Source), path.Base(temporaryThumbnailPath))
	err = files.PipeTo(fs, finalThumbnailPath, temporaryThumbnailStream)
	if err != nil {
		return video, errors.Wrap(err, "failed to upload temporary thumbnail")
	}

	video.Thumbnail = finalThumbnailPath

	return video, nil
}

func eventuallyMakeThumbnail(video Video, repo VideoRepo, fs files.FileSystem) {
	_, err := GenerateThumbnail(video, repo, fs)
	if err != nil {
		log.Printf("failed to make thumbnail: %+v", err)
	}
}
