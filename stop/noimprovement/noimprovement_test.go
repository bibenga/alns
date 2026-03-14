package noimprovement

import (
	"math"
	"testing"

	"github.com/bibenga/alns/internal/testutil"
)

func TestNoImprovement(t *testing.T) {
	t.Run("Simple", func(t *testing.T) {
		stop := NoImprovement{
			MaxIterations: 10,
		}
		best := testutil.FakeState{Value: 1}
		curr := testutil.FakeState{Value: 1}
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
		stop := NoImprovement{
			MaxIterations: 10,
		}
		best := testutil.FakeState{Value: 100}
		curr := testutil.FakeState{Value: 100}
		i := 0
		for {
			if done, err := stop.IsDone(nil, best, curr); err != nil {
				t.Fatal(err)
			} else if done {
				break
			}
			best.Value = max(math.Round(curr.Value-1), 0)
			curr.Value = max(math.Round(curr.Value-1), 0)
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
