package cmd

import (
	"log"
	"strconv"

	"github.com/AlbinoDrought/creamy-videos/videostore"
	"github.com/spf13/cobra"
)

var regenerateAllThumbnails = false

type videoGetter func() []videostore.Video

var thumbnailCommand = &cobra.Command{
	Use:   "thumbnail [-a to regen all] [video ids]",
	Short: "Regenerate thumbnail for given video, or all videos",
	Run: func(cmd *cobra.Command, args []string) {
		var getter videoGetter
		// loop over all videos
		if regenerateAllThumbnails {
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
				_, err := videostore.GenerateThumbnail(video, app.repo, app.fs)
				if err == nil {
					log.Printf("generated thumbnail for %+v", video.ID)
				} else {
					log.Printf("failed to generate for %+v: %+v", video.ID, err)
				}
			}
		}
	},
}

func init() {
	thumbnailCommand.Flags().BoolVarP(&regenerateAllThumbnails, "all", "a", false, "if true, regenerate _all_ thumbnails")

	rootCmd.AddCommand(thumbnailCommand)
}
