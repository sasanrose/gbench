package result

import (
	"io"
	"sync"
	"time"
)

type Result struct {
	Output io.Writer

	receivedDataLength                                                  map[string]int64
	responseStatusCode, failedResponseStatusCode                        map[string]map[int]int
	timedoutResponse, failedResponse                                    map[string]int
	totalRequests, successfulRequests, failedRequests, timedOutRequests int

	totalTime                                   time.Duration
	responseTime                                map[string]time.Duration
	shortestResponseTimes, longestResponseTimes map[string]time.Duration
	shortestResponseTime, longestResponseTime   time.Duration

	lock *sync.Mutex
}

func (r *Result) Init() {
	r.responseTime = make(map[string]time.Duration)
	r.receivedDataLength = make(map[string]int64)
	r.responseStatusCode = make(map[string]map[int]int)
	r.failedResponseStatusCode = make(map[string]map[int]int)
	r.timedoutResponse = make(map[string]int)
	r.failedResponse = make(map[string]int)
	r.shortestResponseTimes = make(map[string]time.Duration)
	r.longestResponseTimes = make(map[string]time.Duration)

	r.lock = &sync.Mutex{}
}

// Add content length received to total amount for a specific URL.
func (r *Result) AddReceivedDataLength(url string, contentLength int64) {
	r.lock.Lock()
	defer r.lock.Unlock()

	// Nothing to do
	if contentLength <= 0 {
		return
	}

	if _, ok := r.receivedDataLength[url]; ok {
		r.receivedDataLength[url] += contentLength
		return
	}

	r.receivedDataLength[url] = contentLength
}

// Set the total time elapsed since the begining of the bench marck.
func (r *Result) SetTotalDuration(duration time.Duration) {
	r.lock.Lock()
	defer r.lock.Unlock()

	r.totalTime = duration
}

// Add response time duration to total amount for a specific URL.
func (r *Result) AddResponseTime(url string, responseTime time.Duration) {
	r.lock.Lock()
	defer r.lock.Unlock()

	if r.shortestResponseTime == 0*time.Second || r.shortestResponseTime > responseTime {
		r.shortestResponseTime = responseTime
	}

	if r.longestResponseTime == 0*time.Second || r.longestResponseTime < responseTime {
		r.longestResponseTime = responseTime
	}

	r.updateUrlShortestTime(url, responseTime)
	r.updateUrlLongestTime(url, responseTime)

	if _, ok := r.responseTime[url]; ok {
		r.responseTime[url] += responseTime
		return
	}

	r.responseTime[url] = responseTime
}

func (r *Result) updateUrlShortestTime(url string, time time.Duration) {
	if _, ok := r.shortestResponseTimes[url]; !ok {
		r.shortestResponseTimes[url] = time
		return
	}

	if r.shortestResponseTimes[url] > time {
		r.shortestResponseTimes[url] = time
	}
}

func (r *Result) updateUrlLongestTime(url string, time time.Duration) {
	if _, ok := r.longestResponseTimes[url]; !ok {
		r.longestResponseTimes[url] = time
		return
	}

	if r.longestResponseTimes[url] < time {
		r.longestResponseTimes[url] = time
	}
}

// Increament the number of responses with a specific status code for a
// specific url.
func (r *Result) AddResponseStatusCode(url string, statusCode int, failed bool) {
	r.lock.Lock()
	defer r.lock.Unlock()

	r.totalRequests++

	if failed {
		r.failedRequests++
		updateStatusCode(&(*r).failedResponseStatusCode, url, statusCode)
		return
	}

	r.successfulRequests++
	updateStatusCode(&(*r).responseStatusCode, url, statusCode)
}

func updateStatusCode(statusCodeMap *map[string]map[int]int, url string, statusCode int) {
	if _, ok := (*statusCodeMap)[url]; !ok {
		(*statusCodeMap)[url] = make(map[int]int)
	}

	if _, ok := (*statusCodeMap)[url][statusCode]; ok {
		(*statusCodeMap)[url][statusCode]++
		return
	}

	(*statusCodeMap)[url][statusCode] = 1
}

// Increament the number of timed out responses for a url
func (r *Result) AddTimedoutResponse(url string) {
	r.lock.Lock()
	defer r.lock.Unlock()

	r.timedOutRequests++
	r.totalRequests++

	if _, ok := r.timedoutResponse[url]; ok {
		r.timedoutResponse[url]++
		return
	}

	r.timedoutResponse[url] = 1
}

// Increament the number of failed responses for a url.
func (r *Result) AddFailedResponse(url string) {
	r.lock.Lock()
	defer r.lock.Unlock()

	r.failedRequests++
	r.totalRequests++

	if _, ok := r.failedResponse[url]; ok {
		r.failedResponse[url]++
		return
	}

	r.failedResponse[url] = 1
}
