package bench

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/sasanrose/gbench/result"
)

const URL_ERROR_MESSAGE = "%s\nWrong URL format. Example: GET|www.google.com?search=test or POST|www.google.com|search=test or HEAD|www.google.com"

// Create a config to set number of concurrent requests.
func WithConcurrency(n int) func(*Bench) {
	return func(b *Bench) {
		b.Concurrency = n
	}
}

// Create a config to set total number of requests to send.
func WithRequests(n int) func(*Bench) {
	return func(b *Bench) {
		b.Requests = n
	}
}

// Add an endpoint to benchmark.
func WithURL(u *Url) func(*Bench) {
	return func(b *Bench) {
		b.Urls = append(b.Urls, u)
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
	parsedUrl, err := parseUrl(u)

	if err != nil {
		return nil, err
	}

	f := WithURL(parsedUrl)

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

	urls := make([]*Url, 0)

	for s.Scan() {
		parsedUrl, err := parseUrl(s.Text())

		if err != nil {
			return nil, err
		}

		urls = append(urls, parsedUrl)
	}

	return func(b *Bench) {
		b.Urls = append(b.Urls, urls...)
	}, nil
}

// Sets a result renderer
func WithRenderer(renderer result.Renderer) func(*Bench) {
	return func(b *Bench) {
		b.Renderer = renderer
	}
}

func parseUrl(u string) (*Url, error) {
	u = strings.Trim(u, "\" ")
	parts := strings.Split(u, "|")

	if len(parts) != 2 && len(parts) != 3 {
		return nil, fmt.Errorf(URL_ERROR_MESSAGE, u)
	}

	method := strings.ToUpper(parts[0])
	if method != http.MethodGet && method != http.MethodPost && method != http.MethodHead {
		return nil, fmt.Errorf("Not an allowed method provided for: %s", u)
	}

	urlStruct, err := url.ParseRequestURI(parts[1])
	if err != nil {
		return nil, fmt.Errorf("Invalid URL provided: %v", err)
	}

	if urlStruct.Scheme != "http" && urlStruct.Scheme != "https" {
		return nil, fmt.Errorf("Only http and https schemes are supported for now: %s", u)
	}

	if (method == http.MethodGet || method == http.MethodHead) && len(parts) > 2 {
		return nil, fmt.Errorf("GET and HEAD do not need any data")
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

	return &Url{
		Addr:   urlStruct.String(),
		Method: method,
		Data:   data,
	}, nil
}
