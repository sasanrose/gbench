package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/sasanrose/gbench/bench"
	"github.com/sasanrose/gbench/result"
)

func main() {
	if len(urls) == 0 && file == "" {
		fmt.Fprintln(os.Stderr, "You should provide at least a url or a file")
		GbenchUsage()
		os.Exit(2)
	}

	renderer := result.NewStdoutRenderer()
	renderer.Init(concurrent)

	configurations := []func(*bench.Bench){
		bench.WithConcurrency(concurrent),
		bench.WithRequests(requests),
		bench.WithConnectionTimeout(connectionTimeout),
		bench.WithResponseTimeout(responseTimeout),
		bench.WithRenderer(renderer),
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

	renderer.Render()
}
