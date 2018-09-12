package cmd

import (
	"log"

	_ "github.com/lib/pq" // For sqlx
	"github.com/spf13/cobra"
	"github.com/voidfiles/artarchive/doers"
)

func init() {
	rootCmd.AddCommand(syncCmd)
}

var syncCmd = &cobra.Command{
	Use:   "sync",
	Short: "Sync slide stores",
	Long:  `Sync slide stores`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Print("cmd: sync")
		doers.RunSlideSync()
	},
}
