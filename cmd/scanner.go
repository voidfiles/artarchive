package cmd

import (
	"log"

	"github.com/spf13/cobra"
	"github.com/voidfiles/artarchive/doers"
)

var renderer string

func init() {
	rootCmd.AddCommand(scannerCmd)
	scannerCmd.Flags().StringVarP(&renderer, "renderer", "r", "slideshow", "What renderer")
}

var scannerCmd = &cobra.Command{
	Use:   "scanner",
	Short: "Turns items in s3 into a slideshow",
	Long:  `Turns items in s3 into a slideshow`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Print("cmd: scanner")
		doers.RunScanner(renderer)
	},
}
