package cmd

import (
	"log"
	"strconv"

	"github.com/AlbinoDrought/creamy-videos/videostore"
	"github.com/spf13/cobra"
)

var transcodeAllVideos = false

var transcodeCommand = &cobra.Command{
	Use:   "transcode [-a to regen all] [video ids]",
	Short: "Transcode given video(s), or all videos",
	Run: func(cmd *cobra.Command, args []string) {
		var getter videoGetter
		// loop over all videos
		if transcodeAllVideos {
			currentOffset := uint(0)
			limit := uint(100)
			getter = func() []videostore.Video {
				videos, err := app.repo.All(videostore.VideoFilter{}, limit, currentOffset)
				if err != nil {
					log.Fatalf("error fetching videos: %+v", err)
				}
				currentOffset += limit
				return videos
			}
		} else {
			// loop over selected ids
			ids := make([]uint, len(args))
			for i, arg := range args {
				id, err := strconv.Atoi(arg)
				if err != nil {
					log.Fatalf("error converting to int: %+v", err)
				}
				ids[i] = uint(id)
			}

			getter = func() []videostore.Video {
				if len(ids) == 0 {
					return []videostore.Video{}
				}

				id := ids[0]
				ids = ids[1:]

				video, err := app.repo.FindById(id)
				if err != nil {
					log.Fatalf("error fetching video: %+v", err)
				}

				return []videostore.Video{video}
			}
		}

		var videos []videostore.Video
		for {
			videos = getter()
			if len(videos) == 0 {
				break
			}
			for _, video := range videos {
				log.Printf("transcoding %+v", video.ID)
				err := processTrancodeJob(app, video.ID)
				if err == nil {
					log.Printf("transcoded %+v", video.ID)
				} else {
					log.Printf("failed to transcode %+v: %+v", video.ID, err)
				}
			}
		}
	},
}

func init() {
	transcodeCommand.Flags().BoolVarP(&transcodeAllVideos, "all", "a", false, "if true, transcode _all_ videos")

	rootCmd.AddCommand(transcodeCommand)
}
