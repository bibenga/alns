package hillclimbing

import (
	"math/rand/v2"

	"github.com/bibenga/alns"
)

type HillClimbing struct {
}

var _ alns.AcceptanceCriterion = &HillClimbing{}

func NewHillClimbing() HillClimbing {
	return HillClimbing{}
}

func (a *HillClimbing) Accept(rnd *rand.Rand, best, current, candidate alns.State) (bool, error) {
	return candidate.Objective() <= current.Objective(), nil
}
