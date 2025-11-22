package alns

import (
	"math/rand/v2"
	"time"
)

type Listener func(outcome Outcome, cand State) error

type ALNS struct {
	Rnd               *rand.Rand
	CollectObjectives bool
	Listener          Listener
	DestroyOperators  []Operator
	RepairOperators   []Operator
}

// def iterate(initial_solution, select, accept, stop)
func (a *ALNS) Iterate(
	initSol State,
	selector OperatorSelectionScheme,
	acceptor AcceptanceCriterion,
	stop StoppingCriterion,
) (*Result, error) {
	if len(a.DestroyOperators) == 0 || len(a.RepairOperators) == 0 {
		panic("Missing destroy or repair operators.")
	}

	curr := initSol
	best := initSol

	numIterations := 0
	if a.CollectObjectives {
		if maxIterations, ok := stop.(*MaxIterations); ok {
			numIterations = maxIterations.MaxIterations + 1
		}
	}
	stats := newStatistics(numIterations, len(a.DestroyOperators), len(a.RepairOperators))

	started := time.Now()
	if a.CollectObjectives {
		stats.collectObjective(0, initSol.Objective())
	}

	for {
		if done, err := stop.IsDone(a.Rnd, best, curr); err != nil {
			return nil, err
		} else if done {
			break
		}
		dIdx, rIdx, err := selector.Select(a.Rnd, best, curr)
		if err != nil {
			return nil, err
		}
		destroyOp := a.DestroyOperators[dIdx]
		repairOp := a.RepairOperators[rIdx]

		destroyed, err := destroyOp(curr, a.Rnd)
		if err != nil {
			return nil, err
		}
		cand, err := repairOp(destroyed, a.Rnd)
		if err != nil {
			return nil, err
		}

		var outcome Outcome
		best, curr, outcome, err = a.evalCand(acceptor, best, curr, cand)
		if err != nil {
			return nil, err
		}

		err = selector.Update(cand, dIdx, rIdx, outcome)
		if err != nil {
			return nil, err
		}

		stats.IterationCount++
		if a.CollectObjectives {
			stats.collectObjective(time.Since(started), curr.Objective())
		}
		stats.collectOperators(dIdx, rIdx, outcome)
	}
	stats.TotalRuntime = time.Since(started)

	result := Result{
		BestState:  best,
		Statistics: stats,
	}
	return &result, nil
}

func (a *ALNS) evalCand(acceptor AcceptanceCriterion, best, curr, cand State) (State, State, Outcome, error) {
	outcome, err := a.determineOutcome(acceptor, best, curr, cand)
	if err != nil {
		return nil, nil, 0, err
	}

	if a.Listener != nil {
		if err := a.Listener(outcome, cand); err != nil {
			return nil, nil, 0, err
		}
	}

	switch outcome {
	case Best:
		return cand, cand, outcome, nil
	case Reject:
		return best, curr, outcome, nil
	default:
		return best, cand, outcome, nil
	}
}

func (a *ALNS) determineOutcome(acceptor AcceptanceCriterion, best, curr, cand State) (Outcome, error) {
	outcome := Reject

	if accepted, err := acceptor.Accept(a.Rnd, best, curr, cand); err != nil {
		return 0, err
	} else if accepted {
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

	return outcome, nil
}
