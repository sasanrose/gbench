package result

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
	if r.receivedDataLength["testUrl1"] != 10 {
		t.Fatal("Expected to get '10' for 'testUrl1'")
	}

	r.AddReceivedDataLength("testUrl1", -1)
	if r.receivedDataLength["testUrl1"] != 10 {
		t.Fatal("Expected to get '10' for 'testUrl1'")
	}

	r.AddReceivedDataLength("testUrl1", 2)
	if r.receivedDataLength["testUrl1"] != 12 {
		t.Fatal("Expected to get '12' for 'testUrl1'")
	}

	r.AddReceivedDataLength("testUrl2", 20)
	if r.receivedDataLength["testUrl2"] != 20 {
		t.Fatal("Expected to get '20' for 'testUrl2'")
	}

	if r.receivedDataLength["testUrl1"] != 12 {
		t.Fatal("Expected to get '12' for 'testUrl1'")
	}

	if r.totalReceivedDataLength != 32 {
		t.Fatal("Expected to get 32 as totalReceivedDataLength")
	}
}

func TestResponseTime(t *testing.T) {
	r := getTestResultStruct()

	r.AddResponseTime("testUrl1", 10*time.Second)
	if r.responseTime["testUrl1"] != 10*time.Second {
		t.Fatal("Expected to get '10' seconds for 'testUrl1'")
	}

	r.AddResponseTime("testUrl1", 2*time.Second)
	if r.responseTime["testUrl1"] != 12*time.Second {
		t.Fatal("Expected to get '12' seconds for 'testUrl1'")
	}

	r.AddResponseTime("testUrl2", 20*time.Second)
	if r.responseTime["testUrl2"] != 20*time.Second {
		t.Fatal("Expected to get '20' seconds for 'testUrl2'")
	}

	if r.totalResponseTime != 32*time.Second {
		t.Fatal("Expected to get 32 as totalResponseTime")
	}

	if r.responseTime["testUrl1"] != 12*time.Second {
		t.Fatal("Expected to get '12' seconds for 'testUrl1'")
	}

	if r.shortestResponseTime != 2*time.Second {
		t.Fatalf("Expected to get '2' as the shortest response time but got %v", r.shortestResponseTime)
	}

	if r.longestResponseTime != 20*time.Second {
		t.Fatalf("Expected to get '20' as the longest response time but got %v", r.longestResponseTime)
	}

	if r.shortestResponseTimes["testUrl1"] != 2*time.Second {
		t.Fatalf("Expected to get '2' as the shortest response time for 'testUrl1' but got %v", r.shortestResponseTimes["testUrl1"])
	}

	if r.shortestResponseTimes["testUrl2"] != 20*time.Second {
		t.Fatalf("Expected to get '20' as the shortest response time for 'testUrl2' but got %v", r.shortestResponseTimes["testUrl2"])
	}

	if r.longestResponseTimes["testUrl1"] != 10*time.Second {
		t.Fatalf("Expected to get '2' as the longest response time for 'testUrl1' but got %v", r.longestResponseTimes["testUrl1"])
	}

	if r.longestResponseTimes["testUrl2"] != 20*time.Second {
		t.Fatalf("Expected to get '20' as the longest response time for 'testUrl2' but got %v", r.longestResponseTimes["testUrl2"])
	}

	if r.responseTimesTotalCount != 3 {
		t.Errorf("Expected to get 3 for responseTimesTotalCount but got %d", r.responseTimesTotalCount)
	}

	if r.responseTimesCount["testUrl1"] != 2 {
		t.Errorf("Expected to get 2 for responseTimesCount for testUrl1 but got %d", r.responseTimesCount["testUrl1"])
	}

	if r.responseTimesCount["testUrl2"] != 1 {
		t.Errorf("Expected to get 1 for responseTimesCount for testUrl2 but got %d", r.responseTimesCount["testUrl1"])
	}
}

func TestTotalTime(t *testing.T) {
	r := getTestResultStruct()
	r.SetTotalDuration(100 * time.Second)
	if r.totalTime != 100*time.Second {
		t.Fatal("Expected to get '100' seconds for total time")
	}
}

func TestTimedoutResponse(t *testing.T) {
	r := getTestResultStruct()

	r.AddTimedoutResponse("testUrl1")
	if r.timedoutResponse["testUrl1"] != 1 {
		t.Fatal("Expected to get '1' for 'testUrl1'")
	}

	r.AddTimedoutResponse("testUrl1")
	if r.timedoutResponse["testUrl1"] != 2 {
		t.Fatal("Expected to get '2' for 'testUrl1'")
	}

	r.AddTimedoutResponse("testUrl2")
	if r.timedoutResponse["testUrl2"] != 1 {
		t.Fatal("Expected to get '1' for 'testUrl2'")
	}

	if r.timedoutResponse["testUrl1"] != 2 {
		t.Fatal("Expected to get '2' for 'testUrl1'")
	}

	if r.totalRequests != 3 || r.timedOutRequests != 3 {
		t.Error("Number of total requets and timed out requests are not set as expected")
	}

	if len(r.concurrencyResult["testUrl1"]) != 1 || len(r.concurrencyResult["testUrl2"]) != 1 {
		t.Error("Wrong concurrencyResult")
	}

	checkConcurrencyResult(t, r.concurrencyResult["testUrl1"][0], 2, 0, 0, 2)
	checkConcurrencyResult(t, r.concurrencyResult["testUrl2"][0], 1, 0, 0, 1)

	checkUrls(t, r, []string{"testUrl1", "testUrl2"})
}

