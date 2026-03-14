package maxiterations

import (
	"testing"
)

func TestMaxIterations(t *testing.T) {
	stop := MaxIterations{MaxIterations: 10}

	i := 0
	for {
		if done, err := stop.IsDone(nil, nil, nil); err != nil {
			t.Fatal(err)
		} else if done {
			break
		}
		i++
	}

	if i != 10 {
		t.Fatalf("10 iterations expected, actual %d iterations", i)
	}
	if stop.currentIteration != 11 {
		t.Fatalf("number 11 expected, actual number %d", stop.currentIteration)
	}
}
