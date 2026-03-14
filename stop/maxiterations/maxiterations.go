package maxiterations

import (
	"math/rand/v2"

	"github.com/bibenga/alns"
)

type MaxIterations struct {
	MaxIterations    int
	currentIteration int
}

var _ alns.StoppingCriterion = &MaxIterations{}

func NewMaxIterations(maxIterations int) MaxIterations {
	return MaxIterations{
		MaxIterations: maxIterations,
	}
}

func (s *MaxIterations) IsDone(rnd *rand.Rand, best, current alns.State) (bool, error) {
	s.currentIteration++
	return s.currentIteration > s.MaxIterations, nil
}
