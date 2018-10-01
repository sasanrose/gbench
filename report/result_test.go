package report

import (
	"testing"
	"time"
)

func getTestResultStruct() *Result {
	r := &Result{}
	r.Init(2)

	return r
}

func TestTimes(t *testing.T) {
	r := getTestResultStruct()

	s := time.Now()

	r.SetStartTime(s)
	r.SetEndTime(s.Add(10 * time.Second))

	if r.StartTime != s {
		t.Error("Enexpected start time")
	}

	if r.EndTime.Sub(s) != 10*time.Second {
		t.Error("Enexpected end time")
	}
}

func TestContentLength(t *testing.T) {
	r := getTestResultStruct()

	r.AddReceivedDataLength("testURL1", 10)
	if r.ReceivedDataLength["testURL1"] != 10 {
		t.Fatal("Expected to get '10' for 'testURL1'")
	}

	r.AddReceivedDataLength("testURL1", -1)
	if r.ReceivedDataLength["testURL1"] != 10 {
		t.Fatal("Expected to get '10' for 'testURL1'")
	}

	r.AddReceivedDataLength("testURL1", 2)
	if r.ReceivedDataLength["testURL1"] != 12 {
		t.Fatal("Expected to get '12' for 'testURL1'")
	}

	r.AddReceivedDataLength("testURL2", 20)
	if r.ReceivedDataLength["testURL2"] != 20 {
		t.Fatal("Expected to get '20' for 'testURL2'")
	}

	if r.ReceivedDataLength["testURL1"] != 12 {
		t.Fatal("Expected to get '12' for 'testURL1'")
	}

	if r.TotalReceivedDataLength != 32 {
		t.Fatal("Expected to get 32 as TotalReceivedDataLength")
	}
}

func TestResponseTime(t *testing.T) {
	r := getTestResultStruct()

	r.AddResponseTime("testURL1", 10*time.Second)
	if r.ResponseTime["testURL1"] != 10*time.Second {
		t.Fatal("Expected to get '10' seconds for 'testURL1'")
	}

	r.AddResponseTime("testURL1", 2*time.Second)
	if r.ResponseTime["testURL1"] != 12*time.Second {
		t.Fatal("Expected to get '12' seconds for 'testURL1'")
	}

	r.AddResponseTime("testURL2", 20*time.Second)
	if r.ResponseTime["testURL2"] != 20*time.Second {
		t.Fatal("Expected to get '20' seconds for 'testURL2'")
	}

	if r.TotalResponseTime != 32*time.Second {
		t.Fatal("Expected to get 32 as TotalResponseTime")
	}

	if r.ResponseTime["testURL1"] != 12*time.Second {
		t.Fatal("Expected to get '12' seconds for 'testURL1'")
	}

	if r.ShortestResponseTime != 2*time.Second {
		t.Fatalf("Expected to get '2' as the shortest response time but got %v", r.ShortestResponseTime)
	}

	if r.LongestResponseTime != 20*time.Second {
		t.Fatalf("Expected to get '20' as the longest response time but got %v", r.LongestResponseTime)
	}

	if r.ShortestResponseTimes["testURL1"] != 2*time.Second {
		t.Fatalf("Expected to get '2' as the shortest response time for 'testURL1' but got %v", r.ShortestResponseTimes["testURL1"])
	}

	if r.ShortestResponseTimes["testURL2"] != 20*time.Second {
		t.Fatalf("Expected to get '20' as the shortest response time for 'testURL2' but got %v", r.ShortestResponseTimes["testURL2"])
	}

	if r.LongestResponseTimes["testURL1"] != 10*time.Second {
		t.Fatalf("Expected to get '2' as the longest response time for 'testURL1' but got %v", r.LongestResponseTimes["testURL1"])
	}

	if r.LongestResponseTimes["testURL2"] != 20*time.Second {
		t.Fatalf("Expected to get '20' as the longest response time for 'testURL2' but got %v", r.LongestResponseTimes["testURL2"])
	}

	if r.ResponseTimesTotalCount != 3 {
		t.Errorf("Expected to get 3 for ResponseTimesTotalCount but got %d", r.ResponseTimesTotalCount)
	}

	if r.ResponseTimesCount["testURL1"] != 2 {
		t.Errorf("Expected to get 2 for ResponseTimesCount for testURL1 but got %d", r.ResponseTimesCount["testURL1"])
	}

	if r.ResponseTimesCount["testURL2"] != 1 {
		t.Errorf("Expected to get 1 for ResponseTimesCount for testURL2 but got %d", r.ResponseTimesCount["testURL1"])
	}
}

