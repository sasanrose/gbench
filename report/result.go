package report

import (
	"math"
	"sync"
	"time"
)

// Init initializes a result to be used to store the benchmark result. Method
// only accepts an integer to set as the number of concurrent requests which
// are supposed to be sent.
func (r *Result) Init(concurrency int) {
	r.ResponseTime = make(map[string]time.Duration)
	r.ReceivedDataLength = make(map[string]int64)
	r.ResponseStatusCode = make(map[string]map[int]int)
	r.FailedResponseStatusCode = make(map[string]map[int]int)
	r.TimedoutResponse = make(map[string]int)
	r.FailedResponse = make(map[string]int)
	r.ShortestResponseTimes = make(map[string]time.Duration)
	r.LongestResponseTimes = make(map[string]time.Duration)
	r.concurrency = concurrency
	r.Urls = make(map[string]bool)
	r.ResponseTimesCount = make(map[string]int)

	r.ConcurrencyResult = make(map[string][]*ConcurrencyResult)
	r.concurrencyCounter = make(map[string]int)

	r.lock = &sync.Mutex{}
}

// SetStartTime sets benchmark's start time.
func (r *Result) SetStartTime(t time.Time) {
	r.lock.Lock()
	defer r.lock.Unlock()

	r.StartTime = t
}

// SetEndTime sets benchmark's end time.
func (r *Result) SetEndTime(t time.Time) {
	r.lock.Lock()
	defer r.lock.Unlock()

	r.EndTime = t
}

// AddReceivedDataLength adds content length received to the Total
// amount for a specific URL.
func (r *Result) AddReceivedDataLength(url string, contentLength int64) {
	r.lock.Lock()
	defer r.lock.Unlock()

	r.Urls[url] = true

	// Nothing to do
	if contentLength <= 0 {
		return
	}

	r.TotalReceivedDataLength += contentLength

	if _, ok := r.ReceivedDataLength[url]; ok {
		r.ReceivedDataLength[url] += contentLength
		return
	}

	r.ReceivedDataLength[url] = contentLength
}

// SetTotalDuration sets the total time elapsed since the beginning of the bench marck.
func (r *Result) SetTotalDuration(duration time.Duration) {
	r.lock.Lock()
	defer r.lock.Unlock()

	r.TotalTime = duration
}

// AddResponseTime add response time duration to total amount for a specific URL.
func (r *Result) AddResponseTime(url string, responseTime time.Duration) {
	r.lock.Lock()
	defer r.lock.Unlock()

	r.Urls[url] = true
	r.ResponseTimesTotalCount++
	r.TotalResponseTime += responseTime

	if r.ShortestResponseTime == 0*time.Second || r.ShortestResponseTime > responseTime {
		r.ShortestResponseTime = responseTime
	}

	if r.LongestResponseTime == 0*time.Second || r.LongestResponseTime < responseTime {
		r.LongestResponseTime = responseTime
	}

	r.updateURLShortestTime(url, responseTime)
	r.updateURLLongestTime(url, responseTime)

	if _, ok := r.ResponseTime[url]; ok {
		r.ResponseTime[url] += responseTime
		r.ResponseTimesCount[url]++
		return
	}

	r.ResponseTimesCount[url] = 1
	r.ResponseTime[url] = responseTime
}

func (r *Result) updateURLShortestTime(url string, time time.Duration) {
	if _, ok := r.ShortestResponseTimes[url]; !ok {
		r.ShortestResponseTimes[url] = time
		return
	}

	if r.ShortestResponseTimes[url] > time {
		r.ShortestResponseTimes[url] = time
	}
}

func (r *Result) updateURLLongestTime(url string, time time.Duration) {
	if _, ok := r.LongestResponseTimes[url]; !ok {
		r.LongestResponseTimes[url] = time
		return
	}

	if r.LongestResponseTimes[url] < time {
		r.LongestResponseTimes[url] = time
	}
}

// AddResponseStatusCode increaments the number of responses with a specific
// status code for a specific url.
func (r *Result) AddResponseStatusCode(url string, statusCode int, failed bool) {
	r.lock.Lock()
	defer r.lock.Unlock()

	r.Urls[url] = true

	r.TotalRequests++

	if failed {
		r.updateConcurrencyResult(url, 0, 1, 0)
		r.FailedRequests++
		updateStatusCode(&(*r).FailedResponseStatusCode, url, statusCode)
		return
	}

	r.updateConcurrencyResult(url, 1, 0, 0)
	r.SuccessfulRequests++
	updateStatusCode(&(*r).ResponseStatusCode, url, statusCode)
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

// AddTimedoutResponse increaments the number of timed out responses for a url.
func (r *Result) AddTimedoutResponse(url string) {
	r.lock.Lock()
	defer r.lock.Unlock()

	r.Urls[url] = true

	r.updateConcurrencyResult(url, 0, 0, 1)

	r.TimedOutRequests++
	r.TotalRequests++

	if _, ok := r.TimedoutResponse[url]; ok {
		r.TimedoutResponse[url]++
		return
	}

	r.TimedoutResponse[url] = 1
}

// AddFailedResponse increaments the number of failed responses for a url.
func (r *Result) AddFailedResponse(url string) {
	r.lock.Lock()
	defer r.lock.Unlock()

	r.Urls[url] = true

	r.updateConcurrencyResult(url, 0, 1, 0)

	r.FailedRequests++
	r.TotalRequests++

	if _, ok := r.FailedResponse[url]; ok {
		r.FailedResponse[url]++
		return
	}

	r.FailedResponse[url] = 1
}

func (r *Result) updateConcurrencyResult(url string, successfulRequests, failedRequests, timedOutRequests int) {
	if r.concurrency == 0 {
		return
	}

	defer func() {
		r.concurrencyCounter[url]++
	}()

	if _, ok := r.ConcurrencyResult[url]; !ok {
		r.ConcurrencyResult[url] = make([]*ConcurrencyResult, 0)
		r.concurrencyCounter[url] = 0
	}

	if r.concurrencyCounter[url] == 0 || int(math.Mod(float64(r.concurrencyCounter[url]), float64(r.concurrency))) == 0 {
		r.ConcurrencyResult[url] = append(r.ConcurrencyResult[url], &ConcurrencyResult{
			TotalRequests:      1,
			SuccessfulRequests: successfulRequests,
			FailedRequests:     failedRequests,
			TimedOutRequests:   timedOutRequests,
		})
		return
	}

	lenResult := len(r.ConcurrencyResult[url])

	r.ConcurrencyResult[url][lenResult-1].TotalRequests++
	r.ConcurrencyResult[url][lenResult-1].FailedRequests += failedRequests
	r.ConcurrencyResult[url][lenResult-1].SuccessfulRequests += successfulRequests
	r.ConcurrencyResult[url][lenResult-1].TimedOutRequests += timedOutRequests
}
