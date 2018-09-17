package result

import (
	"io"
	"math"
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

	concurrencyResult  map[string][]*concurrencyResult
	concurrencyCounter map[string]int
	concurrency        int

	lock *sync.Mutex
}

type concurrencyResult struct {
	totalRequests, successfulRequests, failedRequests, timedOutRequests int
}

func (r *Result) Init(concurrency int) {
	r.responseTime = make(map[string]time.Duration)
	r.receivedDataLength = make(map[string]int64)
	r.responseStatusCode = make(map[string]map[int]int)
	r.failedResponseStatusCode = make(map[string]map[int]int)
	r.timedoutResponse = make(map[string]int)
	r.failedResponse = make(map[string]int)
	r.shortestResponseTimes = make(map[string]time.Duration)
	r.longestResponseTimes = make(map[string]time.Duration)
	r.concurrency = concurrency

	r.concurrencyResult = make(map[string][]*concurrencyResult)
	r.concurrencyCounter = make(map[string]int)

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
		r.updateConcurrencyResult(url, 0, 1, 0)
		r.failedRequests++
		updateStatusCode(&(*r).failedResponseStatusCode, url, statusCode)
		return
	}

	r.updateConcurrencyResult(url, 1, 0, 0)
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

	r.updateConcurrencyResult(url, 0, 0, 1)

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

	r.updateConcurrencyResult(url, 0, 1, 0)

	r.failedRequests++
	r.totalRequests++

	if _, ok := r.failedResponse[url]; ok {
		r.failedResponse[url]++
		return
	}

	r.failedResponse[url] = 1
}

func (r *Result) updateConcurrencyResult(url string, successfulRequests, failedRequests, timedOutRequests int) {
	if r.concurrency == 0 {
		return
	}

	defer func() {
		r.concurrencyCounter[url]++
	}()

	if _, ok := r.concurrencyResult[url]; !ok {
		r.concurrencyResult[url] = make([]*concurrencyResult, 0)
		r.concurrencyCounter[url] = 0
	}

	if r.concurrencyCounter[url] == 0 || int(math.Mod(float64(r.concurrencyCounter[url]), float64(r.concurrency))) == 0 {
		r.concurrencyResult[url] = append(r.concurrencyResult[url], &concurrencyResult{
			totalRequests:      1,
			successfulRequests: successfulRequests,
			failedRequests:     failedRequests,
			timedOutRequests:   timedOutRequests,
		})
		return
	}

	lenResult := len(r.concurrencyResult[url])

	r.concurrencyResult[url][lenResult-1].totalRequests++
	r.concurrencyResult[url][lenResult-1].failedRequests += failedRequests
	r.concurrencyResult[url][lenResult-1].successfulRequests += successfulRequests
	r.concurrencyResult[url][lenResult-1].timedOutRequests += timedOutRequests
}
