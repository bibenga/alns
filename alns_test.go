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
	const total = 10000

	opSelect, err := NewRouletteWheel([4]float64{3, 2, 1, 0.5}, 0.8, 1, 1, nil)
	if err != nil {
		t.Fatal(err)
	}
	accept := HillClimbing{}
	stop := MaxIterations{MaxIterations: total}

	lastBest := rand.Float64()
	initialSolution := FakeState{objective: lastBest}

	a := ALNS{
		Rnd:               RuntimeRand,
		CollectObjectives: true,
		Selector:          &opSelect,
		Acceptor:          &accept,
		Stop:              &stop,
		InitialSolution:   &initialSolution,
	}

	bestCount := 0
	destroyCalled := 0
	a.AddDestroyOperator(func(state State, rnd *rand.Rand) (State, error) {
		destroyCalled++
		current := state.(*FakeState)
		destroyed := current.Clone()
		return destroyed, nil
	})

	repairCalled := 0
	a.AddRepairOperator(func(state State, rnd *rand.Rand) (State, error) {
		repairCalled++
		current := state.(*FakeState)
		current.objective = rand.Float64()
		if current.objective < lastBest {
			lastBest = current.objective
			bestCount++
		}
		return current, nil
	})

	res, err := a.Iterate()
	if err != nil {
		t.Fatal(err)
	}

	if destroyCalled != total {
		t.Errorf("%d destroy calls expected, actual %d calls", total, destroyCalled)
	}
	if repairCalled != total {
		t.Errorf("%d repair calls expected, actual %d calls", total, repairCalled)
	}
	if stop.currentIteration != total+1 {
		t.Errorf("%d iterations expected, actual %d calls", total+1, stop.currentIteration)
	}
	if len(res.Statistics.Objectives) != total+1 {
		t.Errorf("%d objectives expected, actual %d objectives", total+1, len(res.Statistics.Objectives))
	}
	repairOperatorCounts := OperatorStatistics{bestCount, 0, 0, total - bestCount}
	if res.Statistics.RepairOperatorCounts[0] != repairOperatorCounts {
		t.Errorf("expected repair opeator statistics %v, actual %v",
			repairOperatorCounts, res.Statistics.RepairOperatorCounts[0])
	}
	rejectOperatorCounts := OperatorStatistics{bestCount, 0, 0, total - bestCount}
	if res.Statistics.DestroyOperatorCounts[0] != rejectOperatorCounts {
		t.Errorf("expected destory opeator statistics %v, actual %v",
			rejectOperatorCounts, res.Statistics.DestroyOperatorCounts[0])
	}
}

func TestAlnsCollectObjectives(t *testing.T) {
	solve := func(collectObjectives bool) *Result {
		opSelect, _ := NewRouletteWheel([4]float64{3, 2, 1, 0.5}, 0.8, 1, 1, nil)
		accept := HillClimbing{}
		stop := MaxIterations{MaxIterations: 10}
		initialSolution := FakeState{objective: 1}
		a := ALNS{
			Rnd:               rand.New(rand.NewPCG(1, 2)),
			CollectObjectives: collectObjectives,
			DestroyOperators: []Operator{
				func(state State, rnd *rand.Rand) (State, error) { return state, nil },
			},
			RepairOperators: []Operator{
				func(state State, rnd *rand.Rand) (State, error) { return state, nil },
			},
			Selector:        &opSelect,
			Acceptor:        &accept,
			Stop:            &stop,
			InitialSolution: initialSolution,
		}
		res, err := a.Iterate()
		if err != nil {
			t.Fatal(err)
		}
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
