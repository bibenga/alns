package alns

import (
	"testing"
	"time"
)

func TestMaxIterations(t *testing.T) {
	stop := MaxIterations{MaxIterations: 10}

	i := 0
	for !stop.Call(nil, nil, nil) {
		i++
	}

	if i != 10 {
		t.Fatalf("10 iterations expected, actual %d iterations", i)
	}
	if stop.currentIteration != 11 {
		t.Fatalf("number 11 expected, actual number %d", stop.currentIteration)
	}
}

func TestMaxRuntime(t *testing.T) {
	stop := MaxRuntime{MaxRuntime: 100 * time.Millisecond}

	started := time.Now()
	for !stop.Call(nil, nil, nil) {
		time.Sleep(1 * time.Millisecond)
	}
	elapsed := time.Since(started).Milliseconds()
	if !(100 <= elapsed && elapsed <= 105) {
		t.Fatalf("expected duration 100ms, actual %d", elapsed)
	}
}
