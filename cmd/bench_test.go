package cmd

import (
	"os"
	"os/exec"
	"testing"
	"time"

	"github.com/sasanrose/gbench/bench"
	"github.com/sasanrose/gbench/report"
)

func TestGlobalConfigurations(t *testing.T) {
	result := setSharedVars()

	headers = []string{"X-CustomHeader: Value;"}
	authUserPass = "user:pass"

	configurations := make([]func(*bench.Bench), 0)
	configurations, _ = appendGlobalConfigurations(configurations, result)

	b := bench.NewBench(configurations...)

	if b.Concurrency != concurrency {
		t.Errorf("Expected concurrency of %d but got %d", concurrency, b.Concurrency)
	}

	if b.Requests != requests {
		t.Errorf("Expected requests of %d but got %d", requests, b.Requests)
	}

	if b.ConnectionTimeout != connectionTimeout {
		t.Errorf("Expected ConnectionTimeout of %s but got %s", connectionTimeout, b.ConnectionTimeout)
	}

	if b.ResponseTimeout != responseTimeout {
		t.Errorf("Expected ResponseTimeout of %s but got %s", responseTimeout, b.ResponseTimeout)
	}

	if b.Proxy != proxyURL {
		t.Errorf("Expected Proxy of %s but got %s", proxyURL, b.Proxy)
	}

	if b.RawCookie != rawCookie {
		t.Errorf("Expected RawCookie of %s but got %s", rawCookie, b.RawCookie)
	}

	if b.OutputWriter != os.Stdout {
		t.Error("Expected output writer to be os.Stdout")
	}

	if b.Report != result {
		t.Error("Unexpected report.Result")
	}

	if len(b.SuccessStatusCodes) != len(successStatusCodes) {
		t.Errorf("Expected SuccessStatusCodes of %v but got %v", successStatusCodes, b.SuccessStatusCodes)
	}

	for index, statusCode := range successStatusCodes {
		if b.SuccessStatusCodes[index] != statusCode {
			t.Errorf("Expected SuccessStatusCode of %d at index %d but got %d", statusCode, index, b.SuccessStatusCodes[index])
		}
	}
}

func TestWrongHeader(t *testing.T) {
	expected := "Error with header: WrongHeader is not a correct 'key;' format"
	result := setSharedVars()
	headers = []string{"WrongHeader"}

	configurations := make([]func(*bench.Bench), 0)
	configurations, err := appendGlobalConfigurations(configurations, result)

	if err == nil || err.Error() != expected {
		t.Errorf("Expected to get %q but got %v", expected, err)
	}
}

func TestWrongAuth(t *testing.T) {
	expected := "Error with authentication credentials: Wrong auth credentials format: wronguserpass"
	result := setSharedVars()
	headers = []string{}
	authUserPass = "wronguserpass"

	configurations := make([]func(*bench.Bench), 0)
	configurations, err := appendGlobalConfigurations(configurations, result)

	if err == nil || err.Error() != expected {
		t.Errorf("Expected to get %q but got %v", expected, err)
	}
}

func setSharedVars() *report.Result {
	concurrency = 5
	requests = 100
	connectionTimeout = 1 * time.Second
	responseTimeout = 5 * time.Second
	successStatusCodes = []int{200, 201}
	proxyURL = "http://proxy.com"
	rawCookie = "rawCookie"

	return &report.Result{}
}

func testExit(t *testing.T, testName, errMsg string) {
	cmd := exec.Command(os.Args[0], "-test.run="+testName)
	cmd.Env = append(os.Environ(), "CRASH_TEST=1")
	_, err := cmd.Output()

	if e, ok := err.(*exec.ExitError); ok &&
		e.Exited() &&
		!e.Success() &&
		e.Error() == "exit status 2" {

		if errMsg != "" && string(e.Stderr) != errMsg {
			t.Fatalf("Expected error message %q but got %q", errMsg, e.Stderr)
		}

		return
	}

	t.Fatalf("process ran with err %v, want exit status 2", err)
}
