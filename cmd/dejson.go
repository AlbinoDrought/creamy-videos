package cmd

import (
	"log"

	"github.com/spf13/cobra"
)

// dejsonCmd represents the dejson command
var dejsonCmd = &cobra.Command{
	Use:   "dejson",
	Short: "Migrate data from the JSON video repository to the active repository",
	Run: func(cmd *cobra.Command, args []string) {
		dummyRepo := app.makeDummyRepo()

		videos, err := dummyRepo.All(10000, 0)
		if err != nil {
			log.Fatalf("error fetching all videos from JSON repo: %+v", err)
		}

		for _, video := range videos {
			videoID := video.ID
			video.ID = 0 // trigger a re-save instead of update
			savedVideo, err := app.repo.Save(video)
			if err != nil {
				log.Fatalf("error saving video %+v: %+v", video, err)
			}
			video.ID = videoID
			log.Printf("saved video %+v as %+v", video, savedVideo)
		}
	},
}

func init() {
	rootCmd.AddCommand(dejsonCmd)
}
