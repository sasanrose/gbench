package main

import (
    "os"
    "time"
    "os/signal"
    "fmt"
    "net/url"
)

var (
    urls []Url
    err error
    urlFile *os.File
    requests, concurrent int = 0, 1
    totalTransactions, failedTransactions, delay int
    cookies []string
    responseStats, urlFailedStats map[string]int
    urlsResponseTimes map[string]time.Duration
    totalResponseTime, shortestResponseTime, longestResponseTime, responseTimeout, connectionTimeout time.Duration
    averageResponseTime, transactionRate, transferredData float64
    basicAuthUsername, basicAuthPassword string
    proxyUrl *url.URL
    totalLength int64
    headers map[string]string
    verbose bool = false
    disableKeepAlive bool = false
)

type Url struct {
    url string
    method string
    data map[string]string
}

// Main func
func main() {
    interruptSignal := make(chan os.Signal, 1);
    signal.Notify(interruptSignal, os.Interrupt, os.Kill);

    go func() {
        for sig := range interruptSignal {
            showResult();
            fmt.Printf("Got Signal: %v\n", sig);
            os.Exit(1);
        }
    }()

    parsArgs();
    startBench();
    showResult();
}
