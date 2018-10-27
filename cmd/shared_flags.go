package cmd

import "github.com/spf13/cobra"

var (
	forceOverWrite bool
	outputPath     string
)

func initSharedFlags(cmd *cobra.Command) {
	cmd.Flags().BoolVarP(&forceOverWrite, "force", "F", false, "Force overwrite for the report file.")
	cmd.Flags().StringVarP(&outputPath, "output", "o", "./report.json", "The path to store the report of benchmark.")
}
