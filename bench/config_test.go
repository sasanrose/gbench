package bench

import (
	"bytes"
	"testing"
	"time"

	"github.com/sasanrose/gbench/report"
)

func TestHeaders(t *testing.T) {
	_, err := WithHeaderString("wrongformat")

	if err == nil {
		t.Error("Expected to get an error for wrong header format")
	}

	_, err = WithHeaderString("Accept:text/html")

	if err != nil {
		t.Error("Unexpected error for correct header format")
	}

	_, err = WithHeaderString("Accept: text/html")

	if err != nil {
		t.Error("Unexpected error for correct header format")
	}

	_, err = WithHeaderString("Accept: text/html;")

	if err != nil {
		t.Error("Unexpected error for correct header format")
	}

	_, err = WithHeaderString("X-CustomHeader;")

	if err != nil {
		t.Error("Unexpected error for correct header format")
	}
}

func TestAuthUserPass(t *testing.T) {
	_, err := WithAuthUserPass("user:pass")

	if err != nil {
		t.Errorf("Did not expect error for 'user:pass' but got: %v", err)
	}

	_, err = WithAuthUserPass("user:pass@1234:1234")

	if err != nil {
		t.Errorf("Did not expect error for 'user:pass@1234:1234' but got: %v", err)
	}

	_, err = WithAuthUserPass("usernopass")

	if err == nil {
		t.Error("Expected to get an error for 'usernopass' but got nothing")
	}

	config, _ := WithAuthUserPass("user:pass")

	if config == nil {
		t.Error("Expected to get a config func for 'user:pass'")
	}
}

func TestConfigurations(t *testing.T) {
	var buf bytes.Buffer

	u := &URL{
		Addr: "testAddr", Method: "GET",
	}

	r := &report.Result{}

	configurations := []func(*Bench){
		WithConcurrency(2),
		WithRequests(4),
		WithURL(u),
		WithAuth("user", "pass"),
		WithOutput(&buf),
		WithProxy("testProxy"),
		WithConnectionTimeout(2 * time.Second),
		WithResponseTimeout(3 * time.Second),
		WithRawCookie("testCookie"),
		WithHeader("testKey", "testVal"),
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

	if b.OutputWriter == nil {
		t.Error("Output writer is not set as expected")
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

	if len(b.SuccessStatusCodes) != 2 || b.SuccessStatusCodes[0] != 100 || b.SuccessStatusCodes[1] != 101 {
		t.Errorf("Success status codes were expected to be set as '100' and '101' but got %v", b.SuccessStatusCodes)
	}
}
