package bench

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	"github.com/sasanrose/gbench/result"
)

type testHttp struct {
	lock                      *sync.Mutex
	requests                  []*testRequest
	totalRequests, statusCode int
}

type testRequest struct {
	path, method string
	data         map[string]string
	headers      map[string]string
	cookie       string
}

type expectedResult struct {
	receivedDataLength                                                  map[string]int64
	responseStatusCode, failedResponseStatusCode                        map[string]map[int]int
	timedoutResponse, failedResponse                                    map[string]int
	totalRequests, successfulRequests, failedRequests, timedOutRequests int
	concurrencyResult                                                   map[string][]*expectedConcurrencyResult
}

type expectedConcurrencyResult struct {
	totalRequests, successfulRequests, failedRequests, timedOutRequests int
}

func newTestHttp(statusCode int) *testHttp {
	return &testHttp{
		lock:       &sync.Mutex{},
		requests:   make([]*testRequest, 0),
		statusCode: statusCode,
	}
}

func (h *testHttp) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.lock.Lock()
	defer h.lock.Unlock()

	h.totalRequests++

	testReq := &testRequest{
		path:    r.RequestURI,
		method:  r.Method,
		data:    make(map[string]string),
		headers: make(map[string]string),
		cookie:  r.Header.Get("Set-Cookie"),
	}

	for k, v := range r.Header {
		testReq.headers[k] = v[0]
	}

	r.ParseForm()
	for k, v := range r.PostForm {
		testReq.data[k] = v[0]
	}

	h.requests = append(h.requests, testReq)

	w.WriteHeader(h.statusCode)
	w.Write([]byte("Test data"))
}

func TestExec(t *testing.T) {
	hOk := newTestHttp(http.StatusOK)
	ts1 := httptest.NewServer(hOk)
	defer ts1.Close()

	hCreated := newTestHttp(http.StatusCreated)
	ts2 := httptest.NewServer(hCreated)
	defer ts2.Close()

	hNotFound := newTestHttp(http.StatusNotFound)
	ts3 := httptest.NewServer(hNotFound)
	defer ts3.Close()

	r := result.NewMockRenderer()
	r.Init(2)

	url1 := ts1.URL + "/one"
	url2 := ts2.URL + "/two"
	url3 := ts3.URL + "/three"

	withUrl1, _ := WithURLString("GET|" + url1)
	withUrl2, _ := WithURLString("POST|" + url2 + "|foo=bar&foo2=bar2")
	withUrl3, _ := WithURLString("HEAD|" + url3)

	var buf bytes.Buffer

	configurations := []func(*Bench){
		WithConcurrency(2),
		WithRequests(4),
		withUrl1,
		withUrl2,
		withUrl3,
		WithVerbosity(&buf),
		WithRawCookie("testCookie"),
		WithHeader("Test-Key", "testVal"),
		WithRenderer(r),
	}

	b := NewBench(configurations...)
	b.Exec(context.Background())

	r.Render()

	if buf.Len() == 0 {
		t.Errorf("Verbosity writer is empty")
	}

	expected := &expectedResult{
		receivedDataLength: map[string]int64{
			url1: 36,
			url2: 36,
			url3: 36,
		},
		responseStatusCode: map[string]map[int]int{
			url1: map[int]int{http.StatusOK: 4},
			url2: map[int]int{http.StatusCreated: 4},
		},
		failedResponseStatusCode: map[string]map[int]int{
			url3: map[int]int{http.StatusNotFound: 4},
		},
		timedoutResponse:   map[string]int{},
		failedResponse:     map[string]int{},
		totalRequests:      12,
		successfulRequests: 8,
		failedRequests:     4,
		timedOutRequests:   0,
		concurrencyResult: map[string][]*expectedConcurrencyResult{
			url1: []*expectedConcurrencyResult{
				&expectedConcurrencyResult{2, 2, 0, 0},
				&expectedConcurrencyResult{2, 2, 0, 0},
			},
			url2: []*expectedConcurrencyResult{
				&expectedConcurrencyResult{2, 2, 0, 0},
				&expectedConcurrencyResult{2, 2, 0, 0},
			},
			url3: []*expectedConcurrencyResult{
				&expectedConcurrencyResult{2, 0, 2, 0},
				&expectedConcurrencyResult{2, 0, 2, 0},
			},
		},
	}

	checkResult(t, r.BenchResult, expected)

	if hOk.totalRequests != 4 || hCreated.totalRequests != 4 || hNotFound.totalRequests != 4 {
		t.Errorf("Wrong number of requests are sent to the servers")
	}

	expectedHeaders := map[string]string{
		"User-Agent": "Gbench",
		"Test-Key":   "testVal",
	}

	expectedRequest := &testRequest{
		path:    "/one",
		method:  http.MethodGet,
		cookie:  "testCookie",
		headers: expectedHeaders,
		data:    make(map[string]string),
	}

	checkRequest(t, hOk, expectedRequest)

	expectedRequest = &testRequest{
		path:    "/two",
		method:  http.MethodPost,
		cookie:  "testCookie",
		headers: expectedHeaders,
		data:    map[string]string{"foo": "bar", "foo2": "bar2"},
	}

	checkRequest(t, hCreated, expectedRequest)

	expectedRequest = &testRequest{
		path:    "/three",
		method:  http.MethodHead,
		cookie:  "testCookie",
		headers: expectedHeaders,
		data:    make(map[string]string),
	}

	checkRequest(t, hNotFound, expectedRequest)
}

