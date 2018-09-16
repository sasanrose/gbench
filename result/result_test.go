package result

import (
	"testing"
	"time"
)

func getTestResultStruct() *result {
	r := &result{}
	r.Init()

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

	if r.responseTime["testUrl1"] != 12*time.Second {
		t.Fatal("Expected to get '12' seconds for 'testUrl1'")
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

	r.AddResponseStatusCode("testUrl1", 201, false)
	if r.responseStatusCode["testUrl1"][201] != 1 {
		t.Fatal("Expected count '1' for status cude '201' for testUrl1")
	}

	r.AddResponseStatusCode("testUrl2", 200, false)
	if r.responseStatusCode["testUrl2"][200] != 1 {
		t.Fatal("Expected count '1' for status cude '200' for testUrl2")
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
}
