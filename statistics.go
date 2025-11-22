package alns

import (
	"fmt"
	"time"
)

type Statistics struct {
	IterationCount        int                  // the number of iterations
	TotalRuntime          time.Duration        // the total runtime
	Runtimes              []time.Duration      // run times
	Objectives            []float64            // previous objective values, tracking progress
	DestroyOperatorCounts []OperatorStatistics // the destroy operator counts
	RepairOperatorCounts  []OperatorStatistics // the repair operator counts
}

func newStatistics(numIterations int, numDestroy, numRepair int) Statistics {
	var runtimes []time.Duration
	var objectives []float64
	if numIterations > 0 {
		runtimes = make([]time.Duration, 0, numIterations+1)
		objectives = make([]float64, 0, numIterations+1)
	}
	return Statistics{
		Runtimes:              runtimes,
		Objectives:            objectives,
		DestroyOperatorCounts: make([]OperatorStatistics, numDestroy),
		RepairOperatorCounts:  make([]OperatorStatistics, numRepair),
	}
}

func (s *Statistics) collectObjective(t time.Duration, objective float64) {
	s.Runtimes = append(s.Runtimes, t)
	s.Objectives = append(s.Objectives, objective)
}

func (s *Statistics) collectOperators(dIdx, rIdx int, outcome Outcome) {
	s.DestroyOperatorCounts[dIdx][outcome]++
	s.RepairOperatorCounts[rIdx][outcome]++
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