func checkRequest(t *testing.T, h *testHttp, expected *testRequest) {
	for _, request := range h.requests {
		if request.cookie != expected.cookie {
			t.Errorf("Expected %s as cookie but got %s", expected.cookie, request.cookie)
		}

		if request.path != expected.path {
			t.Errorf("Expected %s as path but got %s", expected.path, request.path)
		}

		if request.method != expected.method {
			t.Errorf("Expected %s as method but got %s", expected.method, request.method)
		}

		if len(request.data) != len(expected.data) {
			t.Error("Wrong number of data")
		}

		for key, expectedValue := range expected.data {
			if value, ok := request.data[key]; !ok || value != expectedValue {
				t.Errorf("Expected %s at %s in request data for path %s", expectedValue, key, request.path)
			}
		}

		for key, expectedValue := range expected.headers {
			if value, ok := request.headers[key]; !ok || value != expectedValue {
				t.Errorf("Expected %s at %s in request headers for path %s", expectedValue, key, request.path)
			}
		}
	}
}

func checkResult(t *testing.T, r *result.MockedResult, expected *expectedResult) {
	if r.TotalRequests != expected.totalRequests ||
		r.SuccessfulRequests != expected.successfulRequests ||
		r.FailedRequests != expected.failedRequests ||
		r.TimedOutRequests != expected.timedOutRequests {
		t.Error("Unexpected count of requests")
	}

	if len(r.ReceivedDataLength) != len(expected.receivedDataLength) {
		t.Error("Unexpected receivedDataLength")
	}

	for url, expectedVal := range expected.receivedDataLength {
		if value, ok := r.ReceivedDataLength[url]; !ok || expectedVal != value {
			t.Errorf("Expected value %v for %s in receivedDataLength", expectedVal, url)
		}
	}

	checkStatusCodes(t, r.ResponseStatusCode, expected.responseStatusCode, "responseStatusCode")
	checkStatusCodes(t, r.FailedResponseStatusCode, expected.failedResponseStatusCode, "failedResponseStatusCode")

	for url, expectedConcurrencyResults := range expected.concurrencyResult {
		if _, ok := r.ConcurrencyResult[url]; !ok {
			t.Errorf("Expected to get a concurrencyResult for %s but got nothing", url)
			continue
		}

		for key, expectedConcurrencyResult := range expectedConcurrencyResults {
			if expectedConcurrencyResult.failedRequests != r.ConcurrencyResult[url][key].FailedRequests ||
				expectedConcurrencyResult.totalRequests != r.ConcurrencyResult[url][key].TotalRequests ||
				expectedConcurrencyResult.successfulRequests != r.ConcurrencyResult[url][key].SuccessfulRequests ||
				expectedConcurrencyResult.timedOutRequests != r.ConcurrencyResult[url][key].TimedOutRequests {
				t.Errorf("Expected to get %+v for %s in concurrency Result index %d but got %+v",
					expectedConcurrencyResult,
					url,
					key,
					r.ConcurrencyResult[url][key])

			}
		}
	}
}

func checkStatusCodes(t *testing.T, statusCodes, expectedStatusCodes map[string]map[int]int, fieldName string) {
	if len(expectedStatusCodes) != len(statusCodes) {
		t.Errorf("Unexpected %s", fieldName)
	}

	for url, expectedMethods := range expectedStatusCodes {
		if _, ok := statusCodes[url]; !ok {
			t.Errorf("Expected url %s in %s but got nothing", url, fieldName)
			continue
		}

		for method, expectedVal := range expectedMethods {
			if value, ok := statusCodes[url][method]; !ok || expectedVal != value {
				t.Errorf("Expected value %v for %s in %s", expectedVal, url, fieldName)
			}
		}
	}
}
