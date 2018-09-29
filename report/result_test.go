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

func TestContentLength(t *testing.T) {
	r := getTestResultStruct()

	r.AddReceivedDataLength("testUrl1", 10)
	if r.ReceivedDataLength["testUrl1"] != 10 {
		t.Fatal("Expected to get '10' for 'testUrl1'")
	}

	r.AddReceivedDataLength("testUrl1", -1)
	if r.ReceivedDataLength["testUrl1"] != 10 {
		t.Fatal("Expected to get '10' for 'testUrl1'")
	}

	r.AddReceivedDataLength("testUrl1", 2)
	if r.ReceivedDataLength["testUrl1"] != 12 {
		t.Fatal("Expected to get '12' for 'testUrl1'")
	}

	r.AddReceivedDataLength("testUrl2", 20)
	if r.ReceivedDataLength["testUrl2"] != 20 {
		t.Fatal("Expected to get '20' for 'testUrl2'")
	}

	if r.ReceivedDataLength["testUrl1"] != 12 {
		t.Fatal("Expected to get '12' for 'testUrl1'")
	}

	if r.TotalReceivedDataLength != 32 {
		t.Fatal("Expected to get 32 as TotalReceivedDataLength")
	}
}

func TestResponseTime(t *testing.T) {
	r := getTestResultStruct()

	r.AddResponseTime("testUrl1", 10*time.Second)
	if r.ResponseTime["testUrl1"] != 10*time.Second {
		t.Fatal("Expected to get '10' seconds for 'testUrl1'")
	}

	r.AddResponseTime("testUrl1", 2*time.Second)
	if r.ResponseTime["testUrl1"] != 12*time.Second {
		t.Fatal("Expected to get '12' seconds for 'testUrl1'")
	}

	r.AddResponseTime("testUrl2", 20*time.Second)
	if r.ResponseTime["testUrl2"] != 20*time.Second {
		t.Fatal("Expected to get '20' seconds for 'testUrl2'")
	}

	if r.TotalResponseTime != 32*time.Second {
		t.Fatal("Expected to get 32 as TotalResponseTime")
	}

	if r.ResponseTime["testUrl1"] != 12*time.Second {
		t.Fatal("Expected to get '12' seconds for 'testUrl1'")
	}

	if r.ShortestResponseTime != 2*time.Second {
		t.Fatalf("Expected to get '2' as the shortest response time but got %v", r.ShortestResponseTime)
	}

	if r.LongestResponseTime != 20*time.Second {
		t.Fatalf("Expected to get '20' as the longest response time but got %v", r.LongestResponseTime)
	}

	if r.ShortestResponseTimes["testUrl1"] != 2*time.Second {
		t.Fatalf("Expected to get '2' as the shortest response time for 'testUrl1' but got %v", r.ShortestResponseTimes["testUrl1"])
	}

	if r.ShortestResponseTimes["testUrl2"] != 20*time.Second {
		t.Fatalf("Expected to get '20' as the shortest response time for 'testUrl2' but got %v", r.ShortestResponseTimes["testUrl2"])
	}

	if r.LongestResponseTimes["testUrl1"] != 10*time.Second {
		t.Fatalf("Expected to get '2' as the longest response time for 'testUrl1' but got %v", r.LongestResponseTimes["testUrl1"])
	}

	if r.LongestResponseTimes["testUrl2"] != 20*time.Second {
		t.Fatalf("Expected to get '20' as the longest response time for 'testUrl2' but got %v", r.LongestResponseTimes["testUrl2"])
	}

	if r.ResponseTimesTotalCount != 3 {
		t.Errorf("Expected to get 3 for ResponseTimesTotalCount but got %d", r.ResponseTimesTotalCount)
	}

	if r.ResponseTimesCount["testUrl1"] != 2 {
		t.Errorf("Expected to get 2 for ResponseTimesCount for testUrl1 but got %d", r.ResponseTimesCount["testUrl1"])
	}

	if r.ResponseTimesCount["testUrl2"] != 1 {
		t.Errorf("Expected to get 1 for ResponseTimesCount for testUrl2 but got %d", r.ResponseTimesCount["testUrl1"])
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

	r.AddTimedoutResponse("testUrl1")
	if r.TimedoutResponse["testUrl1"] != 1 {
		t.Fatal("Expected to get '1' for 'testUrl1'")
	}

	r.AddTimedoutResponse("testUrl1")
	if r.TimedoutResponse["testUrl1"] != 2 {
		t.Fatal("Expected to get '2' for 'testUrl1'")
	}

	r.AddTimedoutResponse("testUrl2")
	if r.TimedoutResponse["testUrl2"] != 1 {
		t.Fatal("Expected to get '1' for 'testUrl2'")
	}

	if r.TimedoutResponse["testUrl1"] != 2 {
		t.Fatal("Expected to get '2' for 'testUrl1'")
	}

	if r.TotalRequests != 3 || r.TimedOutRequests != 3 {
		t.Error("Number of total requets and timed out requests are not set as expected")
	}

	if len(r.ConcurrencyResult["testUrl1"]) != 1 || len(r.ConcurrencyResult["testUrl2"]) != 1 {
		t.Error("Wrong ConcurrencyResult")
	}

	checkConcurrencyResult(t, r.ConcurrencyResult["testUrl1"][0], 2, 0, 0, 2)
	checkConcurrencyResult(t, r.ConcurrencyResult["testUrl2"][0], 1, 0, 0, 1)

	checkUrls(t, r, []string{"testUrl1", "testUrl2"})
}

func TestFailedResponse(t *testing.T) {
	r := getTestResultStruct()

	r.AddFailedResponse("testUrl1")
	if r.FailedResponse["testUrl1"] != 1 {
		t.Fatal("Expected to get '1' for 'testUrl1'")
	}

	r.AddFailedResponse("testUrl1")
	if r.FailedResponse["testUrl1"] != 2 {
		t.Fatal("Expected to get '2' for 'testUrl1'")
	}

	r.AddFailedResponse("testUrl2")
	if r.FailedResponse["testUrl2"] != 1 {
		t.Fatal("Expected to get '1' for 'testUrl2'")
	}

	if r.FailedResponse["testUrl1"] != 2 {
		t.Fatal("Expected to get '2' for 'testUrl1'")
	}

	if r.TotalRequests != 3 || r.FailedRequests != 3 {
		t.Error("Number of total requets and failed requests are not set as expected")
	}

	if len(r.ConcurrencyResult["testUrl1"]) != 1 || len(r.ConcurrencyResult["testUrl2"]) != 1 {
		t.Error("Wrong ConcurrencyResult")
	}

	checkConcurrencyResult(t, r.ConcurrencyResult["testUrl1"][0], 2, 0, 2, 0)
	checkConcurrencyResult(t, r.ConcurrencyResult["testUrl2"][0], 1, 0, 1, 0)

	checkUrls(t, r, []string{"testUrl1", "testUrl2"})
}

func TestStatusCode(t *testing.T) {
	r := getTestResultStruct()

	r.AddResponseStatusCode("testUrl1", 200, false)
	if r.ResponseStatusCode["testUrl1"][200] != 1 {
		t.Fatal("Expected count '1' for status cude '200' for testUrl1")
	}

	r.AddResponseStatusCode("testUrl1", 200, false)
	if r.ResponseStatusCode["testUrl1"][200] != 2 {
		t.Fatal("Expected count '2' for status cude '200' for testUrl1")
	}

	r.AddResponseStatusCode("testUrl2", 200, false)
	if r.ResponseStatusCode["testUrl2"][200] != 1 {
		t.Fatal("Expected count '1' for status cude '200' for testUrl2")
	}

	r.AddResponseStatusCode("testUrl1", 201, false)
	if r.ResponseStatusCode["testUrl1"][201] != 1 {
		t.Fatal("Expected count '1' for status cude '201' for testUrl1")
	}

	if len(r.FailedResponseStatusCode) != 0 {
		t.Fatal("Expected empty for FailedResponseStatusCode")
	}

	r.AddResponseStatusCode("testUrl2", 500, true)
	if r.FailedResponseStatusCode["testUrl2"][500] != 1 {
		t.Fatal("Expected count '1' for status cude '500' for testUrl2")
	}

	r.AddResponseStatusCode("testUrl2", 500, true)
	if r.FailedResponseStatusCode["testUrl2"][500] != 2 {
		t.Fatal("Expected count '2' for status cude '500' for testUrl2")
	}

	r.AddResponseStatusCode("testUrl3", 500, true)
	if r.FailedResponseStatusCode["testUrl3"][500] != 1 {
		t.Fatal("Expected count '1' for status cude '500' for testUrl3")
	}

	if r.TotalRequests != 7 || r.SuccessfulRequests != 4 || r.FailedRequests != 3 {
		t.Error("Number of total requests, successful requests and failed requests are not set as expected")
	}

	if len(r.ConcurrencyResult["testUrl1"]) != 2 ||
		len(r.ConcurrencyResult["testUrl2"]) != 2 ||
		len(r.ConcurrencyResult["testUrl3"]) != 1 {
		t.Error("Wrong ConcurrencyResult")
	}

	checkConcurrencyResult(t, r.ConcurrencyResult["testUrl1"][0], 2, 2, 0, 0)
	checkConcurrencyResult(t, r.ConcurrencyResult["testUrl1"][1], 1, 1, 0, 0)
	checkConcurrencyResult(t, r.ConcurrencyResult["testUrl2"][0], 2, 1, 1, 0)
	checkConcurrencyResult(t, r.ConcurrencyResult["testUrl2"][1], 1, 0, 1, 0)
	checkConcurrencyResult(t, r.ConcurrencyResult["testUrl3"][0], 1, 0, 1, 0)

	checkUrls(t, r, []string{"testUrl1", "testUrl2", "testUrl3"})
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

func checkUrls(t *testing.T, result *Result, expectedUrls []string) {
	for _, url := range expectedUrls {
		if _, ok := result.Urls[url]; !ok {
			t.Errorf("Expected to see %s in the list of urls", url)
		}
	}
}