func TestTotalTime(t *testing.T) {
	r := getTestResultStruct()
	r.SetTotalDuration(100 * time.Second)
	if r.TotalTime != 100*time.Second {
		t.Fatal("Expected to get '100' seconds for total time")
	}
}

func TestTimedoutResponse(t *testing.T) {
	r := getTestResultStruct()

	r.AddTimedoutResponse("testURL1")
	if r.TimedoutResponse["testURL1"] != 1 {
		t.Fatal("Expected to get '1' for 'testURL1'")
	}

	r.AddTimedoutResponse("testURL1")
	if r.TimedoutResponse["testURL1"] != 2 {
		t.Fatal("Expected to get '2' for 'testURL1'")
	}

	r.AddTimedoutResponse("testURL2")
	if r.TimedoutResponse["testURL2"] != 1 {
		t.Fatal("Expected to get '1' for 'testURL2'")
	}

	if r.TimedoutResponse["testURL1"] != 2 {
		t.Fatal("Expected to get '2' for 'testURL1'")
	}

	if r.TotalRequests != 3 || r.TimedOutRequests != 3 {
		t.Error("Number of total requets and timed out requests are not set as expected")
	}

	if len(r.ConcurrencyResult["testURL1"]) != 1 || len(r.ConcurrencyResult["testURL2"]) != 1 {
		t.Error("Wrong ConcurrencyResult")
	}

	checkConcurrencyResult(t, r.ConcurrencyResult["testURL1"][0], 2, 0, 0, 2)
	checkConcurrencyResult(t, r.ConcurrencyResult["testURL2"][0], 1, 0, 0, 1)

	checkURLs(t, r, []string{"testURL1", "testURL2"})
}

func TestFailedResponse(t *testing.T) {
	r := getTestResultStruct()

	r.AddFailedResponse("testURL1")
	if r.FailedResponse["testURL1"] != 1 {
		t.Fatal("Expected to get '1' for 'testURL1'")
	}

	r.AddFailedResponse("testURL1")
	if r.FailedResponse["testURL1"] != 2 {
		t.Fatal("Expected to get '2' for 'testURL1'")
	}

	r.AddFailedResponse("testURL2")
	if r.FailedResponse["testURL2"] != 1 {
		t.Fatal("Expected to get '1' for 'testURL2'")
	}

	if r.FailedResponse["testURL1"] != 2 {
		t.Fatal("Expected to get '2' for 'testURL1'")
	}

	if r.TotalRequests != 3 || r.FailedRequests != 3 {
		t.Error("Number of total requets and failed requests are not set as expected")
	}

	if len(r.ConcurrencyResult["testURL1"]) != 1 || len(r.ConcurrencyResult["testURL2"]) != 1 {
		t.Error("Wrong ConcurrencyResult")
	}

	checkConcurrencyResult(t, r.ConcurrencyResult["testURL1"][0], 2, 0, 2, 0)
	checkConcurrencyResult(t, r.ConcurrencyResult["testURL2"][0], 1, 0, 1, 0)

	checkURLs(t, r, []string{"testURL1", "testURL2"})
}

