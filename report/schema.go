package report

import (
	"sync"
	"time"
)

// Result struct implements Report interface and stores all the result
// information for a specific benchmark. This struct is used to encode the
// result to json and vice versa.
type Result struct {
	Urls map[string]bool `json:"urls"`

	TotalReceivedDataLength  int64                  `json:"total-received-data-length"`
	ReceivedDataLength       map[string]int64       `json:"received-data-length"`
	ResponseStatusCode       map[string]map[int]int `json:"response-status-code"`
	FailedResponseStatusCode map[string]map[int]int `json:"failed-response-status-code"`
	TimedoutResponse         map[string]int         `json:"timedout-response"`
	FailedResponse           map[string]int         `json:"failed-response"`
	TotalRequests            int                    `json:"total-requests"`
	SuccessfulRequests       int                    `json:"successful-requests"`
	FailedRequests           int                    `json:"failed-requests"`
	TimedOutRequests         int                    `json:"timedout-requests"`

	StartTime               time.Time                `json:"start-time"`
	EndTime                 time.Time                `json:"end-time"`
	TotalTime               time.Duration            `json:"total-time"`
	TotalResponseTime       time.Duration            `json:"total-response-time"`
	ResponseTimesTotalCount int                      `json:"response-times-total-count"`
	ResponseTime            map[string]time.Duration `json:"response-time"`
	ResponseTimesCount      map[string]int           `json:"response-times-count"`
	ShortestResponseTimes   map[string]time.Duration `json:"shortest-response-times"`
	LongestResponseTimes    map[string]time.Duration `json:"longest-response-times"`
	ShortestResponseTime    time.Duration            `json:"shortest-response-time"`
	LongestResponseTime     time.Duration            `json:"longest-response-time"`

	ConcurrencyResult  map[string][]*ConcurrencyResult `json:"concurrency-result"`
	concurrencyCounter map[string]int
	concurrency        int

	lock *sync.Mutex
}

// ConcurrencyResult struct store the result for each batch of concurrent
// requests.
type ConcurrencyResult struct {
	TotalRequests      int `json:"total-request"`
	SuccessfulRequests int `json:"successful-requests"`
	FailedRequests     int `json:"failed-requests"`
	TimedOutRequests   int `json:"timedout-requests"`
}
