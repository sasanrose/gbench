package cmd

import (
	"net/http"
	"time"
)

var (
	isVerbose                                             bool
	urls, headers                                         []string
	file, authUsername, authPassword, proxyUrl, rawCookie string
	concurrent, requests                                  int
	successStatusCodes                                    []int
	connectionTimeout, responseTimeout                    time.Duration
)

func initFlags() {
	execCmd.Flags().BoolVarP(&isVerbose, "verbose", "v", false, "Turn on verbosity mode.")
	execCmd.Flags().StringSliceVarP(&urls, "url", "u", []string{}, "Url to benchmark. This can be used multiple times.")
	execCmd.Flags().StringVarP(&file, "file", "f", "", "Path to the file containing list of urls to benchmark.")
	execCmd.Flags().IntVarP(&concurrent, "concurrent", "c", 1, "Number of concurrent requests.")
	execCmd.Flags().IntVarP(&requests, "requests", "r", 1, "Number of requests to send.")
	execCmd.Flags().IntSliceVarP(&successStatusCodes,
		"status-codes",
		"s",
		[]int{http.StatusOK, http.StatusAccepted, http.StatusCreated},
		"Define what should be considered as a successful status code.")
	execCmd.Flags().StringVar(&authUsername, "auth-username", "", "Username for basic HTTP authentication.")
	execCmd.Flags().StringVar(&authPassword, "auth-password", "", "Password for basic HTTP authentication.")
	execCmd.Flags().StringVar(&proxyUrl, "proxy-url", "", "Proxy server url.")
	execCmd.Flags().DurationVarP(&connectionTimeout, "connection-timeout", "C", 0, "Connection timeout (0 means no timeout).")
	execCmd.Flags().DurationVarP(&responseTimeout, "response-timeout", "R", 0, "Response timeout (0 means no timeout).")
	execCmd.Flags().StringSliceVarP(&headers, "header", "H", []string{}, "HTTP header in format of key=value. This can be used multiple times.")
	execCmd.Flags().StringVar(&rawCookie, "raw-cookie", "", "A string to be sent as raw cookie (In the format of Set-Cookie HTTP header).")
}