func TestStatusCode(t *testing.T) {
	r := getTestResultStruct()

	r.AddResponseStatusCode("testURL1", 200, false)
	if r.ResponseStatusCode["testURL1"][200] != 1 {
		t.Fatal("Expected count '1' for status cude '200' for testURL1")
	}

	r.AddResponseStatusCode("testURL1", 200, false)
	if r.ResponseStatusCode["testURL1"][200] != 2 {
		t.Fatal("Expected count '2' for status cude '200' for testURL1")
	}

	r.AddResponseStatusCode("testURL2", 200, false)
	if r.ResponseStatusCode["testURL2"][200] != 1 {
		t.Fatal("Expected count '1' for status cude '200' for testURL2")
	}

	r.AddResponseStatusCode("testURL1", 201, false)
	if r.ResponseStatusCode["testURL1"][201] != 1 {
		t.Fatal("Expected count '1' for status cude '201' for testURL1")
	}

	if len(r.FailedResponseStatusCode) != 0 {
		t.Fatal("Expected empty for FailedResponseStatusCode")
	}

	r.AddResponseStatusCode("testURL2", 500, true)
	if r.FailedResponseStatusCode["testURL2"][500] != 1 {
		t.Fatal("Expected count '1' for status cude '500' for testURL2")
	}

	r.AddResponseStatusCode("testURL2", 500, true)
	if r.FailedResponseStatusCode["testURL2"][500] != 2 {
		t.Fatal("Expected count '2' for status cude '500' for testURL2")
	}

	r.AddResponseStatusCode("testURL3", 500, true)
	if r.FailedResponseStatusCode["testURL3"][500] != 1 {
		t.Fatal("Expected count '1' for status cude '500' for testURL3")
	}

	if r.TotalRequests != 7 || r.SuccessfulRequests != 4 || r.FailedRequests != 3 {
		t.Error("Number of total requests, successful requests and failed requests are not set as expected")
	}

	if len(r.ConcurrencyResult["testURL1"]) != 2 ||
		len(r.ConcurrencyResult["testURL2"]) != 2 ||
		len(r.ConcurrencyResult["testURL3"]) != 1 {
		t.Error("Wrong ConcurrencyResult")
	}

	checkConcurrencyResult(t, r.ConcurrencyResult["testURL1"][0], 2, 2, 0, 0)
	checkConcurrencyResult(t, r.ConcurrencyResult["testURL1"][1], 1, 1, 0, 0)
	checkConcurrencyResult(t, r.ConcurrencyResult["testURL2"][0], 2, 1, 1, 0)
	checkConcurrencyResult(t, r.ConcurrencyResult["testURL2"][1], 1, 0, 1, 0)
	checkConcurrencyResult(t, r.ConcurrencyResult["testURL3"][0], 1, 0, 1, 0)

	checkURLs(t, r, []string{"testURL1", "testURL2", "testURL3"})
}

func checkConcurrencyResult(t *testing.T,
	result *ConcurrencyResult,
	TotalRequests,
	SuccessfulRequests,
	FailedRequests,
	TimedOutRequests int) {

	if result.FailedRequests != FailedRequests {
		t.Errorf("Expected to get %d for FailedRequests but got %d", FailedRequests, result.FailedRequests)
	}

	if result.TotalRequests != TotalRequests {
		t.Errorf("Expected to get %d for TotalRequests but got %d", TotalRequests, result.TotalRequests)
	}

	if result.SuccessfulRequests != SuccessfulRequests {
		t.Errorf("Expected to get %d for SuccessfulRequests but got %d", SuccessfulRequests, result.SuccessfulRequests)
	}

	if result.TimedOutRequests != TimedOutRequests {
		t.Errorf("Expected to get %d for TimedOutRequests but got %d", TimedOutRequests, result.TimedOutRequests)
	}
}

func checkURLs(t *testing.T, result *Result, expectedURLs []string) {
	for _, url := range expectedURLs {
		if _, ok := result.URLs[url]; !ok {
			t.Errorf("Expected to see %s in the list of urls", url)
		}
	}
}
