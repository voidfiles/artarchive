package cmd

import (
	"log"

	"github.com/spf13/cobra"
	"github.com/voidfiles/artarchive/doers"
)

func init() {
	rootCmd.AddCommand(scannerCmd)
}

var scannerCmd = &cobra.Command{
	Use:   "scanner",
	Short: "Turns items in s3 into a slideshow",
	Long:  `Turns items in s3 into a slideshow`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Print("cmd: scanner")
		doers.RunScanner()
	},
}
