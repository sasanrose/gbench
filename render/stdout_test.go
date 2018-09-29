package render

import (
	"bytes"
	"strings"
	"testing"

	"github.com/sasanrose/gbench/report"
)

// The strings should be in the order of appearance
var expectedStringsInOutput []string = []string{
	"Final benchmark result",
	"Total requests sent",
	"Total data received",
	"Total successful requests",
	"Total failed requests",
	"Total timedout requests",
	"Success rate",
	"Failure rate",
	"Timedout rate",
	"Total benchmark time",
	"Sum of all response times",
	"Shortest response time",
	"Longest response time",
	"Average response time",
	"Final result for http://testurl1.com",
	"Total data received",
	"Response with status code 200",
	"Response with status code 201",
	"Response with status code 500",
	"Failed requests",
	"Timedout requests",
	"Sum response times",
	"Shortest response time",
	"Longest response time",
	"Average response time",
	"Final result for http://testurl2.com",
	"Total data received",
	"Response with status code 200",
	"Response with status code 201",
	"Response with status code 500",
	"Failed requests",
	"Timedout requests",
	"Sum response times",
	"Shortest response time",
	"Longest response time",
	"Average response time",
	"Final result for http://testurl3.com",
	"Total data received",
	"Response with status code 500",
	"Response with status code 404",
	"Failed requests",
	"Timedout requests",
	"Sum response times",
	"Shortest response time",
	"Longest response time",
	"Average response time",
	"Result for concurrent requests batch 1",
	"Url",
	"http://testurl1.com",
	"http://testurl2.com",
	"http://testurl3.com",
	"Result for concurrent requests batch 2",
	"Url",
	"http://testurl1.com",
	"http://testurl2.com",
	"http://testurl3.com",
	"Result for concurrent requests batch 3",
	"Url",
	"http://testurl1.com",
	"http://testurl2.com",
	"http://testurl3.com",
}

func TestNew(t *testing.T) {
	r := NewStdoutRenderer()

	if _, ok := r.(Renderer); !ok {
		t.Error("Expected to get a var of Renderer interface type")
	}
}

func TestOutput(t *testing.T) {
	buf := bytes.NewBuffer([]byte{})

	r := &stdout{}
	r.output = buf

	result := &report.Result{}
	result.Init(2)

	addTestData(result)

	r.Render(result)

	index := 0
	output := buf.Bytes()

	for _, str := range expectedStringsInOutput {
		index := strings.Index(string(output[index:]), str)

		if index == -1 {
			t.Fatalf("Could not find %s in the output", str)
		}
	}
}
