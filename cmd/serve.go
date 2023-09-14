package cmd

import (
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/mux"

	"github.com/AlbinoDrought/creamy-videos/files"
	"github.com/AlbinoDrought/creamy-videos/web"
	"github.com/spf13/cobra"
)

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Provide videos, UI, and API over HTTP",
	Run: func(cmd *cobra.Command, args []string) {
		fileServer := http.FileServer(files.AdaptToHTTPFileSystem(app.fs, false))

		r := mux.NewRouter()

		// mount api:
		publicUrlGenerator := func(relativeURL string) string {
			return app.config.AppURL + app.config.HTTPVideoDirectory + relativeURL
		}
		var apiHandler http.Handler
		if app.config.ReadOnly {
			apiHandler = web.NewReadOnlyAPI(publicUrlGenerator, app.fs, app.repo)
		} else {
			apiHandler = web.NewWriteableAPI(publicUrlGenerator, app.fs, app.repo)
		}
		r.PathPrefix("/api/").Handler(apiHandler)

		// mount video files:
		r.PathPrefix(app.config.HTTPVideoDirectory).Handler(
			http.StripPrefix(
				strings.TrimRight(app.config.HTTPVideoDirectory, "/"),
				fileServer,
			),
		)

		// there was previously an SPA UI @ commit 7c6cf4e199ef3baa03dde8cbe16ced4e21251be8
		// mount non-SPA UI:
		var cUI2Handler http.Handler
		if app.config.ReadOnly {
			cUI2Handler = web.NewReadOnlyCUI2(publicUrlGenerator, app.repo)
		} else {
			cUI2Handler = web.NewWriteableCUI2(publicUrlGenerator, app.fs, app.repo, app.config.XSRFKey)
		}
		r.PathPrefix("/").Handler(cUI2Handler)

		http.Handle("/", r)

		log.Printf("Remote URL: %s\n", app.config.AppURL)
		log.Printf("Serving videos from %s on %s\n", app.config.LocalVideoDirectory, app.config.HTTPVideoDirectory)
		log.Printf("Listening on %s\n", app.config.Port)
		log.Fatal(http.ListenAndServe(":"+app.config.Port, nil))
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
}
