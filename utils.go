package alns

import (
	"cmp"
	"fmt"
	"math"
	"math/rand/v2"
)

var UseIsCloseComparison = true
var RelativeTolerance float64 = 1e-12 // for numbers that can be of large order
var AbsoulteTolerance float64 = 1e-12 // for numbers near zero

func weightedRandomIndex(rnd *rand.Rand, weights []float64) int {
	if len(weights) == 0 {
		panic("invalid weights")
	}
	if len(weights) == 1 {
		return 0
	}
	sum := sum(weights)
	value := rnd.Float64() * sum // adjusted value
	for i, weight := range weights {
		value -= weight
		if Compare(value, 0) <= 0 {
			return i
		}
	}
	panic(fmt.Sprintf("arithmetic error: sum=%f, value=%f, weights=%v", sum, value, weights))
}

func Compare(a, b float64) int {
	if UseIsCloseComparison && IsClose(a, b) {
		return 0
	}
	return cmp.Compare(a, b)
}

func IsClose(a, b float64) bool {
	// return math.Abs(a-b) <= atol
	// see https://numpy.org/doc/stable/reference/generated/numpy.isclose.html#numpy.isclose
	// https://docs.python.org/3/library/math.html#math.isclose
	return a == b || math.Abs(a-b) <= (AbsoulteTolerance+RelativeTolerance*max(math.Abs(a), math.Abs(b)))
}

func sum(weights []float64) float64 {
	sum := 0.0
	for _, w := range weights {
		sum += w
	}
	return sum
}
