package alns

import (
	"math/rand/v2"
	"testing"
)

type FakeState struct {
	objective float64
}

func (s FakeState) Clone() *FakeState {
	return &FakeState{}
}

func (s FakeState) Objective() float64 {
	return s.objective
}

func TestAlns(t *testing.T) {
	a := NewWithPCGRandom(1, 2)
	a.CollectObjectives = true

	destroyCalled := 0
	a.AddDestroyOperator(func(state State, rnd *rand.Rand) State {
		destroyCalled++
		current := state.(*FakeState)
		destroyed := current.Clone()
		return destroyed
	})
	repairCalled := 0
	a.AddRepairOperator(func(state State, rnd *rand.Rand) State {
		repairCalled++
		current := state.(*FakeState)
		current.objective = rand.Float64()
		return current
	})

	initialSolution := FakeState{objective: 1}
	opSelect := NewRouletteWheel([4]float64{3, 2, 1, 0.5}, 0.8, 1, 1, nil)
	accept := HillClimbing{}
	stop := MaxIterations{MaxIterations: 10}
	res := a.Iterate(&initialSolution, &opSelect, &accept, &stop)

	if destroyCalled != 10 {
		t.Fatalf("10 destroy calls expected, actual %d calls", destroyCalled)
	}
	if repairCalled != 10 {
		t.Fatalf("10 repair calls expected, actual %d calls", repairCalled)
	}
	if stop.currentIteration != 11 {
		t.Fatalf("10 iterations expected, actual %d calls", stop.currentIteration)
	}
	if len(res.Statistics.Objectives) != 11 {
		t.Fatalf("11 objectives expected, actual %d objectives", len(res.Statistics.Objectives))
	}
	repairOperatorCounts := OperatorStatistics{3, 0, 0, 7}
	if res.Statistics.RepairOperatorCounts[0] != repairOperatorCounts {
		t.Fatalf("expected repair opeator statistics %v, actual %v",
			repairOperatorCounts, res.Statistics.RepairOperatorCounts[0])
	}
	rejectOperatorCounts := OperatorStatistics{3, 0, 0, 7}
	if res.Statistics.DestroyOperatorCounts[0] != rejectOperatorCounts {
		t.Fatalf("expected destory opeator statistics %v, actual %v",
			rejectOperatorCounts, res.Statistics.DestroyOperatorCounts[0])
	}
}

func TestAlnsCollectObjectives(t *testing.T) {
	solve := func(collectObjectives bool) Result {
		a := NewDefault()
		a.CollectObjectives = collectObjectives
		a.AddDestroyOperator(func(state State, rnd *rand.Rand) State { return state })
		a.AddRepairOperator(func(state State, rnd *rand.Rand) State { return state })
		initialSolution := FakeState{objective: 1}
		opSelect := NewRouletteWheel([4]float64{3, 2, 1, 0.5}, 0.8, 1, 1, nil)
		accept := HillClimbing{}
		stop := MaxIterations{MaxIterations: 10}
		res := a.Iterate(&initialSolution, &opSelect, &accept, &stop)
		return res
	}
	t.Run("With", func(t *testing.T) {
		res := solve(true)
		if len(res.Statistics.Objectives) != 11 { // initial + 10 iterations
			t.Fatalf("11 objectives expected, actual %d objectives", len(res.Statistics.Objectives))
		}
	})
	t.Run("Without", func(t *testing.T) {
		res := solve(false)
		if len(res.Statistics.Objectives) != 0 {
			t.Fatalf("0 objectives expected, actual %d objectives", len(res.Statistics.Objectives))
		}
	})
}
