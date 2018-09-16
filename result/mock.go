package result

import (
	"time"
)

type mockedRenderer struct {
	Result
	BenchResult *MockedResult
}

type MockedResult struct {
	ReceivedDataLength                                                  map[string]int64
	ResponseStatusCode, FailedResponseStatusCode                        map[string]map[int]int
	TimedoutResponse, FailedResponse                                    map[string]int
	TotalRequests, SuccessfulRequests, FailedRequests, TimedOutRequests int

	TotalTime                                   time.Duration
	ResponseTime                                map[string]time.Duration
	ShortestResponseTimes, LongestResponseTimes map[string]time.Duration
	ShortestResponseTime, LongestResponseTime   time.Duration
}

func NewMockRenderer() *mockedRenderer {
	return &mockedRenderer{}
}

func (m *mockedRenderer) Render() error {
	m.BenchResult = &MockedResult{
		TotalTime:                m.totalTime,
		ResponseTime:             m.responseTime,
		ReceivedDataLength:       m.receivedDataLength,
		ResponseStatusCode:       m.responseStatusCode,
		FailedResponseStatusCode: m.failedResponseStatusCode,
		TimedoutResponse:         m.timedoutResponse,
		FailedResponse:           m.failedResponse,
		TotalRequests:            m.totalRequests,
		SuccessfulRequests:       m.successfulRequests,
		FailedRequests:           m.failedRequests,
		TimedOutRequests:         m.timedOutRequests,
		ShortestResponseTime:     m.shortestResponseTime,
		LongestResponseTime:      m.longestResponseTime,
		ShortestResponseTimes:    m.shortestResponseTimes,
		LongestResponseTimes:     m.longestResponseTimes,
	}

	return nil
}
