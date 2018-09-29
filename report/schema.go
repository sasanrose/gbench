package report

import (
	"sync"
	"time"
)

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

type ConcurrencyResult struct {
	TotalRequests      int `json:"total-request"`
	SuccessfulRequests int `json:"sucessful-requests"`
	FailedRequests     int `json:"failed-requests"`
	TimedOutRequests   int `json:"timedout-requests"`
}