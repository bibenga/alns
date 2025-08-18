package alns

import "time"

type Statistics struct {
	Objectives            []float64
	Runtimes              []time.Time
	DestroyOperatorCounts map[string][]int // operator name -> outcome -> count
	RepairOperatorCounts  map[string][]int // operator name -> outcome -> count
}

func (s *Statistics) collectObjective(objective float64) {
	s.Objectives = append(s.Objectives, objective)
}

func (s *Statistics) collectRuntime(t time.Time) {
	s.Runtimes = append(s.Runtimes, t)
}

func (s *Statistics) collectDestroyOperator(name string, outcome Outcome) {
	if s.DestroyOperatorCounts == nil {
		s.DestroyOperatorCounts = make(map[string][]int)
	}
	opStats, ok := s.DestroyOperatorCounts[name]
	if !ok {
		opStats = make([]int, 4)
		s.DestroyOperatorCounts[name] = opStats
	}
	opStats[outcome] += 1
}

func (s *Statistics) collectRepairOperator(name string, outcome Outcome) {
	if s.RepairOperatorCounts == nil {
		s.RepairOperatorCounts = make(map[string][]int)
	}
	opStats, ok := s.RepairOperatorCounts[name]
	if !ok {
		opStats = make([]int, 4)
		s.RepairOperatorCounts[name] = opStats
	}
	opStats[outcome] += 1
}

func (s *Statistics) IterationCount() int {
	return len(s.Objectives)
}

func (s *Statistics) TotalRuntime() time.Duration {
	return s.Runtimes[len(s.Runtimes)-1].Sub(s.Runtimes[0])
}
