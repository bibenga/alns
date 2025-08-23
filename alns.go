package alns

import (
	"cmp"
	"math/rand/v2"
	"time"
)

type Comparator[O any] func(a, b O) int

type Listener[O any] func(outcome Outcome, cand State[O])

type ALNS[O any] struct {
	Rnd               *rand.Rand
	Compare           Comparator[O]
	CollectObjectives bool
	Listener          Listener[O]
	DestroyOperators  []Operator[O]
	RepairOperators   []Operator[O]
}

func NewOrdered[O cmp.Ordered]() ALNS[O] {
	return ALNS[O]{
		Rnd:               RuntimeRand,
		Compare:           cmp.Compare[O],
		CollectObjectives: true,
	}
}

func (a *ALNS[O]) AddDestroyOperator(ops ...Operator[O]) {
	a.DestroyOperators = append(a.DestroyOperators, ops...)
}

func (a *ALNS[O]) AddRepairOperator(ops ...Operator[O]) {
	a.RepairOperators = append(a.RepairOperators, ops...)
}

func (a *ALNS[O]) Iterate(
	initialSolution State[O],
	selector OperatorSelectionScheme[O],
	accept AcceptanceCriterion[O],
	stop StoppingCriterion[O],
) Result[O] {
	if len(a.DestroyOperators) == 0 || len(a.RepairOperators) == 0 {
		panic("Missing destroy or repair operators.")
	}

	curr := initialSolution
	best := initialSolution

	numIterations := 0
	if a.CollectObjectives {
		if maxIterations, ok := stop.(*MaxIterations[O]); ok {
			numIterations = maxIterations.MaxIterations + 1
		}
	}
	stats := newStatistics[O](numIterations, len(a.DestroyOperators), len(a.RepairOperators))

	started := time.Now()
	stats.IterationCount++
	if a.CollectObjectives {
		stats.collectObjective(0, initialSolution.Objective())
	}

	for !stop.IsDone(a.Rnd, best, curr) {
		dIdx, rIdx := selector.Select(a.Rnd, best, curr)
		destroyOp := a.DestroyOperators[dIdx]
		repairOp := a.RepairOperators[rIdx]

		destroyed := destroyOp(curr, a.Rnd)
		cand := repairOp(destroyed, a.Rnd)

		var outcome Outcome
		best, curr, outcome = a.evalCand(accept, best, curr, cand)

		selector.Update(cand, dIdx, rIdx, outcome)

		stats.IterationCount++
		if a.CollectObjectives {
			stats.collectObjective(time.Since(started), curr.Objective())
		}
		stats.collectOperators(dIdx, rIdx, outcome)
	}
	stats.TotalRuntime = time.Since(started)

	return Result[O]{BestState: best, Statistics: stats}
}

func (a *ALNS[O]) evalCand(accept AcceptanceCriterion[O], best, curr, cand State[O]) (State[O], State[O], Outcome) {
	outcome := a.determineOutcome(accept, best, curr, cand)

	if a.Listener != nil {
		a.Listener(outcome, cand)
	}

	switch outcome {
	case Best:
		return cand, cand, outcome
	case Reject:
		return best, curr, outcome
	default:
		return best, cand, outcome
	}
}

func (a *ALNS[O]) determineOutcome(accept AcceptanceCriterion[O], best, curr, cand State[O]) Outcome {
	outcome := Reject

	if accept.Accept(a.Rnd, best, curr, cand) {
		// accept candidate
		outcome = Accept

		if a.Compare(cand.Objective(), curr.Objective()) < 0 {
			outcome = Better
		}
	}

	if a.Compare(cand.Objective(), best.Objective()) < 0 {
		// candidate is new best
		outcome = Best
	}

	return outcome
}
