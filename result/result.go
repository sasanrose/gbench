package result

import (
	"sync"
	"time"
)

type result struct {
	totalTime                                    time.Duration
	responseTime                                 map[string]time.Duration
	receivedDataLength                           map[string]int64
	responseStatusCode, failedResponseStatusCode map[string]map[int]int
	timedoutResponse                             map[string]int

	totalTimeLock, responseTimeLock, dataLengthLock, statusCodeLock, timedoutLock *sync.Mutex
}

func (r *result) Init() {
	r.responseTime = make(map[string]time.Duration)
	r.receivedDataLength = make(map[string]int64)
	r.responseStatusCode = make(map[string]map[int]int)
	r.failedResponseStatusCode = make(map[string]map[int]int)
	r.timedoutResponse = make(map[string]int)

	r.totalTimeLock = &sync.Mutex{}
	r.responseTimeLock = &sync.Mutex{}
	r.dataLengthLock = &sync.Mutex{}
	r.statusCodeLock = &sync.Mutex{}
	r.timedoutLock = &sync.Mutex{}
}

// Add content length received to total amount for a specific URL.
func (r *result) AddReceivedDataLength(url string, contentLength int64) {
	r.dataLengthLock.Lock()
	defer r.dataLengthLock.Unlock()

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
func (r *result) SetTotalDuration(duration time.Duration) {
	r.totalTimeLock.Lock()
	defer r.totalTimeLock.Unlock()

	r.totalTime = duration
}

// Add response time duration to total amount for a specific URL.
func (r *result) AddResponseTime(url string, time time.Duration) {
	r.responseTimeLock.Lock()
	defer r.responseTimeLock.Unlock()

	if _, ok := r.responseTime[url]; ok {
		r.responseTime[url] += time
		return
	}

	r.responseTime[url] = time
}

// Increament the number of responses with a specific status code for a
// specific url.
func (r *result) AddResponseStatusCode(url string, statusCode int, failed bool) {
	r.statusCodeLock.Lock()
	defer r.statusCodeLock.Unlock()

	if failed {
		updateStatusCode(&(*r).failedResponseStatusCode, url, statusCode)
		return
	}

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
func (r *result) AddTimedoutResponse(url string) {
	r.timedoutLock.Lock()
	defer r.timedoutLock.Unlock()

	if _, ok := r.timedoutResponse[url]; ok {
		r.timedoutResponse[url]++
		return
	}

	r.timedoutResponse[url] = 1
}
