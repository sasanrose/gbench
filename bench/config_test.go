package bench

import (
	"bytes"
	"testing"
	"time"
)

func TestConfigurations(t *testing.T) {
	var buf bytes.Buffer

	u := &Url{
		Addr: "testAddr", Method: "GET",
	}

	r := &mockedRenderer{}

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
		WithRenderer(r),
		WithSuccessStatusCode(100),
		WithSuccessStatusCode(101),
	}

	b := NewBench(configurations...)

	if b.Concurrency != 2 {
		t.Error("Concurrency is not set as expected")
	}

	if b.Requests != 4 {
		t.Error("Number of requests is not set as expected")
	}

	if len(b.Urls) != 1 || b.Urls[0].Addr != "testAddr" || b.Urls[0].Method != "GET" {
		t.Error("URL is not set as expected")
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

	if b.RawCookie != "testCookie" {
		t.Error("Raw cookie is not set as expected")
	}

	if val, ok := b.Headers["testKey"]; !ok || val != "testVal" {
		t.Error("Header is not set as expected")
	}

	if b.Renderer == nil {
		t.Error("Renderer is not set as expected")
	}

	if len(b.SuccessStatusCodes) != 2 || b.SuccessStatusCodes[0] != 100 || b.SuccessStatusCodes[1] != 101 {
		t.Errorf("Success status codes were expected to be set as '100' and '101' but got %v", b.SuccessStatusCodes)
	}
}
