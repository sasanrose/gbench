package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/sasanrose/gbench/bench"
	"github.com/sasanrose/gbench/report"
	"github.com/spf13/cobra"
)

var execCmd = &cobra.Command{
	Use:   "exec",
	Short: "Executes the benchmark",
	Long: `Executes the benchmark using the given urls.
URL format should be as follow:

METHOD|URL|POSTDATA or METHOD|URL

Sample URLs: GET|www.google.com?search=test or POST|www.google.com|search=test or HEAD|www.google.com

Examples:

gbench exec -file ~/benchmarkurl.txt -r 100 -c 10 -v
gbench exec -url 'GET|www.google.com' -url 'GET|www.google.com/path2' -r 100 -c 10 -v`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(urls) == 0 && file == "" {
			fmt.Fprintln(os.Stderr, "You should provide at least a url or a file")
			cmd.Usage()
			os.Exit(2)
		}

		if _, err := os.Stat(outputPath); err == nil && !forceOverWrite {
			fmt.Fprintf(os.Stderr, "%s already exists. Use -F to overwrite.\n", outputPath)
			os.Exit(2)
		}

		outputFile, err := os.Create(outputPath)

		if err != nil {
			fmt.Fprintf(os.Stderr, "Could not open %s: %v\n", outputPath, err)
			os.Exit(2)
		}

		defer outputFile.Close()

		result := &report.Result{}
		result.Init(concurrent)

		configurations := []func(*bench.Bench){
			bench.WithConcurrency(concurrent),
			bench.WithRequests(requests),
			bench.WithConnectionTimeout(connectionTimeout),
			bench.WithResponseTimeout(responseTimeout),
			bench.WithReport(result),
		}

		if file != "" {
			fileConfig, err := bench.WithFile(file)

			if err != nil {
				log.Fatalf("Error reading file: %v", err)
			}

			configurations = append(configurations, fileConfig)
		}

		for _, statusCode := range successStatusCodes {
			statusCodeConfig := bench.WithSuccessStatusCode(statusCode)

			configurations = append(configurations, statusCodeConfig)
		}

		for _, url := range urls {
			urlConfig, err := bench.WithURLString(url)

			if err != nil {
				log.Fatalf("Error with url: %v", err)
			}

			configurations = append(configurations, urlConfig)
		}

		for _, header := range headers {
			headerConfig, err := bench.WithHeaderString(header)

			if err != nil {
				log.Fatalf("Error with header: %v", err)
			}

			configurations = append(configurations, headerConfig)
		}

		if isVerbose {
			configurations = append(configurations, bench.WithVerbosity(os.Stdout))
		}

		if authUsername != "" && authPassword != "" {
			configurations = append(configurations, bench.WithAuth(authUsername, authPassword))
		}

		if proxyUrl != "" {
			configurations = append(configurations, bench.WithProxy(proxyUrl))
		}

		if rawCookie != "" {
			configurations = append(configurations, bench.WithRawCookie(rawCookie))
		}

		ctx, cancelFunc := context.WithCancel(context.Background())
		sigs := make(chan os.Signal, 1)

		signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

		go func() {
			sig := <-sigs
			log.Printf("Got signal %v. Stopping the benchmark...", sig)
			cancelFunc()
		}()

		b := bench.NewBench(configurations...)
		b.Exec(ctx)

		encoder := json.NewEncoder(outputFile)
		encoder.Encode(result)
	},
}

func init() {
	rootCmd.AddCommand(execCmd)

	initExecFlags()
}
