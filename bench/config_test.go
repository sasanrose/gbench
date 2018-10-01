package bench

import (
	"bytes"
	"testing"
	"time"

	"github.com/sasanrose/gbench/report"
)

func TestConfigurations(t *testing.T) {
	var buf bytes.Buffer

	u := &URL{
		Addr: "testAddr", Method: "GET",
	}

	r := &report.Result{}

	_, err := WithHeaderString("wrongformat")

	if err == nil {
		t.Error("Expected to get an error for wrong header format")
	}

	hConfig, err := WithHeaderString("fooKey=barVal")

	if err != nil {
		t.Error("Unexpected error for correct header format")
	}

	configurations := []func(*Bench){
		WithConcurrency(2),
		WithRequests(4),
		WithURL(u),
		WithAuth("user", "pass"),
		WithVerbosity(&buf),
		WithProxy("testProxy"),
		WithConnectionTimeout(2 * time.Second),
		WithResponseTimeout(3 * time.Second),
		WithRawCookie("testCookie"),
		WithHeader("testKey", "testVal"),
		hConfig,
		WithReport(r),
		WithSuccessStatusCode(100),
		WithSuccessStatusCode(101),
	}

	b := NewBench(configurations...)

	checkBenchGeneralConfig(b, t)
	checkBenchURLs(b, t)
}

func checkBenchGeneralConfig(b *Bench, t *testing.T) {
	if b.Concurrency != 2 {
		t.Error("Concurrency is not set as expected")
	}

	if b.Requests != 4 {
		t.Error("Number of requests is not set as expected")
	}

	if b.Auth.Username != "user" || b.Auth.Password != "pass" {
		t.Error("Auth is not set as expected")
	}

	if b.VerbosityWriter == nil {
		t.Error("Verbosity writer is not set as expected")
	}

	if b.Proxy != "testProxy" {
		t.Error("Proxy is not set as expected")
	}

	if b.ConnectionTimeout != 2*time.Second {
		t.Error("Connection timeout is not set as expected")
	}

	if b.ResponseTimeout != 3*time.Second {
		t.Error("Response timeout is not set as expected")
	}

	if b.Report == nil {
		t.Error("Report is not set as expected")
	}
}

func checkBenchURLs(b *Bench, t *testing.T) {
	if len(b.URLs) != 1 || b.URLs[0].Addr != "testAddr" || b.URLs[0].Method != "GET" {
		t.Error("URL is not set as expected")
	}

	if b.RawCookie != "testCookie" {
		t.Error("Raw cookie is not set as expected")
	}

	if val, ok := b.Headers["testKey"]; !ok || val != "testVal" {
		t.Error("Header is not set as expected")
	}

	if val, ok := b.Headers["fooKey"]; !ok || val != "barVal" {
		t.Error("Header is not set as expected")
	}

	if len(b.SuccessStatusCodes) != 2 || b.SuccessStatusCodes[0] != 100 || b.SuccessStatusCodes[1] != 101 {
		t.Errorf("Success status codes were expected to be set as '100' and '101' but got %v", b.SuccessStatusCodes)
	}
}
