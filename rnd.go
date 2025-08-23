package alns

import "math/rand/v2"

type randomSource struct{}

func (r *randomSource) Uint64() uint64 {
	return rand.Uint64()
}

var _ rand.Source = &randomSource{}

var RuntimeRand *rand.Rand = rand.New(&randomSource{})
