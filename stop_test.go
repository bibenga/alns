package alns

import (
	"cmp"
	"context"
	"math"
	"testing"
	"time"
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

func TestMaxRuntime(t *testing.T) {
	stop := MaxRuntime{MaxRuntime: 100 * time.Millisecond}

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

func TestNoImprovement(t *testing.T) {
	t.Run("Simple", func(t *testing.T) {
		stop := NoImprovement{
			Compare:       cmp.Compare[float64],
			MaxIterations: 10,
		}
		best := FakeState{objective: 1}
		curr := FakeState{objective: 1}
		i := 0
		for {
			if done, err := stop.IsDone(nil, best, curr); err != nil {
				t.Fatal(err)
			} else if done {
				break
			}
			i++
		}
		if i != 10 {
			t.Fatalf("10 iterations expected, actual %d iterations", i)
		}
		if stop.counter != 10 {
			t.Fatalf("number 10 expected, actual number %d", stop.counter)
		}
	})

	t.Run("SimulatedDecrease", func(t *testing.T) {
		stop := NoImprovement{MaxIterations: 10}
		best := FakeState{objective: 100}
		curr := FakeState{objective: 100}
		i := 0
		for {
			if done, err := stop.IsDone(nil, best, curr); err != nil {
				t.Fatal(err)
			} else if done {
				break
			}
			best.objective = max(math.Round(curr.objective-1), 0)
			curr.objective = max(math.Round(curr.objective-1), 0)
			i++
		}
		if i != 110 {
			t.Fatalf("110 iterations expected, actual %d iterations", i)
		}
		if stop.counter != 10 {
			t.Fatalf("number 110 expected, actual number %d", stop.counter)
		}
	})
}

func TestStoppingCriterions(t *testing.T) {
	stop := StoppingCriterions{
		&MaxIterations{MaxIterations: 10},
		&MaxIterations{MaxIterations: 20},
	}

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
}

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
