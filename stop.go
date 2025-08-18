package alns

import (
	"math/rand/v2"
	"time"
)

type StoppingCriterion interface {
	Call(rnd *rand.Rand, best, current State) bool
}

type MaxIterations struct {
	MaxIterations    int
	currentIteration int
}

var _ StoppingCriterion = &MaxIterations{}

func (s *MaxIterations) Call(rnd *rand.Rand, best, current State) bool {
	s.currentIteration++
	return s.currentIteration > s.MaxIterations
}

type NoImprovement struct {
	MaxIterations int
	counter       int
	isInitialized bool
	target        float64
}

var _ StoppingCriterion = &NoImprovement{}

func (n *NoImprovement) Call(rnd *rand.Rand, best, current State) bool {
	if !n.isInitialized || best.Objective() < n.target {
		n.target = best.Objective()
		n.counter = 0
	} else {
		n.counter++
	}
	return n.counter >= n.MaxIterations
}

type MaxRuntime struct {
	MaxRuntime time.Duration
	started    time.Time
}

var _ StoppingCriterion = &MaxRuntime{}

func (s *MaxRuntime) Call(rnd *rand.Rand, best, current State) bool {
	if s.started.IsZero() {
		s.started = time.Now()
		return false
	}
	return time.Since(s.started) > s.MaxRuntime
}
