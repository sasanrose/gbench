package bench

import (
	"net/http"
	"testing"
)

func TestNewBenchWithConfig(t *testing.T) {
	b1 := NewBench(func(b *Bench) {
		b.Proxy = "testProxy"
	})

	if b1.Proxy != "testProxy" {
		t.Errorf("Expected the proxy to be %s", "testProxy")
	}

	configs := []func(b *Bench){
		func(b *Bench) {
			b.Concurrency = 1
		},
		func(b *Bench) {
			b.Proxy = "testProxy2"
		},
	}

	b2 := NewBench(configs...)

	if b2.Proxy != "testProxy2" {
		t.Errorf("Expected the proxy to be %s", "testProxy2")
	}

	if b2.Concurrency != 1 {
		t.Errorf("Expected the concurrency to be %d", 1)
	}
}

func TestDefaultSuccessStatusCodes(t *testing.T) {
	b := NewBench()
	expected := []int{
		http.StatusOK,
		http.StatusAccepted,
		http.StatusCreated,
	}

	if len(expected) != len(b.SuccessStatusCodes) {
		t.Errorf("Wrong default status codes: %v", b.SuccessStatusCodes)
	}

	for k, v := range expected {
		if v != b.SuccessStatusCodes[k] {
			t.Errorf("Wrong default status code. Expected %d but got %d", v, b.SuccessStatusCodes[k])
		}
	}
}
