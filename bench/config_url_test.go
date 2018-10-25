package bench

import (
	"strings"
	"testing"
)

type testCase struct {
	testURL    string
	testMethod string
	testData   []string
	headers    []string
	rawCookie  string
	userPass   string
}

type urlTest struct {
	testCase
	expectedURL *URL
}

var expectedURLs = []urlTest{
	{
		testCase{
			"http://www.google.com?search=test",
			"GET",
			[]string{},
			[]string{},
			"",
			"",
		},
		&URL{
			Addr: "http://www.google.com?search=test", Method: "GET", Data: map[string]string{},
		},
	},
	{
		testCase{
			"https://www.google.com",
			"POST",
			[]string{"search=test"},
			[]string{},
			"",
			"",
		},
		&URL{
			Addr: "https://www.google.com", Method: "POST", Data: map[string]string{"search": "test"},
		},
	},
	{
		testCase{
			"https://www.google.com?query=string",
			"POST",
			[]string{"search=test", "foo=bar"},
			[]string{"header1: value1", "header2: value2;"},
			"",
			"username:password",
		},
		&URL{
			Addr:   "https://www.google.com?query=string",
			Method: "POST",
			Data: map[string]string{
				"search": "test",
				"foo":    "bar",
			},
			Headers: map[string]string{
				"header1": "value1",
				"header2": "value2",
			},
			Auth: &Auth{"username", "password"},
		},
	},
	{
		testCase{
			"http://www.google.com",
			"HEAD",
			[]string{},
			[]string{},
			"cookie",
			"",
		},
		&URL{
			Addr: "http://www.google.com", Method: "HEAD", Data: map[string]string{}, RawCookie: "cookie",
		},
	},
}

var failTestCases = []struct {
	testCase
	expectedErrorMsg string
}{
	{
		testCase{
			"ftp://www.google.com",
			"GET",
			[]string{},
			[]string{},
			"",
			"",
		},
		"Only http and https schemes are supported",
	},
	{
		testCase{
			"https://www.google.com",
			"GET",
			[]string{"foo=bar"},
			[]string{},
			"",
			"",
		},
		"Request data is only allowed with POST, PUT and PATCH request methods",
	},
	{
		testCase{
			"https://www.google.com",
			"POST",
			[]string{"test"},
			[]string{},
			"",
			"",
		},
		"Wrong key value format for request data",
	},
	{
		testCase{
			"https://www.google.com",
			"POST",
			[]string{"key1=val1&key2"},
			[]string{},
			"",
			"",
		},
		"Wrong key value format for request data",
	},
	{
		testCase{
			"https://www.google.com",
			"POST",
			[]string{},
			[]string{"Header"},
			"",
			"",
		},
		"Header is not a correct 'key;' format",
	},
	{
		testCase{
			"https://www.google.com",
			"POST",
			[]string{},
			[]string{"header: value1; header2: value2"},
			"",
			"",
		},
		"header: value1; header2: value2 is not a correct 'key: value;' format",
	},
	{
		testCase{
			"https://www.google.com",
			"GET",
			[]string{},
			[]string{},
			"",
			"user",
		},
		"Wrong auth credentials format: user",
	},
}

/*func TestFileError(t *testing.T) {
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
}*/

/*func TestFile(t *testing.T) {
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
}*/

/*func TestFileNoURLs(t *testing.T) {
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
}*/

func TestURLString(t *testing.T) {
	configurations := []func(*Bench){}

	for _, expectedURL := range expectedURLs {
		urlConfig, err := WithURLSettings(expectedURL.testURL,
			expectedURL.testMethod,
			expectedURL.testData,
			expectedURL.headers,
			expectedURL.rawCookie,
			expectedURL.userPass)

		if err != nil {
			t.Fatalf("Unexpected error in parsing %s: %v", expectedURL.testURL, err)
		}

		configurations = append(configurations, urlConfig)
	}

	b := NewBench(configurations...)

	checkTestResult(b, t)
}

func TestWrongURLs(t *testing.T) {
	for _, testCase := range failTestCases {
		urlConfig, err := WithURLSettings(testCase.testURL,
			testCase.testMethod,
			testCase.testData,
			testCase.headers,
			testCase.rawCookie,
			testCase.userPass)

		if err == nil || !strings.Contains(err.Error(), testCase.expectedErrorMsg) {
			t.Errorf(`Expected to get an error containing "%s" for url string "%s" but got: %v`, testCase.expectedErrorMsg, testCase.testURL, err)
		}

		if urlConfig != nil {
			t.Errorf("Expected to get a nil in resp for %s", testCase.testURL)
		}
	}
}

func checkTestResult(b *Bench, t *testing.T) {
	if len(b.URLs) != len(expectedURLs) {
		t.Fatalf("Wrong number of urls")
	}

	for index, url := range expectedURLs {
		addr := b.URLs[index].Addr
		method := b.URLs[index].Method
		data := b.URLs[index].Data
		headers := b.URLs[index].Headers
		rawCookie := b.URLs[index].RawCookie
		auth := b.URLs[index].Auth

		if url.expectedURL.Addr != addr {
			t.Errorf("Expected address %s but got %s on url index %d", url.expectedURL.Addr, addr, index)
		}

		if url.expectedURL.Method != method {
			t.Errorf("Expected method %s but got %s on url index %d", url.expectedURL.Method, method, index)
		}

		if (url.expectedURL.Data == nil || len(url.expectedURL.Data) == 0) && len(data) > 0 {
			t.Fatalf("Did not expect any data at url index %d but got %v", index, data)
		}

		if url.expectedURL.RawCookie != rawCookie {
			t.Errorf("Unexpected %q as RawCookie but got %q", url.expectedURL.RawCookie, rawCookie)
		}

		if url.expectedURL.Auth == nil && auth != nil {
			t.Errorf("Did not expect to get any Auth config nothing for index %d", index)
		}

		if url.expectedURL.Auth != nil && auth == nil {
			t.Errorf("Expected to get Auth config but got nothing for index %d", index)
		}

		if url.expectedURL.Auth != nil && auth != nil &&
			(url.expectedURL.Auth.Username != auth.Username ||
				url.expectedURL.Auth.Password != auth.Password) {

			t.Errorf("Expected to get %+v as auth but got %+v for index %d", url.expectedURL.Auth, auth, index)
		}

		checkMap(index, data, url.expectedURL.Data, t, "data")
		checkMap(index, headers, url.expectedURL.Headers, t, "headers")
	}
}

func checkMap(index int, got, expected map[string]string, t *testing.T, name string) {
	if expected != nil && len(expected) != 0 {
		if len(expected) != len(got) {
			t.Fatalf("Wrong number of %s at url index %d", name, index)
		}

		for k, expectedValue := range expected {
			if val, ok := got[k]; !ok || val != expectedValue {
				t.Errorf("Expected to have %s=%s in the %s", k, expectedValue, name)
			}
		}

		return
	}

	if got != nil && len(got) > 0 {
		t.Errorf("Did not expect to get anything for %s at %d but got %v", name, index, got)
	}
}
