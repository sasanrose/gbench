package bench

import (
	"bytes"
	"errors"
	"strings"
	"testing"
)

type urlTest struct {
	testStr     string
	expectedURL *URL
}

var expectedURLs = []urlTest{
	{
		"GET|http://www.google.com?search=test",
		&URL{
			Addr: "http://www.google.com?search=test", Method: "GET",
		},
	},
	{
		"POST|https://www.google.com|search=test",
		&URL{
			Addr: "https://www.google.com", Method: "POST", Data: map[string]string{"search": "test"},
		},
	},
	{
		"POST|https://www.google.com?query=string|search=test&foo=bar",
		&URL{
			Addr:   "https://www.google.com?query=string",
			Method: "POST",
			Data: map[string]string{
				"search": "test",
				"foo":    "bar",
			},
		},
	},
	{
		"HEAD|http://www.google.com",
		&URL{
			Addr: "http://www.google.com", Method: "HEAD",
		},
	},
}

func TestFileError(t *testing.T) {
	oldFs := fs
	mfs := &mockedFSType{}
	fs = mfs

	defer func() {
		fs = oldFs
	}()

	mfs.err = errors.New("Test error")

	f, err := WithFile("/test/file")

	if err == nil || err.Error() != "Test error" {
		t.Errorf("Expected to receive error 'Test error' but received %v", err)
	}

	if f != nil {
		t.Error("Expected to receive a nil funcion")
	}
}

func TestFile(t *testing.T) {
	oldFs := fs
	mfs := &mockedFSType{}
	fs = mfs

	defer func() {
		fs = oldFs
	}()

	fileContent := ""

	for _, expectedURL := range expectedURLs {
		fileContent += expectedURL.testStr + "\n"
	}

	mockedFile := &mockedFileType{bytes.NewBufferString(fileContent)}
	mfs.file = mockedFile

	f, err := WithFile("/test/file")

	if err != nil {
		t.Errorf("Expected no error but received %v", err)
	}

	b := NewBench(f)

	checkTestResult(b, t)
}

func TestFileNoURLs(t *testing.T) {
	oldFs := fs
	mfs := &mockedFSType{}
	fs = mfs

	defer func() {
		fs = oldFs
	}()

	fileContent := ""

	mockedFile := &mockedFileType{bytes.NewBufferString(fileContent)}
	mfs.file = mockedFile

	_, err := WithFile("/test/file")

	if err == nil || !strings.Contains(err.Error(), "Did not find any url in the") {
		t.Errorf("Expected to get an error for empty urls slice")
	}
}

func TestURLString(t *testing.T) {
	configurations := []func(*Bench){}

	for _, expectedURL := range expectedURLs {
		urlConfig, err := WithURLString(expectedURL.testStr)

		if err != nil {
			t.Errorf("Unexpected error in parsing %s: %v", expectedURL.testStr, err)
		}

		configurations = append(configurations, urlConfig)
	}

	b := NewBench(configurations...)

	checkTestResult(b, t)
}

func checkTestResult(b *Bench, t *testing.T) {
	if len(b.URLs) != len(expectedURLs) {
		t.Fatalf("Wrong number of urls")
	}

	for index, url := range expectedURLs {
		addr := b.URLs[index].Addr
		method := b.URLs[index].Method
		data := b.URLs[index].Data

		if url.expectedURL.Addr != addr {
			t.Errorf("Expected address %s but got %s on url index %d", url.expectedURL.Addr, addr, index)
		}

		if url.expectedURL.Method != method {
			t.Errorf("Expected method %s but got %s on url index %d", url.expectedURL.Method, method, index)
		}

		if (url.expectedURL.Data == nil || len(url.expectedURL.Data) == 0) && len(data) > 0 {
			t.Fatalf("Did not expect any data at url index %d but got %v", index, data)
		}

		if url.expectedURL.Data != nil && len(url.expectedURL.Data) != 0 {
			if len(url.expectedURL.Data) != len(data) {
				t.Fatalf("Wrong number of data at url index %d", index)
			}

			for k, expectedValue := range url.expectedURL.Data {
				if val, ok := data[k]; !ok || val != expectedValue {
					t.Errorf("Expected to have %s=%s in the data", k, expectedValue)
				}
			}
		}
	}
}

func TestWrongURLs(t *testing.T) {
	failTestCases := map[string]string{
		"wrong":                                      "Wrong URL format",
		"part1|part2|part3|wrongpart":                "Wrong URL format",
		"wrongmethod|part2|part3":                    "Method not allowed",
		"GET|ftp://www.google.com":                   "Only http and https schemes are supported",
		"GET|https://www.google.com|part3":           "GET and HEAD do not need any data",
		"POST|https://www.google.com":                "You need to provide post data",
		"POST|https://www.google.com|test":           "Wrong key value format for post data",
		"POST|https://www.google.com|key1=val1&key2": "Wrong key value format for post data",
	}

	for testCase, expectedError := range failTestCases {
		parsedURL, err := WithURLString(testCase)

		if err == nil || !strings.Contains(err.Error(), expectedError) {
			t.Errorf(`Expected to get an error containing "%s" for url string "%s"`, expectedError, testCase)
		}

		if parsedURL != nil {
			t.Errorf("Expected to get a nil in resp for %s", testCase)
		}
	}
}
