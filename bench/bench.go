// Package bench executes a benchmark based on the given configurations.
package bench

import (
	"io"
	"net/http"
	"sync"
	"time"

	"github.com/sasanrose/gbench/report"
)

// Bench represents a new benchmark that we want to execute.
type Bench struct {
	// Number of concurrent requests as well as total number of requests to
	// send.
	Concurrency, Requests int
	// Benchmarking endpoints.
	URLs []*URL
	// Optional basic HTTP authentication.
	Auth *Auth
	// Optional proxy address to use (Does not support authentication).
	Proxy string
	// Optional HTTP request headers.
	Headers map[string]string
	// Optional definition of status codes (Default is 200 and 201).
	SuccessStatusCodes []int
	// Output writer.
	OutputWriter io.Writer
	// Output writer writer lock.
	OutputWriterLock *sync.Mutex
	// Connection and response timeouts
	ResponseTimeout, ConnectionTimeout time.Duration
	// Optional HTTP raw cookie string (i.e. the result of document.cookie).
	RawCookie string
	// Report to use
	Report report.Report
}

// URL represents an endpoint that we want to benchmark.
type URL struct {
	// Address and method to use for the endpoint.
	Addr, Method string
	// Optional data to send in the format of key-value.
	Data map[string]string
	// Optional URL specific HTTP request headers.
	Headers map[string]string
	// Optional URL specific HTTP raw cookie string (i.e. the result of document.cookie).
	RawCookie string
	// Optional URL specific basic HTTP authentication.
	Auth *Auth
}

// Auth is used for a basic HTTP authentication.
type Auth struct {
	// Username and password to use with basic HTTP authentication.
	Username, Password string
}

// NewBench creates a new benchmark given a list of configurations. A
// config can be created on the fly or using the predefined functions.
func NewBench(configurations ...func(*Bench)) *Bench {
	b := &Bench{
		Headers:            make(map[string]string),
		URLs:               make([]*URL, 0),
		SuccessStatusCodes: make([]int, 0),
	}

	for _, config := range configurations {
		config(b)
	}

	if len(b.SuccessStatusCodes) == 0 {
		b.SuccessStatusCodes = []int{
			http.StatusOK,
			http.StatusAccepted,
			http.StatusCreated,
		}
	}

	if b.OutputWriter != nil {
		b.OutputWriterLock = &sync.Mutex{}
	}

	if b.Concurrency == 0 {
		b.Concurrency = 1
	}

	if b.Requests == 0 {
		b.Requests = 1
	}

	return b
}
