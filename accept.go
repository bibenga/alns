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

func NewHillClimbing() HillClimbing {
	return HillClimbing{}
}

func (a *HillClimbing) Accept(rnd *rand.Rand, best, current, candidate State) (bool, error) {
	return candidate.Objective() <= current.Objective(), nil
}
