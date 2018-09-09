package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "aarchive",
	Short: "aarchive is an art archive ",
	Long:  `aarchive is an art archive `,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func Execute(lambda string) {
	if lambda != "" {
		rootCmd.SetArgs([]string{lambda})
	}
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
