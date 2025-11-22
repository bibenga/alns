package alns

import (
	"math/rand/v2"
)

type AcceptanceCriterion interface {
	Accept(rnd *rand.Rand, best, current, candidate State) (bool, error)
}

type HillClimbing struct {
	Compare CompareFunc
}

var _ AcceptanceCriterion = &HillClimbing{}

func NewHillClimbing(compare CompareFunc) HillClimbing {
	return HillClimbing{
		Compare: compare,
	}
}

func (a *HillClimbing) Accept(rnd *rand.Rand, best, current, candidate State) (bool, error) {
	return a.Compare(candidate.Objective(), current.Objective()) <= 0, nil
}
