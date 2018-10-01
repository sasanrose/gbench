// Package report helps to create a report file from a bench mark result.
package report

import "time"

// Report defines the interface for a type report that can be used with
// benchmarks to store the result.
type Report interface {
	AddReceivedDataLength(url string, contentLength int64)
	SetTotalDuration(duration time.Duration)
	AddResponseTime(url string, time time.Duration)
	AddResponseStatusCode(url string, statusCode int, failed bool)
	AddTimedoutResponse(url string)
	AddFailedResponse(url string)
	Init(concurrency int)
	SetStartTime(t time.Time)
	SetEndTime(t time.Time)
}
