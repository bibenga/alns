package alns

import (
	"math/rand/v2"
	"testing"
)

func TestRouletteWheel(t *testing.T) {
	t.Run("Validation", func(t *testing.T) {
		_, err := NewRouletteWheel([4]float64{3, 2, 1, 0.5}, 0.8, 2, 3, nil)
		if err != nil {
			t.Fatalf("is not valid: %s", err)
		}

		_, err = NewRouletteWheel([4]float64{-1, 2, 1, 0.5}, 0.8, 2, 3, nil)
		if err == nil || err.Error() != "negative scores are not understood" {
			t.Fatalf("is not valid: %s", err)
		}

		_, err = NewRouletteWheel([4]float64{4, 2, 1, 0.5}, -0.8, 2, 3, nil)
		if err == nil || err.Error() != "decay outside [0, 1] not understood" {
			t.Fatalf("is not valid: %s", err)
		}

		_, err = NewRouletteWheel([4]float64{4, 2, 1, 0.5}, 0.8, 2, 3, [][]bool{{true, true, true}, {true, true, true}})
		if err != nil {
			t.Fatalf("is not valid: %s", err)
		}

		_, err = NewRouletteWheel([4]float64{4, 2, 1, 0.5}, 0.8, 2, 3, [][]bool{})
		if err == nil || err.Error() != "coupling matrix of shape (0, 0), expected (2, 3)" {
			t.Fatalf("is not valid: %s", err)
		}

		_, err = NewRouletteWheel([4]float64{4, 2, 1, 0.5}, 0.8, 2, 3, [][]bool{{true, true}, {true}})
		if err == nil || err.Error() != "the number of columns in a row 1 does not match the expected 2" {
			t.Fatalf("is not valid: %s", err)
		}

		_, err = NewRouletteWheel([4]float64{4, 2, 1, 0.5}, 0.8, 2, 3, [][]bool{{true, true}, {true, true}, {true, true}})
		if err == nil || err.Error() != "coupling matrix of shape (3, 2), expected (2, 3)" {
			t.Fatalf("is not valid: %s", err)
		}
	})

	t.Run("Simple", func(t *testing.T) {
		r := rand.New(rand.NewPCG(1, 3))

		selector, _ := NewRouletteWheel([4]float64{3, 2, 1, 0.5}, 0.8, 3, 3, nil)

		best := FakeState{}
		current := FakeState{}
		candidate := FakeState{}

		dIdx, rIdx := selector.Call(r, best, current)
		selector.Update(candidate, dIdx, rIdx, Reject)
	})

	t.Run("Coupling", func(t *testing.T) {
		t.FailNow()
	})
}
