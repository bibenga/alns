package alns

import (
	"math/rand/v2"
	"testing"
)

func TestWeightedRandomIndex(t *testing.T) {
	t.Run("Simple", func(t *testing.T) {
		r := rand.New(rand.NewPCG(1, 3))

		tests := []struct {
			weights []float64
			want    int
		}{
			{[]float64{1}, 0},       // one element
			{[]float64{0, 1}, 1},    // only the second weight is non-zero
			{[]float64{5, 0, 0}, 0}, // only the first
			{[]float64{0, 0, 3}, 2}, // only the third
		}

		for _, tt := range tests {
			got := weightedRandomIndex(r, tt.weights)
			if got != tt.want {
				t.Errorf("weights=%v: got %d, want %d", tt.weights, got, tt.want)
			}
		}
	})

	t.Run("Distribution", func(t *testing.T) {
		r := rand.New(rand.NewPCG(1, 3))

		weights := []float64{1, 2, 3, 4, 5, 6}
		total := 100000
		counts := make([]int, len(weights))

		for range total {
			idx := weightedRandomIndex(r, weights)
			counts[idx]++
		}

		// Let's check that the frequencies roughly match the weight fractions
		sum := 0.0
		for _, w := range weights {
			sum += w
		}

		for i, w := range weights {
			expected := (w / sum) * float64(total)
			got := float64(counts[i])
			if (got-expected)/expected > 0.01 { // allow 1% error
				t.Errorf("index %d: got %f, expected ~%f", i, got, expected)
			}
		}
	})
}
