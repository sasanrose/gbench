package bench

import (
	"bytes"
	"errors"
	"strings"
	"testing"
)

type urlTest struct {
	testStr     string
	expectedUrl *Url
}

var expectedUrls = []urlTest{
	urlTest{
		"GET|http://www.google.com?search=test",
		&Url{
			Addr: "http://www.google.com?search=test", Method: "GET",
		},
	},
	urlTest{
		"POST|https://www.google.com|search=test",
		&Url{
			Addr: "https://www.google.com", Method: "POST", Data: map[string]string{"search": "test"},
		},
	},
	urlTest{
		"POST|https://www.google.com?query=string|search=test&foo=bar",
		&Url{
			Addr:   "https://www.google.com?query=string",
			Method: "POST",
			Data: map[string]string{
				"search": "test",
				"foo":    "bar",
			},
		},
	},
	urlTest{
		"HEAD|http://www.google.com",
		&Url{
			Addr: "http://www.google.com", Method: "HEAD",
		},
	},
}

type mockedFSType struct {
	err  error // An arbitary error
	file *mockedFileType
}

type mockedFileType struct {
	mockedReader *bytes.Buffer // A mocked reader to return what we want to test
}

func (m mockedFileType) Read(p []byte) (n int, err error) {
	return m.mockedReader.Read(p)
}

func (m mockedFileType) Close() error {
	return nil
}

func (m mockedFSType) Open(name string) (file, error) {
	if m.err != nil {
		return nil, m.err
	}

	return m.file, nil
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
		t.Errorf("Expected to recieve error 'Test error' but recieved %v", err)
	}

	if f != nil {
		t.Error("Expected to recieve a nil funcion")
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

	for _, expectedUrl := range expectedUrls {
		fileContent += expectedUrl.testStr + "\n"
	}

	mockedFile := &mockedFileType{bytes.NewBufferString(fileContent)}
	mfs.file = mockedFile

	f, err := WithFile("/test/file")

	if err != nil {
		t.Errorf("Expected no error but recieved %v", err)
	}

	b := NewBench(f)

	checkTestResult(b, t)
}

func TestUrlString(t *testing.T) {
	configurations := []func(*Bench){}

	for _, expectedUrl := range expectedUrls {
		urlConfig, err := WithURLString(expectedUrl.testStr)

		if err != nil {
			t.Errorf("Unexpected error in parsing %s: %v", expectedUrl.testStr, err)
		}

		configurations = append(configurations, urlConfig)
	}

	b := NewBench(configurations...)

	checkTestResult(b, t)
}

func checkTestResult(b *Bench, t *testing.T) {
	if len(b.Urls) != len(expectedUrls) {
		t.Fatalf("Wrong number of urls")
	}

	for index, url := range expectedUrls {
		addr := b.Urls[index].Addr
		method := b.Urls[index].Method
		data := b.Urls[index].Data

		if url.expectedUrl.Addr != addr {
			t.Errorf("Expected address %s but got %s on url index %d", url.expectedUrl.Addr, addr, index)
		}

		if url.expectedUrl.Method != method {
			t.Errorf("Expected method %s but got %s on url index %d", url.expectedUrl.Method, method, index)
		}

		if (url.expectedUrl.Data == nil || len(url.expectedUrl.Data) == 0) && len(data) > 0 {
			t.Fatalf("Did not expect any data at url index %d but got %v", index, data)
		}

		if url.expectedUrl.Data != nil && len(url.expectedUrl.Data) != 0 {
			if len(url.expectedUrl.Data) != len(data) {
				t.Fatalf("Wrong number of data at url index %d", index)
			}

			for k, expectedValue := range url.expectedUrl.Data {
				if val, ok := data[k]; !ok || val != expectedValue {
					t.Errorf("Expected to have %s=%s in the data", k, expectedValue)
				}
			}
		}
	}
}

func TestWrongUrls(t *testing.T) {
	failTestCases := map[string]string{
		"wrong":                                      "Wrong URL format",
		"part1|part2|part3|wrongpart":                "Wrong URL format",
		"wrongmethod|part2|part3":                    "Not an allowed method provided",
		"GET|wrongurl":                               "Invalid URL provided",
		"GET|ftp://www.google.com":                   "Only http and https schemes are supported",
		"GET|https://www.google.com|part3":           "GET and HEAD do not need any data",
		"POST|https://www.google.com":                "You need to provide post data",
		"POST|https://www.google.com|test":           "Wrong key value format for post data",
		"POST|https://www.google.com|key1=val1&key2": "Wrong key value format for post data",
	}

	for testCase, expectedError := range failTestCases {
		parsedUrl, err := WithURLString(testCase)

		if err == nil || !strings.Contains(err.Error(), expectedError) {
			t.Errorf(`Expected to get an error containing "%s" for url string "%s"`, expectedError, testCase)
		}

		if parsedUrl != nil {
			t.Errorf("Expected to get a nil in resp for %s", testCase)
		}
	}
}
