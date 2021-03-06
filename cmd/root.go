// Package cmd contains all the subcommands used by gbench.
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "gbench",
	Short: "Gbench, HTTP benchmarking and load generating tool",
	Long: `Yet another HTTP benchmarking tool inspired by Apache benchmark and Siege.
For more info please refer to https://github.com/sasanrose/gbench`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Usage()
	},
}

// Execute runs the gbench root command.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
