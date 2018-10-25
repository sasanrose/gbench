package cmd

import (
	"log"
	"os"

	"github.com/sasanrose/gbench/bench"
	"github.com/spf13/cobra"
)

var execCmd = &cobra.Command{
	Use:   "exec",
	Short: "Executes the benchmark",
	Long: `Executes the benchmark using the given urls.
Sample usage:

gbench exec -r 100 -c 10 www.google.com
gbench exec -X post -d "search=gbench" -r 100 -c 10 www.google.com`,
	Run: runExec,
}

func runExec(cmd *cobra.Command, args []string) {
	if len(args) != 1 {
		cmd.Usage()
		os.Exit(0)
	}

	configurations := make([]func(*bench.Bench), 0)

	urlConfig, err := bench.WithURLSettings(args[0], method, data, []string{}, "", "")

	if err != nil {
		log.Fatalf("Error with url: %v", err)
	}

	configurations = append(configurations, urlConfig)

	runBench(configurations)
}

func init() {
	rootCmd.AddCommand(execCmd)

	initExecFlags()
}
