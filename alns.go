package alns

import (
	"math/rand/v2"
	"time"
)

type Listener func(outcome Outcome, cand State)

type ALNS struct {
	Rnd               *rand.Rand
	Listener          Listener
	CollectObjectives bool
	DestroyOperators  []Operator
	RepairOperators   []Operator
}

func New(rnd *rand.Rand) ALNS {
	return ALNS{Rnd: rnd}
}

func NewDefault() ALNS {
	return New(rand.New(&randomSource{}))
}

func NewWithPCGRandom(seed1, seed2 uint64) ALNS {
	rnd := rand.New(rand.NewPCG(seed1, seed2))
	return New(rnd)
}

func (a *ALNS) AddDestroyOperator(op Operator) {
	a.DestroyOperators = append(a.DestroyOperators, op)
}

func (a *ALNS) AddRepairOperator(op Operator) {
	a.RepairOperators = append(a.RepairOperators, op)
}

func (a *ALNS) Iterate(
	initialSolution State,
	opSelect OperatorSelectionScheme,
	accept AcceptanceCriterion,
	stop StoppingCriterion,
) Result {
	if len(a.DestroyOperators) == 0 || len(a.RepairOperators) == 0 {
		panic("Missing destroy or repair operators.")
	}

	curr := initialSolution
	best := initialSolution

	numIterations := 0
	if a.CollectObjectives {
		if maxIterations, ok := stop.(*MaxIterations); ok {
			numIterations = maxIterations.MaxIterations + 1
		}
	}
	stats := newStatistics(numIterations, len(a.DestroyOperators), len(a.RepairOperators))

	started := time.Now()
	if a.CollectObjectives {
		stats.collectObjective(0, initialSolution.Objective())
	}

	for !stop.Call(a.Rnd, best, curr) {
		dIdx, rIdx := opSelect.Call(a.Rnd, best, curr)
		destroyOp := a.DestroyOperators[dIdx]
		repairOp := a.RepairOperators[rIdx]

		destroyed := destroyOp(curr, a.Rnd)
		cand := repairOp(destroyed, a.Rnd)

		var outcome Outcome
		best, curr, outcome = a.evalCand(accept, best, curr, cand)

		opSelect.Update(cand, dIdx, rIdx, outcome)

		if a.CollectObjectives {
			stats.collectObjective(time.Since(started), curr.Objective())
		}
		stats.collectOperators(dIdx, rIdx, outcome)
	}
	stats.TotalRuntime = time.Since(started)

	return Result{BestState: best, Statistics: stats}
}

func (a *ALNS) evalCand(accept AcceptanceCriterion, best, curr, cand State) (State, State, Outcome) {
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

func (a *ALNS) determineOutcome(accept AcceptanceCriterion, best, curr, cand State) Outcome {
	outcome := Reject

	if accept.Call(a.Rnd, best, curr, cand) {
		// accept candidate
		outcome = Accept

		if cand.Objective() < curr.Objective() {
			outcome = Better
		}
	}

	if cand.Objective() < best.Objective() {
		// candidate is new best
		outcome = Best
	}

	return outcome
}
