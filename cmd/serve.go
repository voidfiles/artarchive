package cmd

import (
	_ "github.com/heroku/x/hmetrics/onload" // Keeps track of go stats
	"github.com/spf13/cobra"
	"github.com/voidfiles/artarchive/server"
)

func init() {
	rootCmd.AddCommand(serveCmd)
}

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Runs an http server",
	Long:  `Runs an http server`,
	Run: func(cmd *cobra.Command, args []string) {
		server.Serve()
	},
}
