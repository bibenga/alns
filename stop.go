package alns

import (
	"context"
	"math/rand/v2"
	"time"
)

type StoppingCriterion interface {
	IsDone(rnd *rand.Rand, best, current State) (bool, error)
}

type MaxIterations struct {
	MaxIterations    int
	currentIteration int
}

var _ StoppingCriterion = &MaxIterations{}

func NewMaxIterations(maxIterations int) MaxIterations {
	return MaxIterations{
		MaxIterations: maxIterations,
	}
}

func (s *MaxIterations) IsDone(rnd *rand.Rand, best, current State) (bool, error) {
	s.currentIteration++
	return s.currentIteration > s.MaxIterations, nil
}

type MaxRuntime struct {
	MaxRuntime time.Duration
	started    time.Time
}

var _ StoppingCriterion = &MaxRuntime{}

func NewMaxRuntime(maxRuntime time.Duration) MaxRuntime {
	return MaxRuntime{
		MaxRuntime: maxRuntime,
	}
}

func (s *MaxRuntime) IsDone(rnd *rand.Rand, best, current State) (bool, error) {
	if s.started.IsZero() {
		s.started = time.Now()
		return false, nil
	}
	return time.Since(s.started) > s.MaxRuntime, nil
}

type NoImprovement struct {
	MaxIterations int
	counter       int
	isInitialized bool
	target        float64
}

var _ StoppingCriterion = &NoImprovement{}

func NewNoImprovement(maxIterations int) NoImprovement {
	return NoImprovement{
		MaxIterations: maxIterations,
	}
}

func (s *NoImprovement) IsDone(rnd *rand.Rand, best, current State) (bool, error) {
	if !s.isInitialized || best.Objective() < s.target {
		s.isInitialized = true
		s.target = best.Objective()
		s.counter = 0
	} else {
		s.counter++
	}
	return s.counter >= s.MaxIterations, nil
}

type StoppingCriterions []StoppingCriterion

var _ StoppingCriterion = StoppingCriterions{}

func NewStoppingCriterions(criterions ...StoppingCriterion) StoppingCriterions {
	return criterions
}

func (s StoppingCriterions) IsDone(rnd *rand.Rand, best, current State) (bool, error) {
	if len(s) == 0 {
		panic("no criterias were specified")
	}
	for _, c := range s {
		if done, err := c.IsDone(rnd, best, current); err != nil {
			return true, err
		} else if done {
			return true, nil
		}
	}
	return false, nil
}

type Context struct {
	Context context.Context
}

var _ StoppingCriterion = &Context{}

func NewContext(context context.Context) Context {
	return Context{
		Context: context,
	}
}

func (s *Context) IsDone(rnd *rand.Rand, best, current State) (bool, error) {
	select {
	case <-s.Context.Done():
		return true, nil
	default:
		return false, nil
	}
}
