package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "gbench",
	Short: "Gbench, HTTP benchmarking tool",
	Long: `Yet another HTTP benchmarking tool inspired by Apache benchmark and Siege.
For more info please refer to https://github.com/sasanrose/gbench`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Usage()
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
