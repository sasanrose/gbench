package cmd

var (
	data   []string
	method string
)

func initExecFlags() {
	initSharedFlags(execCmd)

	execCmd.Flags().IntVarP(&concurrency, "concurrency", "c", defaultConcurreny, "Number of concurrent requests.")
	execCmd.Flags().IntVarP(&requests, "total-requests", "r", defaultRequests, "Number of total requests to send.")
	execCmd.Flags().IntSliceVarP(&successStatusCodes,
		"status-codes",
		"s",
		defaultStatusCodes,
		"Define what should be considered as a successful status code.")
	execCmd.Flags().StringVarP(&authUserPass, "user", "u", "", `Specify the user name and password to use for server authentication in the format of user:password. Currently only supports Basic Auth.
The user name and passwords are split up on the first colon, as a result it is impossible to use a colon in the user name.`)
	execCmd.Flags().StringVar(&proxyURL, "proxy", "", "HTTP proxy.")
	execCmd.Flags().DurationVarP(&connectionTimeout, "connect-timeout", "", 0, "Connection timeout (0 means no timeout).")
	execCmd.Flags().DurationVarP(&responseTimeout, "response-timeout", "", 0, "Response timeout (0 means no timeout).")
	execCmd.Flags().StringVarP(&method, "request", "X", defaultMethod, "Specify a custom HTTP method.")
	execCmd.Flags().StringSliceVarP(&headers, "header", "H", []string{}, "HTTP header in format of 'key: value' or 'key: value;' or 'key;'. This can be used multiple times.")
	execCmd.Flags().StringSliceVarP(&data, "data", "d", []string{}, "Sends the specified data in a request. The format should be 'key=val' or 'key1=val1&key2=val2'. This can be used multiple times.")
	execCmd.Flags().StringVarP(&rawCookie, "cookie", "b", "", "A string to be sent as raw cookie (In the format of Set-Cookie HTTP header).")
}
