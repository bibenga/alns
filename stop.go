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

func (mi *MaxIterations) IsDone(rnd *rand.Rand, best, current State) (bool, error) {
	mi.currentIteration++
	return mi.currentIteration > mi.MaxIterations, nil
}

type MaxRuntime struct {
	MaxRuntime time.Duration
	started    time.Time
}

var _ StoppingCriterion = &MaxRuntime{}

func (mr *MaxRuntime) IsDone(rnd *rand.Rand, best, current State) (bool, error) {
	if mr.started.IsZero() {
		mr.started = time.Now()
		return false, nil
	}
	return time.Since(mr.started) > mr.MaxRuntime, nil
}

type NoImprovement struct {
	MaxIterations int
	counter       int
	isInitialized bool
	target        float64
}

var _ StoppingCriterion = &NoImprovement{}

func (ni *NoImprovement) IsDone(rnd *rand.Rand, best, current State) (bool, error) {
	if !ni.isInitialized || Compare(best.Objective(), ni.target) < 0 {
		ni.isInitialized = true
		ni.target = best.Objective()
		ni.counter = 0
	} else {
		ni.counter++
	}
	return ni.counter >= ni.MaxIterations, nil
}

type StoppingCriterions []StoppingCriterion

var _ StoppingCriterion = StoppingCriterions{}

func (sc StoppingCriterions) IsDone(rnd *rand.Rand, best, current State) (bool, error) {
	if len(sc) == 0 {
		panic("no criterias were specified")
	}
	for _, c := range sc {
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

func (c *Context) IsDone(rnd *rand.Rand, best, current State) (bool, error) {
	select {
	case <-c.Context.Done():
		return true, nil
	default:
		return false, nil
	}
}
