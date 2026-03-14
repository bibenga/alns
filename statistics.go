package alns

import (
	"fmt"
	"time"
)

type Statistics struct {
	IterationCount        uint              // the number of iterations
	TotalRuntime          time.Duration     // the total runtime
	Objectives            []ObjectiveRecord // previous objective values, tracking progress
	DestroyOperatorCounts []OperatorCounts  // the destroy operator counts
	RepairOperatorCounts  []OperatorCounts  // the repair operator counts
}

func newStatistics(numDestroy, numRepair int) Statistics {
	return Statistics{
		Objectives:            make([]ObjectiveRecord, 0, 128),
		DestroyOperatorCounts: make([]OperatorCounts, numDestroy),
		RepairOperatorCounts:  make([]OperatorCounts, numRepair),
	}
}

func (s *Statistics) collectObjective(elapsedTime time.Duration, objective float64) {
	s.Objectives = append(s.Objectives, ObjectiveRecord{ElapsedTime: elapsedTime, Objective: objective})
}

func (s *Statistics) collectOperators(dIdx, rIdx int, outcome Outcome) {
	s.DestroyOperatorCounts[dIdx][outcome]++
	s.RepairOperatorCounts[rIdx][outcome]++
}

type ObjectiveRecord struct {
	ElapsedTime time.Duration
	Objective   float64
}

type OperatorCounts [4]uint // see Outcome

func (o OperatorCounts) String() string {
	return fmt.Sprintf(
		"{%s:%d %s:%d %s:%d %s:%d}",
		Best, o[Best],
		Better, o[Better],
		Accept, o[Accept],
		Reject, o[Reject],
	)
}
