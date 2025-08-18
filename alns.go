package alns

import (
	"math/rand/v2"
	"time"
)

type OnOutcome func(outcome Outcome, cand State)

type ALNS struct {
	Rnd              *rand.Rand
	destroyOperators []operator
	repairOperators  []operator
	OnOutcome        OnOutcome
}

func New(rnd *rand.Rand) *ALNS {
	return &ALNS{Rnd: rnd}
}

func NewDefault() *ALNS {
	return New(rand.New(&randomSource{}))
}

func NewPCGRandom(seed1, seed2 uint64) *ALNS {
	rnd := rand.New(rand.NewPCG(seed1, seed2))
	return New(rnd)
}

func (a *ALNS) AddDestroyOperator(op Operator, name string) {
	a.destroyOperators = append(a.destroyOperators, operator{call: op, name: name})
}

func (a *ALNS) AddRepairOperator(op Operator, name string) {
	a.repairOperators = append(a.repairOperators, operator{call: op, name: name})
}

func (a *ALNS) Iterate(
	initialSolution State,
	opSelect OperatorSelectionScheme,
	accept AcceptanceCriterion,
	stop StoppingCriterion,
) Result {
	if len(a.destroyOperators) == 0 || len(a.repairOperators) == 0 {
		panic("Missing destroy or repair operators.")
	}

	curr := initialSolution
	best := initialSolution
	initObj := initialSolution.Objective()

	// logger.debug(f"Initial solution has objective {init_obj:.2f}.")

	stats := Statistics{}
	stats.collectObjective(initObj)
	stats.collectRuntime(time.Now())

	for !stop.Call(a.Rnd, best, curr) {
		dIdx, rIdx := opSelect.Call(a.Rnd, best, curr)
		destroyOp := a.destroyOperators[dIdx]
		repairOp := a.repairOperators[rIdx]

		// logger.debug(f"Selected operators {d_name} and {r_name}.")

		destroyed := destroyOp.call(curr, a.Rnd)
		cand := repairOp.call(destroyed, a.Rnd)

		var outcome Outcome
		best, curr, outcome = a.evalCand(
			accept, best, curr, cand,
		)

		opSelect.Update(cand, dIdx, rIdx, outcome)

		stats.collectObjective(curr.Objective())
		stats.collectDestroyOperator(destroyOp.name, outcome)
		stats.collectRepairOperator(repairOp.name, outcome)
		stats.collectRuntime(time.Now())
	}

	// logger.info(f"Finished iterating in {stats.total_runtime:.2f}s.")

	return Result{BestState: best, Statistics: stats}
}

func (a *ALNS) evalCand(accept AcceptanceCriterion, best, curr, cand State) (State, State, Outcome) {
	outcome := a.determineOutcome(accept, best, curr, cand)

	if a.OnOutcome != nil {
		a.OnOutcome(outcome, cand)
	}

	if outcome == BEST {
		return cand, cand, outcome
	}
	if outcome == REJECT {
		return best, curr, outcome
	}
	return best, cand, outcome
}

func (a *ALNS) determineOutcome(accept AcceptanceCriterion, best, curr, cand State) Outcome {
	outcome := REJECT

	// slog.Info("determine outcome",
	// 	"best", best.Objective(),
	// 	"curr", curr.Objective(),
	// 	"cand", cand.Objective(),
	// )

	if accept.Call(a.Rnd, best, curr, cand) {
		// accept candidate
		outcome = ACCEPT

		if cand.Objective() < curr.Objective() {
			outcome = BETTER
		}
	}

	if cand.Objective() < best.Objective() {
		// candidate is new best
		// slog.Info("New best", "objective", cand.Objective())
		outcome = BEST
	}

	return outcome
}
