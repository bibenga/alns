package alns

import (
	"fmt"
	"time"
)

type Statistics[O any] struct {
	IterationCount        int                  // the number of iterations
	TotalRuntime          time.Duration        // the total runtime
	Runtimes              []time.Duration      // run times
	Objectives            []O                  // previous objective values, tracking progress
	DestroyOperatorCounts []OperatorStatistics // the destroy operator counts
	RepairOperatorCounts  []OperatorStatistics // the repair operator counts
}

func newStatistics[O any](numIterations int, numDestroy, numRepair int) Statistics[O] {
	var runtimes []time.Duration
	var objectives []O
	if numIterations > 0 {
		runtimes = make([]time.Duration, 0, numIterations)
		objectives = make([]O, 0, numIterations)
	}
	return Statistics[O]{
		Runtimes:              runtimes,
		Objectives:            objectives,
		DestroyOperatorCounts: make([]OperatorStatistics, numDestroy),
		RepairOperatorCounts:  make([]OperatorStatistics, numRepair),
	}
}

func (s *Statistics[O]) collectObjective(t time.Duration, objective O) {
	s.Runtimes = append(s.Runtimes, t)
	s.Objectives = append(s.Objectives, objective)
}

func (s *Statistics[O]) collectOperators(dIdx, rIdx int, outcome Outcome) {
	s.DestroyOperatorCounts[dIdx][outcome]++
	s.RepairOperatorCounts[rIdx][outcome]++
}

type IterationObjective[O any] struct {
	Elapsed   time.Duration
	Objective O
}

type OperatorStatistics [4]int // see Outcome

func (o OperatorStatistics) String() string {
	return fmt.Sprintf(
		"{%s:%d %s:%d %s:%d %s:%d}",
		Best, o[Best],
		Better, o[Better],
		Accept, o[Accept],
		Reject, o[Reject],
	)
}
