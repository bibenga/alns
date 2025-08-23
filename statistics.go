package alns

import (
	"fmt"
	"time"
)

type Statistics[O any] struct {
	IterationCount        int
	TotalRuntime          time.Duration
	Objectives            []IterationObjective[O]
	DestroyOperatorCounts []OperatorStatistics
	RepairOperatorCounts  []OperatorStatistics
}

func newStatistics[O any](numIterations int, numDestroy, numRepair int) Statistics[O] {
	var objectives []IterationObjective[O]
	if numIterations > 0 {
		objectives = make([]IterationObjective[O], 0, numIterations)
	}
	return Statistics[O]{
		Objectives:            objectives,
		DestroyOperatorCounts: make([]OperatorStatistics, numDestroy),
		RepairOperatorCounts:  make([]OperatorStatistics, numRepair),
	}
}

func (s *Statistics[O]) collectObjective(t time.Duration, objective O) {
	s.Objectives = append(s.Objectives, IterationObjective[O]{Elapsed: t, Objective: objective})
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
