package cmd

import (
	"fmt"
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
		os.Exit(2)
	}

	configurations, err := getExecConfig(args[0])

	if err != nil {
		exitWithError(err.Error())
	}

	runBench(configurations)
}

func getExecConfig(url string) ([]func(*bench.Bench), error) {
	configurations := make([]func(*bench.Bench), 0)

	urlConfig, err := bench.WithURLSettings(url, method, data, []string{}, "", "")

	if err != nil {
		return []func(*bench.Bench){}, fmt.Errorf("Error with url: %v", err)
	}

	configurations = append(configurations, urlConfig)

	return configurations, nil
}

func init() {
	rootCmd.AddCommand(execCmd)

	initExecFlags()
}
