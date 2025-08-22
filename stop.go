package alns

import (
	"context"
	"math/rand/v2"
	"time"
)

type StoppingCriterion interface {
	IsDone(rnd *rand.Rand, best, current State) bool
}

type MaxIterations struct {
	MaxIterations    int
	currentIteration int
}

var _ StoppingCriterion = &MaxIterations{}

func (mi *MaxIterations) IsDone(rnd *rand.Rand, best, current State) bool {
	mi.currentIteration++
	return mi.currentIteration > mi.MaxIterations
}

type MaxRuntime struct {
	MaxRuntime time.Duration
	started    time.Time
}

var _ StoppingCriterion = &MaxRuntime{}

func (mr *MaxRuntime) IsDone(rnd *rand.Rand, best, current State) bool {
	if mr.started.IsZero() {
		mr.started = time.Now()
		return false
	}
	return time.Since(mr.started) > mr.MaxRuntime
}

type NoImprovement struct {
	MaxIterations int
	counter       int
	isInitialized bool
	target        float64
}

var _ StoppingCriterion = &NoImprovement{}

func (ni *NoImprovement) IsDone(rnd *rand.Rand, best, current State) bool {
	if !ni.isInitialized || best.Objective() < ni.target {
		ni.isInitialized = true
		ni.target = best.Objective()
		ni.counter = 0
	} else {
		ni.counter++
	}
	return ni.counter >= ni.MaxIterations
}

type StoppingCriterions []StoppingCriterion

var _ StoppingCriterion = StoppingCriterions{}

func (sc StoppingCriterions) IsDone(rnd *rand.Rand, best, current State) bool {
	if len(sc) == 0 {
		panic("no criterias were specified")
	}
	for _, c := range sc {
		if c.IsDone(rnd, best, current) {
			return true
		}
	}
	return false
}

type Context struct {
	Context context.Context
}

var _ StoppingCriterion = &Context{}

func (c *Context) IsDone(rnd *rand.Rand, best, current State) bool {
	select {
	case <-c.Context.Done():
		return true
	default:
		return false
	}
}
