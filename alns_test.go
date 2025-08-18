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
	a := NewDefault()
	a.AddDestroyOperator(func(state State, rnd *rand.Rand) State {
		current := state.(*FakeState)
		destroyed := current.Clone()
		return destroyed
	})
	a.AddRepairOperator(func(state State, rnd *rand.Rand) State {
		current := state.(*FakeState)
		current.objective = rand.Float64()
		return current
	})

	initialSolution := FakeState{objective: 1}
	opSelect := NewRouletteWheel([4]float64{3, 2, 1, 0.5}, 0.8, 1, 1, nil)
	accept := HillClimbing{}
	stop := MaxIterations{MaxIterations: 10}
	res := a.Iterate(&initialSolution, &opSelect, &accept, &stop)
	t.Log(*res.BestState.(*FakeState))
	t.Log(res.Statistics.Objectives)
	t.Log(res.Statistics.DestroyOperatorCounts)
	t.Log(res.Statistics.RepairOperatorCounts)
}
