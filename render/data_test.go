package render

import (
	"time"

	"github.com/sasanrose/gbench/report"
)

type testResponse struct {
	contentLength    int64
	timedOut, failed bool
	responseTime     time.Duration
	statusCode       int
}

var testData map[string][]*testResponse = map[string][]*testResponse{
	"http://testurl1.com": []*testResponse{
		&testResponse{contentLength: 15, responseTime: 600 * time.Microsecond, statusCode: 200},
		&testResponse{contentLength: 25, responseTime: 500 * time.Microsecond, statusCode: 200},
		&testResponse{contentLength: 25, responseTime: 550 * time.Microsecond, statusCode: 201},
		&testResponse{contentLength: 25, responseTime: 550 * time.Microsecond, statusCode: 500, failed: true},
		&testResponse{contentLength: 25, responseTime: 550 * time.Microsecond, timedOut: true},
	},
	"http://testurl2.com": []*testResponse{
		&testResponse{contentLength: 35, responseTime: 550 * time.Microsecond, statusCode: 500, failed: true},
		&testResponse{contentLength: 15, responseTime: 600 * time.Microsecond, statusCode: 200},
		&testResponse{contentLength: 25, responseTime: 500 * time.Microsecond, statusCode: 200},
		&testResponse{contentLength: 25, responseTime: 550 * time.Microsecond, timedOut: true},
		&testResponse{contentLength: 25, responseTime: 550 * time.Microsecond, statusCode: 201},
	},
	"http://testurl3.com": []*testResponse{
		&testResponse{contentLength: 10, responseTime: 550 * time.Microsecond, statusCode: 500, failed: true},
		&testResponse{contentLength: 10, responseTime: 550 * time.Microsecond, failed: true},
		&testResponse{contentLength: 25, responseTime: 500 * time.Microsecond, statusCode: 404, failed: true},
		&testResponse{contentLength: 20, responseTime: 550 * time.Microsecond, timedOut: true},
		&testResponse{contentLength: 20, responseTime: 550 * time.Microsecond, timedOut: true},
	},
}

func addTestData(r *report.Result) {
	for url, responses := range testData {
		for _, response := range responses {
			if response.timedOut {
				r.AddTimedoutResponse(url)
				continue
			}

			if response.statusCode == 0 && response.failed {
				r.AddFailedResponse(url)
				continue
			}

			r.AddReceivedDataLength(url, response.contentLength)
			r.AddResponseTime(url, response.responseTime)
			r.AddResponseStatusCode(url, response.statusCode, response.failed)
		}
	}

	r.SetTotalDuration(2 * time.Millisecond)
}
