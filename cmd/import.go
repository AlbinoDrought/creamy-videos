package cmd

import (
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"strings"

	"github.com/AlbinoDrought/creamy-videos/videostore"
	"github.com/BrianAllred/goydl"
	"github.com/imroc/req"
	"github.com/spf13/cobra"
)

var importVideo = videostore.Video{}
var applyImportTag = true
var creamyVideosHostString = "http://localhost:3000"

const CreamyVideosAPIUploadVideo = "/api/upload"

// importCmd represents the dejson command
var importCmd = &cobra.Command{
	Use:   "import [video url]",
	Short: "Import video into creamy-videos",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		creamyVideosHost, err := url.Parse(creamyVideosHostString)

		if err != nil {
			log.Fatalf("failed to parse url %+v: %+v", creamyVideosHostString, err)
		}

		creamyVideosHost.Path = CreamyVideosAPIUploadVideo

		if applyImportTag {
			importVideo.Tags = append(importVideo.Tags, "import")
		}

		needsTitle := importVideo.Title == ""
		needsDescription := importVideo.Description == ""

		youtubeDl := goydl.NewYoutubeDl()
		youtubeDl.Options.NoPlaylist.Value = true
		youtubeDl.VideoURL = args[0]
		info, err := youtubeDl.GetInfo()

		if err != nil {
			log.Fatalf("error fetching info: %+v", err)
		}

		if needsTitle {
			importVideo.Title = info.Title
		}

		if needsDescription {
			importVideo.Description = info.Description
		}

		filename := info.Filename
		if filename == "" {
			filename = "video"
		}

		youtubeDl = goydl.NewYoutubeDl()
		youtubeDl.Options.Output.Value = "-" // stdout

		done := make(chan error, 1)

		downloadCmd, err := youtubeDl.Download(args[0])
		if err != nil {
			log.Fatalf("error running youtubedl: %+v", err)
		}
		defer downloadCmd.Process.Kill()

		go func() {
			r := req.New()

			log.Printf("streaming upload to %+v", creamyVideosHost.String())
			r.SetTimeout(0)
			resp, err := r.Post(creamyVideosHost.String(), req.Param{
				"title":       importVideo.Title,
				"description": importVideo.Description,
				"tags":        strings.Join(importVideo.Tags, ","),
			}, req.FileUpload{
				File:      youtubeDl.Stdout,
				FieldName: "file",
				FileName:  filename,
			})

			if err != nil {
				done <- err
				return
			}

			if resp.Response().StatusCode >= 400 {
				done <- fmt.Errorf("bad status code: %+v", resp.Response().StatusCode)
				return
			}

			err = json.NewDecoder(resp.Response().Body).Decode(&importVideo)
			done <- err
		}()

		err = <-done

		if err != nil {
			log.Fatalf("error uploading video: %+v", err)
		}

		creamyVideosHost.Path = fmt.Sprintf("/watch/%v", importVideo.ID)
		log.Printf("imported %+v: %+v", args[0], creamyVideosHost.String())
	},
}

func init() {
	rootCmd.AddCommand(importCmd)
	importCmd.Flags().StringVarP(&creamyVideosHostString, "api-url", "u", "http://localhost:3000", "url of creamy-videos API")
	importCmd.Flags().StringVarP(&importVideo.Title, "title", "t", "", "new title of video")
	importCmd.Flags().StringVarP(&importVideo.Description, "description", "d", "", "new description of video")
	importCmd.Flags().StringArrayVar(&importVideo.Tags, "tags", []string{}, "new tags of video")
	importCmd.Flags().BoolVarP(&applyImportTag, "tag-as-import", "i", true, "add an \"import\" tag")
}
