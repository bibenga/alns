package testutil

import "math/rand/v2"

func DefaultRandomGenerator() *rand.Rand {
	return rand.New(rand.NewPCG(12, 34))
}
