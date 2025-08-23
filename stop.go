package alns

import (
	"cmp"
	"context"
	"math/rand/v2"
	"time"
)

type StoppingCriterion[O any] interface {
	IsDone(rnd *rand.Rand, best, current State[O]) bool
}

type MaxIterations[O any] struct {
	MaxIterations    int
	currentIteration int
}

var _ StoppingCriterion[int] = &MaxIterations[int]{}

func (mi *MaxIterations[O]) IsDone(rnd *rand.Rand, best, current State[O]) bool {
	mi.currentIteration++
	return mi.currentIteration > mi.MaxIterations
}

type MaxRuntime[O any] struct {
	MaxRuntime time.Duration
	started    time.Time
}

var _ StoppingCriterion[int] = &MaxRuntime[int]{}

func (mr *MaxRuntime[O]) IsDone(rnd *rand.Rand, best, current State[O]) bool {
	if mr.started.IsZero() {
		mr.started = time.Now()
		return false
	}
	return time.Since(mr.started) > mr.MaxRuntime
}

type NoImprovement[O any] struct {
	Compare       Comparator[O]
	MaxIterations int
	counter       int
	isInitialized bool
	target        O
}

var _ StoppingCriterion[int] = &NoImprovement[int]{}

func NewNoImprovement[O cmp.Ordered](maxIterations int) NoImprovement[O] {
	return NoImprovement[O]{
		Compare:       cmp.Compare[O],
		MaxIterations: maxIterations,
	}
}

func (ni *NoImprovement[O]) IsDone(rnd *rand.Rand, best, current State[O]) bool {
	if !ni.isInitialized || ni.Compare(best.Objective(), ni.target) < 0 {
		ni.isInitialized = true
		ni.target = best.Objective()
		ni.counter = 0
	} else {
		ni.counter++
	}
	return ni.counter >= ni.MaxIterations
}

type StoppingCriterions[O any] []StoppingCriterion[O]

var _ StoppingCriterion[int] = StoppingCriterions[int]{}

func (sc StoppingCriterions[O]) IsDone(rnd *rand.Rand, best, current State[O]) bool {
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

type Context[O any] struct {
	Context context.Context
}

var _ StoppingCriterion[int] = &Context[int]{}

func (c *Context[O]) IsDone(rnd *rand.Rand, best, current State[O]) bool {
	select {
	case <-c.Context.Done():
		return true
	default:
		return false
	}
}
