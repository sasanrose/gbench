package main

import (
    "os"
    "time"
    "os/signal"
    "fmt"
)

var (
    urls urlsList
    err error
    urlFile *os.File
    requests, concurrent int
    totalTransactions, failedTransactions, delay int
    cookies cookiesList
    responseStats, urlFailedStats map[string]int
    urlsResponseTimes map[string]time.Duration
    totalResponseTime, shortestResponseTime, longestResponseTime, responseTimeout, connectionTimeout time.Duration
    averageResponseTime, transactionRate, transferredData float64
    basicAuthUsername, basicAuthPassword string
    proxyUrl proxyURL
    totalLength int64
    headers headersList
    verbose bool = false
    disableKeepAlive bool = false
)

type Url struct {
    url string
    method string
    data map[string]string
}

type urlsList []Url;
type cookiesList []string;
type headersList map[string]string;
type proxyURL string;

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


    startBench();
    showResult();
}
