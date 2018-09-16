package cmd

import (
	"log"

	_ "github.com/lib/pq" // For sqlx
	"github.com/spf13/cobra"
	"github.com/voidfiles/artarchive/doers"
)

func init() {
	rootCmd.AddCommand(cronCmd)
}

var cronCmd = &cobra.Command{
	Use:   "cron",
	Short: "Run scheduled tasks",
	Long:  `Run scheduled tasks`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Print("cmd: cron")
		doers.RunCron()
	},
}
