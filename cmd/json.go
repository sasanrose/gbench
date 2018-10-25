package cmd

import (
	"encoding/json"
	"fmt"
	"log"
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
	Run: runJson,
}

func runJson(cmd *cobra.Command, args []string) {
	if len(args) != 1 {
		cmd.Usage()
		os.Exit(0)
	}

	file, err := fs.Open(args[0])

	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not open %s: %v.\n", args[0], err)
		os.Exit(2)
	}

	defer file.Close()

	config := &JsonConfig{}
	decoder := json.NewDecoder(file)

	decoder.Decode(config)

	if config.Host == "" {
		fmt.Fprintln(os.Stderr, "No host is provided")
		os.Exit(2)
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
			log.Fatalf("Error with url: %v", err)
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

	runBench(configurations)
}

func init() {
	rootCmd.AddCommand(jsonCmd)

	initSharedFlags(jsonCmd)
}
