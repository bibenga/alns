package alns

import (
	"cmp"
	"math/rand/v2"
)

type AcceptanceCriterion[O any] interface {
	Accept(rnd *rand.Rand, best, current, candidate State[O]) bool
}

type HillClimbing[O any] struct {
	Compare Comparator[O]
}

var _ AcceptanceCriterion[int] = &HillClimbing[int]{}

func NewOrderedHillClimbing[O cmp.Ordered]() HillClimbing[O] {
	return HillClimbing[O]{Compare: cmp.Compare[O]}
}

func (h *HillClimbing[O]) Accept(rnd *rand.Rand, best, current, candidate State[O]) bool {
	return h.Compare(candidate.Objective(), current.Objective()) <= 0
}