func TestFailedResponse(t *testing.T) {
	r := getTestResultStruct()

	r.AddFailedResponse("testUrl1")
	if r.failedResponse["testUrl1"] != 1 {
		t.Fatal("Expected to get '1' for 'testUrl1'")
	}

	r.AddFailedResponse("testUrl1")
	if r.failedResponse["testUrl1"] != 2 {
		t.Fatal("Expected to get '2' for 'testUrl1'")
	}

	r.AddFailedResponse("testUrl2")
	if r.failedResponse["testUrl2"] != 1 {
		t.Fatal("Expected to get '1' for 'testUrl2'")
	}

	if r.failedResponse["testUrl1"] != 2 {
		t.Fatal("Expected to get '2' for 'testUrl1'")
	}

	if r.totalRequests != 3 || r.failedRequests != 3 {
		t.Error("Number of total requets and failed requests are not set as expected")
	}

	if len(r.concurrencyResult["testUrl1"]) != 1 || len(r.concurrencyResult["testUrl2"]) != 1 {
		t.Error("Wrong concurrencyResult")
	}

	checkConcurrencyResult(t, r.concurrencyResult["testUrl1"][0], 2, 0, 2, 0)
	checkConcurrencyResult(t, r.concurrencyResult["testUrl2"][0], 1, 0, 1, 0)

	checkUrls(t, r, []string{"testUrl1", "testUrl2"})
}

func TestStatusCode(t *testing.T) {
	r := getTestResultStruct()

	r.AddResponseStatusCode("testUrl1", 200, false)
	if r.responseStatusCode["testUrl1"][200] != 1 {
		t.Fatal("Expected count '1' for status cude '200' for testUrl1")
	}

	r.AddResponseStatusCode("testUrl1", 200, false)
	if r.responseStatusCode["testUrl1"][200] != 2 {
		t.Fatal("Expected count '2' for status cude '200' for testUrl1")
	}

	r.AddResponseStatusCode("testUrl2", 200, false)
	if r.responseStatusCode["testUrl2"][200] != 1 {
		t.Fatal("Expected count '1' for status cude '200' for testUrl2")
	}

	r.AddResponseStatusCode("testUrl1", 201, false)
	if r.responseStatusCode["testUrl1"][201] != 1 {
		t.Fatal("Expected count '1' for status cude '201' for testUrl1")
	}

	if len(r.failedResponseStatusCode) != 0 {
		t.Fatal("Expected empty for failedResponseStatusCode")
	}

	r.AddResponseStatusCode("testUrl2", 500, true)
	if r.failedResponseStatusCode["testUrl2"][500] != 1 {
		t.Fatal("Expected count '1' for status cude '500' for testUrl2")
	}

	r.AddResponseStatusCode("testUrl2", 500, true)
	if r.failedResponseStatusCode["testUrl2"][500] != 2 {
		t.Fatal("Expected count '2' for status cude '500' for testUrl2")
	}

	r.AddResponseStatusCode("testUrl3", 500, true)
	if r.failedResponseStatusCode["testUrl3"][500] != 1 {
		t.Fatal("Expected count '1' for status cude '500' for testUrl3")
	}

	if r.totalRequests != 7 || r.successfulRequests != 4 || r.failedRequests != 3 {
		t.Error("Number of total requests, successful requests and failed requests are not set as expected")
	}

	if len(r.concurrencyResult["testUrl1"]) != 2 ||
		len(r.concurrencyResult["testUrl2"]) != 2 ||
		len(r.concurrencyResult["testUrl3"]) != 1 {
		t.Error("Wrong concurrencyResult")
	}

	checkConcurrencyResult(t, r.concurrencyResult["testUrl1"][0], 2, 2, 0, 0)
	checkConcurrencyResult(t, r.concurrencyResult["testUrl1"][1], 1, 1, 0, 0)
	checkConcurrencyResult(t, r.concurrencyResult["testUrl2"][0], 2, 1, 1, 0)
	checkConcurrencyResult(t, r.concurrencyResult["testUrl2"][1], 1, 0, 1, 0)
	checkConcurrencyResult(t, r.concurrencyResult["testUrl3"][0], 1, 0, 1, 0)

	checkUrls(t, r, []string{"testUrl1", "testUrl2", "testUrl3"})
}

func checkConcurrencyResult(t *testing.T,
	result *concurrencyResult,
	totalRequests,
	successfulRequests,
	failedRequests,
	timedOutRequests int) {

	if result.failedRequests != failedRequests {
		t.Errorf("Expected to get %d for failedRequests but got %d", failedRequests, result.failedRequests)
	}

	if result.totalRequests != totalRequests {
		t.Errorf("Expected to get %d for totalRequests but got %d", totalRequests, result.totalRequests)
	}

	if result.successfulRequests != successfulRequests {
		t.Errorf("Expected to get %d for successfulRequests but got %d", successfulRequests, result.successfulRequests)
	}

	if result.timedOutRequests != timedOutRequests {
		t.Errorf("Expected to get %d for timedOutRequests but got %d", timedOutRequests, result.timedOutRequests)
	}
}

func checkUrls(t *testing.T, result *Result, expectedUrls []string) {
	for _, url := range expectedUrls {
		if _, ok := result.urls[url]; !ok {
			t.Errorf("Expected to see %s in the list of urls", url)
		}
	}
}
