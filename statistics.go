package alns

import (
	"fmt"
	"time"
)

type Statistics struct {
	Objectives            []IterationObjective // todo: make optional
	DestroyOperatorCounts []OperatorStatistics
	RepairOperatorCounts  []OperatorStatistics
}

func newStatistics(numIterations int, numDestroy, numRepair int) Statistics {
	var iterations []IterationObjective
	if numIterations > 0 {
		iterations = make([]IterationObjective, 0, numIterations)
	}
	return Statistics{
		Objectives:            iterations,
		DestroyOperatorCounts: make([]OperatorStatistics, numDestroy),
		RepairOperatorCounts:  make([]OperatorStatistics, numRepair),
	}
}

func (s *Statistics) collectObjective(t time.Duration, objective float64) {
	s.Objectives = append(s.Objectives, IterationObjective{Runtime: t, Objective: objective})
}

func (s *Statistics) collectOperators(dIdx, rIdx int, outcome Outcome) {
	s.DestroyOperatorCounts[dIdx][outcome]++
	s.RepairOperatorCounts[rIdx][outcome]++
}

func (s *Statistics) IterationCount() int {
	return len(s.Objectives)
}

func (s *Statistics) TotalRuntime() time.Duration {
	return s.Objectives[len(s.Objectives)-1].Runtime
}

type IterationObjective struct {
	Objective float64
	Runtime   time.Duration
}

type OperatorStatistics [4]int

func (o OperatorStatistics) String() string {
	return fmt.Sprintf(
		"{%s:%d %s:%d %s:%d %s:%d}",
		Best, o[Best],
		Better, o[Better],
		Accept, o[Accept],
		Reject, o[Reject],
	)
}
