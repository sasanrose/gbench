package bench

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/sasanrose/gbench/report"
)

const URL_ERROR_MESSAGE = "%s\nWrong URL format. Example: GET|www.google.com?search=test or POST|www.google.com|search=test or HEAD|www.google.com"

// Create a config to set number of concurrent requests per URL.
func WithConcurrency(n int) func(*Bench) {
	return func(b *Bench) {
		b.Concurrency = n
	}
}

// Create a config to set total number of requests to send per URL.
func WithRequests(n int) func(*Bench) {
	return func(b *Bench) {
		b.Requests = n
	}
}

// Add an endpoint to benchmark.
func WithURL(u *URL) func(*Bench) {
	return func(b *Bench) {
		b.URLs = append(b.URLs, u)
	}
}

// Defines what should be considered as a success status code.
// Default values are: 200, 201, 202
func WithSuccessStatusCode(code int) func(*Bench) {
	return func(b *Bench) {
		b.SuccessStatusCodes = append(b.SuccessStatusCodes, code)
	}
}

// Add basic HTTP authentication. Note: This will be used for all the provided urls.
func WithAuth(username, password string) func(*Bench) {
	return func(b *Bench) {
		b.Auth = &Auth{username, password}
	}
}

// Add an endpoint to benchmark.
func WithVerbosity(w io.Writer) func(*Bench) {
	return func(b *Bench) {
		b.VerbosityWriter = w
	}
}

// Proxy server address. Note: This will be used for all the provided urls.
func WithProxy(addr string) func(*Bench) {
	return func(b *Bench) {
		b.Proxy = addr
	}
}

// Set a benchmarking endpoint using a string.
// Supported formats are:
//
// GET|http://www.google.com?search=test
// POST|https://www.google.com|search=test
// POST|https://www.google.com?query=string|search=test&foo=bar
// HEAD|https://www.google.com
func WithURLString(u string) (func(*Bench), error) {
	parsedURL, err := parseURL(u)

	if err != nil {
		return nil, err
	}

	f := WithURL(parsedURL)

	return f, nil
}

// Set connection timeout.
func WithConnectionTimeout(t time.Duration) func(*Bench) {
	return func(b *Bench) {
		b.ConnectionTimeout = t
	}
}

// Set response timeout.
func WithResponseTimeout(t time.Duration) func(*Bench) {
	return func(b *Bench) {
		b.ResponseTimeout = t
	}
}

// Set a raw cookie string. Note: This will be used for all the provided urls.
func WithRawCookie(cookie string) func(*Bench) {
	return func(b *Bench) {
		b.RawCookie = cookie
	}
}

// Set http headers. Note: This will be used for all the provided urls.
func WithHeader(key, value string) func(*Bench) {
	return func(b *Bench) {
		b.Headers[key] = value
	}
}

// Set http headers using a string in key=value format.
// Note: This will be used for all the provided urls.
func WithHeaderString(header string) (func(*Bench), error) {
	keyValue := strings.Split(header, "=")

	if len(keyValue) != 2 {
		return nil, fmt.Errorf("%s is not a correct key=value format", header)
	}

	return WithHeader(keyValue[0], keyValue[1]), nil
}

// Set endpoints from a file. The file is expected to have a string as defined
// in WithURLString in each line.
func WithFile(path string) (func(*Bench), error) {
	file, err := fs.Open(path)

	if err != nil {
		return nil, err
	}

	defer file.Close()

	s := bufio.NewScanner(file)
	s.Split(bufio.ScanLines)

	urls := make([]*URL, 0)

	for s.Scan() {
		parsedURL, err := parseURL(s.Text())

		if err != nil {
			return nil, err
		}

		urls = append(urls, parsedURL)
	}

	if len(urls) == 0 {
		return nil, fmt.Errorf("Did not find any url in the %s", path)
	}

	return func(b *Bench) {
		b.URLs = append(b.URLs, urls...)
	}, nil
}

// Sets a result report
func WithReport(report report.Report) func(*Bench) {
	return func(b *Bench) {
		b.Report = report
	}
}

func parseURL(u string) (*URL, error) {
	u = strings.Trim(u, "\" ")
	parts := strings.Split(u, "|")

	if len(parts) != 2 && len(parts) != 3 {
		return nil, fmt.Errorf(URL_ERROR_MESSAGE, u)
	}

	method := strings.ToUpper(parts[0])

	urlStruct, err := url.Parse(parts[1])
	if err != nil {
		return nil, fmt.Errorf("Invalid URL provided: %v", err)
	}

	if err := validateURL(urlStruct, method, len(parts)); err != nil {
		return nil, fmt.Errorf("Error for '%s': %v", u, err)
	}

	data := make(map[string]string)
	if method == http.MethodPost {
		if len(parts) != 3 {
			return nil, fmt.Errorf("%s\nWrong URL format. You need to provide post data. Example: POST|www.google.com|name=sasan&lastname=rose", u)
		}

		dataParts := strings.Split(parts[2], "&")
		for _, dataPart := range dataParts {
			keyValue := strings.Split(dataPart, "=")
			if len(keyValue) != 2 {
				return nil, fmt.Errorf("Wrong key value format for post data: %s", dataPart)
			}

			data[keyValue[0]] = keyValue[1]
		}
	}

	return &URL{
		Addr:   urlStruct.String(),
		Method: method,
		Data:   data,
	}, nil
}

func validateURL(urlStruct *url.URL, method string, lenParts int) error {
	if method != http.MethodGet && method != http.MethodPost && method != http.MethodHead {
		return errors.New("Method not allowed")
	}

	if urlStruct.Scheme != "http" && urlStruct.Scheme != "https" {
		return errors.New("Only http and https schemes are supported")
	}

	if (method == http.MethodGet || method == http.MethodHead) && lenParts > 2 {
		return errors.New("GET and HEAD do not need any data")
	}

	return nil
}
