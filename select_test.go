package alns

import (
	"math/rand/v2"
	"testing"
)

func TestRouletteWheel(t *testing.T) {
	t.Run("Validation", func(t *testing.T) {
		_, err := NewRouletteWheel[float64]([4]float64{3, 2, 1, 0.5}, 0.8, 2, 3, nil)
		if err != nil {
			t.Fatalf("is not valid: %s", err)
		}

		_, err = NewRouletteWheel[float64]([4]float64{-1, 2, 1, 0.5}, 0.8, 2, 3, nil)
		if err == nil || err.Error() != "negative scores are not understood" {
			t.Fatalf("is not valid: %s", err)
		}

	})

	t.Run("CouplingValidation", func(t *testing.T) {
		_, err := NewRouletteWheel[float64]([4]float64{4, 2, 1, 0.5}, -0.8, 2, 3, nil)
		if err == nil || err.Error() != "decay outside [0, 1] not understood" {
			t.Fatalf("is not valid: %s", err)
		}

		_, err = NewRouletteWheel[float64]([4]float64{4, 2, 1, 0.5}, 0.8, 2, 3, [][]bool{{true, true, true}, {true, true, true}})
		if err != nil {
			t.Fatalf("is not valid: %s", err)
		}

		_, err = NewRouletteWheel[float64]([4]float64{4, 2, 1, 0.5}, 0.8, 2, 3, [][]bool{})
		if err == nil || err.Error() != "coupling matrix of shape (0, 0), expected (2, 3)" {
			t.Fatalf("is not valid: %s", err)
		}

		_, err = NewRouletteWheel[float64]([4]float64{4, 2, 1, 0.5}, 0.8, 2, 3, [][]bool{{true, true}, {true}})
		if err == nil || err.Error() != "the number of columns in a row 1 does not match the expected 2" {
			t.Fatalf("is not valid: %s", err)
		}

		_, err = NewRouletteWheel[float64]([4]float64{4, 2, 1, 0.5}, 0.8, 2, 3, [][]bool{{true, true}, {true, true}, {true, true}})
		if err == nil || err.Error() != "coupling matrix of shape (3, 2), expected (2, 3)" {
			t.Fatalf("is not valid: %s", err)
		}

		_, err = NewRouletteWheel[float64]([4]float64{4, 2, 1, 0.5}, 0.8, 2, 3, [][]bool{{true, false, false}, {false, false, false}})
		if err == nil || err.Error() != "destroy operator 1 has no coupled repair operators" {
			t.Fatalf("is not valid: %s", err)
		}
	})

	t.Run("Simple", func(t *testing.T) {
		r := rand.New(rand.NewPCG(1, 3))

		selector, _ := NewRouletteWheel[float64]([4]float64{3, 2, 1, 0.5}, 0.8, 3, 2, nil)

		best := FakeState{}
		current := FakeState{}
		candidate := FakeState{}

		dCounter := make([]int, 3)
		rCounter := make([]int, 2)
		total := 10000

		for range total {
			outcome := Outcome(r.IntN(4))
			dIdx, rIdx := selector.Select(r, best, current)
			selector.Update(candidate, dIdx, rIdx, outcome)
			dCounter[dIdx]++
			rCounter[rIdx]++
		}

		for counterNum, counter := range [][]int{dCounter, rCounter} {
			for i := range len(counter) {
				expected := 1 / float64(len(counter)) * float64(total)
				got := float64(counter[i])
				if (got-expected)/expected > 0.05 { // allow 5x% error
					t.Errorf("index (%d, %d): got %f, expected ~%f", counterNum, i, got, expected)
				}
			}
		}
	})

	t.Run("Coupling", func(t *testing.T) {
		r := rand.New(rand.NewPCG(1, 3))

		selector, err := NewRouletteWheel[float64]([4]float64{3, 2, 1, 0.5}, 0.8, 2, 3,
			[][]bool{{true, true, false}, {false, true, true}})
		if err != nil {
			t.Fatal(err)
		}

		best := FakeState{}
		current := FakeState{}
		candidate := FakeState{}

		dCounter := make([]int, selector.numDestroy)
		rCounter := make([]int, selector.numRepair)
		total := 10000

		for range total {
			outcome := Outcome(r.IntN(4))
			dIdx, rIdx := selector.Select(r, best, current)
			if !(0 <= dIdx && dIdx < selector.numDestroy) {
				t.Fatalf("destroy index %d is invalid", dIdx)
			}
			if !(0 <= rIdx && rIdx < selector.numRepair) {
				t.Fatalf("destroy index %d is invalid", dIdx)
			}
			selector.Update(candidate, dIdx, rIdx, outcome)
			dCounter[dIdx]++
			rCounter[rIdx]++
		}
		t.Log(dCounter)
		t.Log(rCounter)

		// fifty-fifty
		expectedDestoryPercent := []float64{0.5, 0.5}
		for i := range len(dCounter) {
			expected := expectedDestoryPercent[i] * float64(total)
			got := float64(dCounter[i])
			if (got-expected)/expected > 0.05 { // allow 5x% error
				t.Errorf("destroy index %d: got %f, expected ~%f", i, got, expected)
			}
		}

		expectedRepairPercent := []float64{0.25, 0.5, 0.25}
		for i := range len(rCounter) {
			expected := expectedRepairPercent[i] * float64(total)
			got := float64(rCounter[i])
			if (got-expected)/expected > 0.05 { // allow 5x% error
				t.Errorf("repair index %d: got %f, expected ~%f", i, got, expected)
			}
		}
	})
}
