package cmd

import (
	"bytes"
	"errors"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/sasanrose/gbench/bench"
)

var testJSON = `{
    "host": "http://localhost:8080",
    "concurrency": 5,
    "requests": 100,
	"status-codes": [200, 201],
	"proxy": "test.proxy.url",
	"user": "user:pass",
	"cookie": "test-raw-cookie",
	"headers": ["X-Custom-Header: TestValue;"],
	"connect-timeout": 1000000000,
	"response-timeout": 5000000000,
	"paths": [
        {
            "path": "/"
        },
        {
            "path": "/test",
            "method": "post",
			"data": ["key1=val1"]
        }
    ]
}`

var testJSONNoHost = `{
    "concurrency": 5,
    "requests": 100,
	"paths": [
        {
            "path": "/"
        }
    ]
}`
var testJSONNoPath = `{
    "concurrency": 5,
    "requests": 100,
	"host": "localhost"
}`

func TestJSONConfig(t *testing.T) {
	oldFs := fs
	mfs := &mockedFSType{}
	fs = mfs

	defer func() {
		fs = oldFs
	}()

	mockedFile := &mockedFileType{bytes.NewBufferString(testJSON)}
	mfs.file = mockedFile

	configurations, err := getJSONConfig("testfile")

	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if mfs.openedName != "testfile" {
		t.Errorf("Wrong file opened")
	}

	if len(successStatusCodes) != 2 || successStatusCodes[0] != 200 || successStatusCodes[1] != 201 {
		t.Error("Unexpected status codes")
	}

	if concurrency != 5 {
		t.Error("Unexpected concurrency")
	}

	if requests != 100 {
		t.Error("Unexpected requests")
	}

	if authUserPass != "user:pass" {
		t.Error("Unexpected userpass")
	}

	if proxyURL != "test.proxy.url" {
		t.Error("Unexpected proxy URL")
	}

	if rawCookie != "test-raw-cookie" {
		t.Error("Unexpected raw cookie")
	}

	if len(headers) != 1 || headers[0] != "X-Custom-Header: TestValue;" {
		t.Errorf("Unecxpected Headers: %+v", headers)
	}

	if connectionTimeout != 1*time.Second || responseTimeout != 5*time.Second {
		t.Error("Unecxpected timeouts")
	}

	b := bench.NewBench(configurations...)

	checkBench(b, t)
}

func checkBench(b *bench.Bench, t *testing.T) {
	if len(b.URLs) != 2 {
		t.Error("Unexpected number of URLs")
	}

	if b.URLs[0].Addr != "http://localhost:8080/" || b.URLs[0].Method != http.MethodGet {
		t.Errorf("Unecxpected first URL: %+v", b.URLs[0])
	}

	if b.URLs[1].Addr != "http://localhost:8080/test" || b.URLs[1].Method != http.MethodPost {
		t.Errorf("Unecxpected first URL: %+v", b.URLs[1])
	}

	if val, ok := b.URLs[1].Data["key1"]; !ok || val != "val1" || len(b.URLs[1].Data) != 1 {
		t.Errorf("Unecxpected data for second URL: %+v", b.URLs[1])
	}
}

func TestWrongFile(t *testing.T) {
	testError(t, "Could not open \"Test file\": Test error", nil, errors.New("Test error"))
}

func TestNoHost(t *testing.T) {
	mockedFile := &mockedFileType{bytes.NewBufferString(testJSONNoHost)}
	testError(t, "No host is provided", mockedFile, nil)
}

func TestNoPaths(t *testing.T) {
	mockedFile := &mockedFileType{bytes.NewBufferString(testJSONNoPath)}
	testError(t, "No path is provided", mockedFile, nil)
}

func testError(t *testing.T, expectedErrorMsg string, mockedFile *mockedFileType, mockedError error) {
	oldFs := fs
	mfs := &mockedFSType{}
	fs = mfs

	defer func() {
		fs = oldFs
	}()

	mfs.file = mockedFile

	if mockedError != nil {
		mfs.err = mockedError
	}

	_, err := getJSONConfig("Test file")

	if err == nil || err.Error() != expectedErrorMsg {
		t.Errorf("Expected to get error %q but got %v", expectedErrorMsg, err)
	}
}

func TestNoFilePath(t *testing.T) {
	if os.Getenv("CRASH_TEST") == "1" {
		runJSON(jsonCmd, []string{})
		return
	}

	testExit(
		t,
		"TestNoFile",
		"",
	)
}
