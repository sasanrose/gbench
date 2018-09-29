package bench

import (
	"io"
	"net/http"
	"sync"
	"time"

	"github.com/sasanrose/gbench/report"
)

// A bnech represents a new benchmark that we want to execute.
type Bench struct {
	// Number of concurrent requests as well as total number of requests to
	// send.
	Concurrency, Requests int
	// Benchmarking endpoints.
	Urls []*Url
	// Optional basic HTTP authentication.
	Auth *Auth
	// Optional proxy address to use (Does not support authentication).
	Proxy string
	// Optional HTTP request headers.
	Headers map[string]string
	// Optional definition of status codes (Default is 200 and 201).
	SuccessStatusCodes []int
	// Verbosity writer.
	VerbosityWriter io.Writer
	// Verbosity writer lock.
	VerbosityWriterLock *sync.Mutex
	// Connection and response timeouts
	ResponseTimeout, ConnectionTimeout time.Duration
	// Http raw cookie string (i.e. the result of document.cookie).
	RawCookie string
	// Report to use
	Report report.Report
}

// A url represents an endpoint that we want to benchmark.
type Url struct {
	// Address and method to use for the endpoint.
	Addr, Method string
	// Optional data to send in the format of key-value.
	Data map[string]string
}

// An auth is used for a basic HTTP authentication.
type Auth struct {
	// Username and password to use with basic HTTP authentication.
	Username, Password string
}

// This function creates a new benchmark given a list of configurations. A
// config can be created on the fly or using the predefined functions.
func NewBench(configurations ...func(*Bench)) *Bench {
	b := &Bench{
		Headers:            make(map[string]string),
		Urls:               make([]*Url, 0),
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

	if b.VerbosityWriter != nil {
		b.VerbosityWriterLock = &sync.Mutex{}
	}

	if b.Concurrency == 0 {
		b.Concurrency = 1
	}

	if b.Requests == 0 {
		b.Requests = 1
	}

	return b
}
