package alns

import "math/rand/v2"

type AcceptanceCriterion interface {
	Call(rnd *rand.Rand, best, current, candidate State) bool
}

type HillClimbing struct{}

var _ AcceptanceCriterion = &HillClimbing{}

func (h *HillClimbing) Call(rnd *rand.Rand, best, current, candidate State) bool {
	return candidate.Objective() <= current.Objective()
}
