package math

import (
	"math/rand/v2"
)

func WeightedRandomIndex(rnd *rand.Rand, weights []float64) int {
	if len(weights) == 0 {
		panic("invalid weights")
	}
	if len(weights) == 1 {
		return 0
	}
	sum := Sum(weights)
	value := rnd.Float64() * sum // adjusted value
	for i, weight := range weights {
		value -= weight
		if value <= 0 {
			return i
		}
	}
	// we will only be here when errors accumulate
	return len(weights) - 1
}

func Sum(weights []float64) float64 {
	sum := 0.0
	for _, w := range weights {
		sum += w
	}
	return sum
}
