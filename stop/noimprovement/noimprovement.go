package noimprovement

import (
	"math/rand/v2"

	"github.com/bibenga/alns"
)

type NoImprovement struct {
	MaxIterations int
	counter       int
	isInitialized bool
	target        float64
}

var _ alns.StoppingCriterion = &NoImprovement{}

func NewNoImprovement(maxIterations int) NoImprovement {
	return NoImprovement{
		MaxIterations: maxIterations,
	}
}

func (s *NoImprovement) IsDone(rnd *rand.Rand, best, current alns.State) (bool, error) {
	if !s.isInitialized || best.Objective() < s.target {
		s.isInitialized = true
		s.target = best.Objective()
		s.counter = 0
	} else {
		s.counter++
	}
	return s.counter >= s.MaxIterations, nil
}
