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
)

func runBench(configurations []func(b *bench.Bench)) {
	if _, err := os.Stat(outputPath); err == nil && !forceOverWrite {
		exitWithError(fmt.Sprintf("%s already exists. Use -F to overwrite.", outputPath))
	}

	outputFile, err := os.Create(outputPath)

	if err != nil {
		exitWithError(fmt.Sprintf("Could not open %s: %v\n", outputPath, err))
	}

	defer outputFile.Close()

	result := &report.Result{}
	result.Init(concurrency)

	configurations, err = appendGlobalConfigurations(configurations, result)

	if err != nil {
		exitWithError(err.Error())
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

	log.Printf("Storing the report in %s...", outputPath)
	encoder := json.NewEncoder(outputFile)
	encoder.Encode(result)
}

func appendGlobalConfigurations(configurations []func(*bench.Bench), result *report.Result) ([]func(*bench.Bench), error) {
	configurations = append(configurations, []func(*bench.Bench){
		bench.WithConcurrency(concurrency),
		bench.WithRequests(requests),
		bench.WithConnectionTimeout(connectionTimeout),
		bench.WithResponseTimeout(responseTimeout),
		bench.WithReport(result),
		bench.WithOutput(os.Stdout),
	}...)

	for _, statusCode := range successStatusCodes {
		statusCodeConfig := bench.WithSuccessStatusCode(statusCode)

		configurations = append(configurations, statusCodeConfig)
	}

	for _, header := range headers {
		headerConfig, err := bench.WithHeaderString(header)

		if err != nil {
			return []func(*bench.Bench){}, fmt.Errorf("Error with header: %v", err)
		}

		configurations = append(configurations, headerConfig)
	}

	if authUserPass != "" {
		authConfig, err := bench.WithAuthUserPass(authUserPass)

		if err != nil {
			return []func(*bench.Bench){}, fmt.Errorf("Error with authentication credentials: %v", err)
		}

		configurations = append(configurations, authConfig)
	}

	if proxyURL != "" {
		configurations = append(configurations, bench.WithProxy(proxyURL))
	}

	if rawCookie != "" {
		configurations = append(configurations, bench.WithRawCookie(rawCookie))
	}

	return configurations, nil
}

func exitWithError(msg string) {
	fmt.Fprintf(os.Stderr, msg)
	os.Exit(2)
}
