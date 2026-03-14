package context

import (
	"context"
	"testing"
	"time"
)

func TestContext(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	stop := Context{Context: ctx}

	started := time.Now()
	for {
		if done, err := stop.IsDone(nil, nil, nil); err != nil {
			t.Fatal(err)
		} else if done {
			break
		}
		time.Sleep(1 * time.Millisecond)
	}
	elapsed := time.Since(started).Milliseconds()
	if !(100 <= elapsed && elapsed <= 105) {
		t.Fatalf("expected duration 100ms, actual %d", elapsed)
	}
}
