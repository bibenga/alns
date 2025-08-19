package alns

import (
	"fmt"
	"math"
	"math/rand/v2"
)

const absTol = 1e-08

func weightedRandomIndex(rnd *rand.Rand, weights []float64) int {
	if len(weights) == 1 {
		return 0
	}
	sum := sum(weights)
	value := rnd.Float64() * sum // adjusted value
	for i, weight := range weights {
		value -= weight
		if value <= 0 || isClose(value, 0, absTol) {
			return i
		}
	}
	panic(fmt.Sprintf("arithmetic error: sum=%f, value=%f, weights=%v", sum, value, weights))
}

func isClose(a, b, tol float64) bool {
	return math.Abs(a-b) <= tol
}

func sum(weights []float64) float64 {
	if len(weights) == 1 {
		return weights[0]
	}
	sum := 0.0
	for _, w := range weights {
		sum += w
	}
	return sum
}
