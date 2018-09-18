package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/spf13/pflag"
)

var (
	isVerbose                                             bool
	urls, headers                                         []string
	file, authUsername, authPassword, proxyUrl, rawCookie string
	concurrent, requests                                  int
	successStatusCodes                                    []int
	connectionTimeout, responseTimeout                    time.Duration
)

var flags *pflag.FlagSet

func init() {
	flags = pflag.NewFlagSet("Gbench", pflag.ContinueOnError)
	flags.Usage = GbenchUsage
	flags.BoolVarP(&isVerbose, "verbose", "v", false, "Turn on verbosity mode.")
	flags.StringSliceVarP(&urls, "url", "u", []string{}, "Url to benchmark. This can be used multiple times.")
	flags.StringVarP(&file, "file", "f", "", "Path to the file containing list of urls to benchmark.")
	flags.IntVarP(&concurrent, "concurrent", "c", 1, "Number of concurrent requests.")
	flags.IntVarP(&requests, "requests", "r", 1, "Number of requests to send.")
	flags.IntSliceVarP(&successStatusCodes,
		"status-codes",
		"s",
		[]int{http.StatusOK, http.StatusAccepted, http.StatusCreated},
		"Define what should be considered as a successful status code.")
	flags.StringVar(&authUsername, "auth-username", "", "Username for basic HTTP authentication.")
	flags.StringVar(&authPassword, "auth-password", "", "Password for basic HTTP authentication.")
	flags.StringVar(&proxyUrl, "proxy-url", "", "Proxy server url.")
	flags.DurationVarP(&connectionTimeout, "connection-timeout", "C", 0, "Connection timeout (0 means no timeout).")
	flags.DurationVarP(&responseTimeout, "response-timeout", "R", 0, "Response timeout (0 means no timeout).")
	flags.StringSliceVarP(&headers, "header", "H", []string{}, "HTTP header in format of key=value. This can be used multiple times.")
	flags.StringVar(&rawCookie, "raw-cookie", "", "A string to be sent as raw cookie (In the format of Set-Cookie HTTP header).")

	err := flags.Parse(os.Args[1:])

	if err != nil {
		if err == pflag.ErrHelp {
			os.Exit(0)
		}

		fmt.Fprintln(os.Stderr, err)
		GbenchUsage()
		os.Exit(2)
	}
}

func GbenchUsage() {
	fmt.Fprintln(os.Stderr, "Usage of Gbench:")
	flags.PrintDefaults()
}
