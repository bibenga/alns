package alns

import (
	"fmt"
	"time"
)

type Statistics struct {
	IterationCount        uint
	TotalRuntime          time.Duration
	Objectives            []IterationObjective
	DestroyOperatorCounts []OperatorStatistics
	RepairOperatorCounts  []OperatorStatistics
}

func newStatistics(numIterations int, numDestroy, numRepair int) Statistics {
	var objectives []IterationObjective
	if numIterations > 0 {
		objectives = make([]IterationObjective, 0, numIterations)
	}
	return Statistics{
		Objectives:            objectives,
		DestroyOperatorCounts: make([]OperatorStatistics, numDestroy),
		RepairOperatorCounts:  make([]OperatorStatistics, numRepair),
	}
}

func (s *Statistics) collectObjective(t time.Duration, objective float64) {
	s.Objectives = append(s.Objectives, IterationObjective{Elapsed: t, Objective: objective})
}

func (s *Statistics) collectOperators(dIdx, rIdx int, outcome Outcome) {
	s.DestroyOperatorCounts[dIdx][outcome]++
	s.RepairOperatorCounts[rIdx][outcome]++
}

type IterationObjective struct {
	Elapsed   time.Duration
	Objective float64
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
