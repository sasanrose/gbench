package result

import "time"

type Renderer interface {
	Render() error
	AddReceivedDataLength(url string, contentLength int64)
	SetTotalDuration(duration time.Duration)
	AddResponseTime(url string, time time.Duration)
	AddResponseStatusCode(url string, statusCode int, failed bool)
	AddTimedoutResponse(url string)
	AddFailedResponse(url string)
}
