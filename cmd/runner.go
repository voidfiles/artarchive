package cmd

import (
	"log"

	"github.com/spf13/cobra"
	"github.com/voidfiles/artarchive/doers"
)

func init() {
	rootCmd.AddCommand(runnerCmd)
}

var runnerCmd = &cobra.Command{
	Use:   "runner",
	Short: "Indexes a set of feeds into s3",
	Long:  `Indexes a set of feeds into s3`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Print("cmd: runner")
		doers.FeedRunner()
	},
}
