package cmd

import (
	"net/http"
	"os"
	"testing"

	"github.com/sasanrose/gbench/bench"
)

func TestExecConfig(t *testing.T) {
	method = http.MethodPost
	data = []string{"key1=val1&key2=val2", "key3=val3"}

	configurations, _ := getExecConfig("http://url")

	b := bench.NewBench(configurations...)

	if len(b.URLs) != 1 {
		t.Fatal("Wrong number of URLs")
	}

	if b.URLs[0].Addr != "http://url" {
		t.Errorf("Expected 'url' as Addr but got %s", b.URLs[0].Addr)
	}

	if b.URLs[0].Method != http.MethodPost {
		t.Errorf("Expected 'post' as Method but got %s", b.URLs[0].Method)
	}

	if len(b.URLs[0].Data) != 3 {
		t.Fatal("Wrong number of data")
	}

	expected := map[string]string{
		"key1": "val1",
		"key2": "val2",
		"key3": "val3",
	}

	for k, v := range expected {
		if v != b.URLs[0].Data[k] {
			t.Errorf("Expected to get %q at %q but got %q", v, k, b.URLs[0].Data[k])
		}
	}
}

func TestUrlError(t *testing.T) {
	expected := "Error with url: Wrong key value format for request data: WrongData"
	method = http.MethodPost
	data = []string{"WrongData"}

	_, err := getExecConfig("http://url")

	if err == nil || err.Error() != expected {
		t.Errorf("Expected to get %q but got %v", expected, err)
	}
}

func TestNoUrl(t *testing.T) {
	if os.Getenv("CRASH_TEST") == "1" {
		runExec(execCmd, []string{})
		return
	}

	testExit(
		t,
		"TestNoUrl",
		"",
	)
}
