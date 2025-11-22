package alns

import (
	"math/rand/v2"
)

type AcceptanceCriterion interface {
	Accept(rnd *rand.Rand, best, current, candidate State) (bool, error)
}

type HillClimbing struct {
}

var _ AcceptanceCriterion = &HillClimbing{}

func (h *HillClimbing) Accept(rnd *rand.Rand, best, current, candidate State) (bool, error) {
	return Compare(candidate.Objective(), current.Objective()) <= 0, nil
}
