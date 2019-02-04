package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var app application

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "creamy-videos",
	Short: "The creamiest selfhosted tubesite",
	Long:  "creamy-videos is a self-service self-hosted tubesite that obscures uploaded media.",
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(func() {
		app = makeApp(makeConfig())
	})
}
