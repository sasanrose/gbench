package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/sasanrose/gbench/bench"
	"github.com/spf13/cobra"
)

var jsonCmd = &cobra.Command{
	Use:   "json",
	Short: "Executes the benchmark using json configuration",
	Long: `Executes the benchmark using a given json configuration.
Sample usage:

gbench json config.json`,
	Run: runJSON,
}

func runJSON(cmd *cobra.Command, args []string) {
	if len(args) != 1 {
		cmd.Usage()
		os.Exit(2)
	}

	configurations, err := getJSONConfig(args[0])

	if err != nil {
		exitWithError(err.Error())
	}

	runBench(configurations)
}

func getJSONConfig(filePath string) ([]func(*bench.Bench), error) {
	file, err := fs.Open(filePath)

	if err != nil {
		return []func(*bench.Bench){}, fmt.Errorf("Could not open %q: %v", filePath, err)
	}

	defer file.Close()

	config := &JSONConfig{}
	decoder := json.NewDecoder(file)

	decoder.Decode(config)

	if config.Host == "" {
		return []func(*bench.Bench){}, errors.New("No host is provided")
	}

	if len(config.Paths) == 0 {
		return []func(*bench.Bench){}, errors.New("No path is provided")
	}

	config.Host = strings.TrimRight(config.Host, "/?&")

	configurations := make([]func(*bench.Bench), 0)

	for _, path := range config.Paths {
		URL := config.Host + "/" + strings.TrimLeft(path.Path, "/")
		urlConfig, err := bench.WithURLSettings(URL,
			path.Method,
			path.Data,
			path.Headers,
			path.RawCookie,
			path.AuthUserPass)

		if err != nil {
			return []func(*bench.Bench){}, fmt.Errorf("Error with url: %v", err)
		}

		configurations = append(configurations, urlConfig)
	}

	if len(config.StatusCodes) == 0 {
		config.StatusCodes = defaultStatusCodes
	}

	if config.Concurrency == 0 {
		config.Concurrency = defaultConcurreny
	}

	if config.Requests == 0 {
		config.Requests = defaultRequests
	}

	successStatusCodes = config.StatusCodes
	concurrency = config.Concurrency
	requests = config.Requests
	headers = config.Headers
	authUserPass = config.AuthUserPass
	proxyURL = config.Proxy
	rawCookie = config.RawCookie

	return configurations, nil
}

func init() {
	rootCmd.AddCommand(jsonCmd)

	initSharedFlags(jsonCmd)
}
