package cmd

import (
	"log"

	"github.com/AlbinoDrought/creamy-videos/videostore"

	"github.com/pkg/errors"
)

var transcodeQueue chan uint

func init() {
	transcodeQueue = make(chan uint)
}

func processTrancodeJob(instance application, id uint) error {
	video, err := instance.repo.FindById(id)
	if err != nil {
		return errors.Wrap(err, "unable to find video to transcode")
	}

	video, err = videostore.Transcode(video, instance.fs)
	if err != nil {
		return errors.Wrap(err, "unable to transcode video")
	}

	freshVideo, err := instance.repo.FindById(id)
	if err != nil {
		return errors.Wrap(err, "unable to find video after transcode")
	}

	// race condition here (dataloss)
	// may overwrite changes made from other areas
	// (thumbnail gen, editing)
	// todo: fix
	freshVideo.Source = video.Source
	freshVideo, err = instance.repo.Save(freshVideo)

	if err != nil {
		return errors.Wrap(err, "unable to update video source after transcode")
	}

	return nil
}

func workTrancodeQueue(instance application) {
	for {
		select {
		case id := <-transcodeQueue:
			log.Printf("transcoding %+v", id)
			err := processTrancodeJob(instance, id)
			if err == nil {
				log.Printf("transcoded %+v", id)
			} else {
				log.Printf("failed to transocde %+v: %+v", id, err)
			}
		}
	}
}

func eventuallyTranscode(id uint) {
	transcodeQueue <- id
}
