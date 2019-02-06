package videostore

import (
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"strconv"

	"github.com/AlbinoDrought/creamy-videos/files"
	"github.com/pkg/errors"
)

func Transcode(video Video, fs files.TransformedFileSystem) (Video, error) {
	// create a temporary directory to store our junk
	tempDir, err := ioutil.TempDir("", "eventual-transcode-"+strconv.Itoa(int(video.ID)))
	if err != nil {
		return video, errors.Wrap(err, "failed to make tempdir for transcode")
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

	temporaryTranscodedPath := path.Join(tempDir, "transcoded.mp4")

	// actually transcode video using ffmpeg
	cmd := exec.Command("ffmpeg", "-i", temporaryVideoPath, "-vcodec", "h264", "-acodec", "aac", "-strict", "-2", temporaryTranscodedPath)
	err = cmd.Run()
	if err != nil {
		return video, errors.Wrap(err, "failed to run ffmpeg")
	}

	// open our transcoded video
	temporaryTranscodedStream, err := os.Open(temporaryTranscodedPath)
	if err != nil {
		return video, errors.Wrap(err, "failed to open created temporary transcoded video")
	}
	defer temporaryTranscodedStream.Close()

	// copy transcoded video to be beside video
	finalTranscodedPath := path.Join(path.Dir(video.Source), path.Base(temporaryTranscodedPath))
	err = fs.PipeTo(finalTranscodedPath, temporaryTranscodedStream)
	if err != nil {
		return video, errors.Wrap(err, "failed to upload temporary transcoded video")
	}

	video.Source = finalTranscodedPath

	return video, nil
}
