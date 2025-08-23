package alns

import (
	"cmp"
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
	a := ALNS[float64]{
		Rnd:               RuntimeRand,
		Compare:           cmp.Compare[float64],
		CollectObjectives: true,
	}

	lastBest := rand.Float64()
	bestCount := 0
	destroyCalled := 0
	a.AddDestroyOperator(func(state State[float64], rnd *rand.Rand) (State[float64], error) {
		destroyCalled++
		current := state.(*FakeState)
		destroyed := current.Clone()
		return destroyed, nil
	})
	repairCalled := 0
	a.AddRepairOperator(func(state State[float64], rnd *rand.Rand) (State[float64], error) {
		repairCalled++
		current := state.(*FakeState)
		current.objective = rand.Float64()
		if current.objective < lastBest {
			lastBest = current.objective
			bestCount++
		}
		return current, nil
	})

	total := 10000

	initialSolution := FakeState{objective: lastBest}
	opSelect, _ := NewRouletteWheel[float64]([4]float64{3, 2, 1, 0.5}, 0.8, 1, 1, nil)
	accept := HillClimbing[float64]{Compare: cmp.Compare[float64]}
	stop := MaxIterations[float64]{MaxIterations: total}
	res, err := a.Iterate(&initialSolution, &opSelect, &accept, &stop)
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
	solve := func(collectObjectives bool) *Result[float64] {
		a := ALNS[float64]{
			Rnd:               rand.New(rand.NewPCG(1, 2)),
			Compare:           cmp.Compare[float64],
			CollectObjectives: collectObjectives,
			DestroyOperators: []Operator[float64]{
				func(state State[float64], rnd *rand.Rand) (State[float64], error) { return state, nil },
			},
			RepairOperators: []Operator[float64]{
				func(state State[float64], rnd *rand.Rand) (State[float64], error) { return state, nil },
			},
		}
		initialSolution := FakeState{objective: 1}
		opSelect, _ := NewRouletteWheel[float64]([4]float64{3, 2, 1, 0.5}, 0.8, 1, 1, nil)
		accept := HillClimbing[float64]{Compare: cmp.Compare[float64]}
		stop := MaxIterations[float64]{MaxIterations: 10}
		res, err := a.Iterate(&initialSolution, &opSelect, &accept, &stop)
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
