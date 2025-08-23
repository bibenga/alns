package alns

import (
	"math/rand/v2"
	"testing"
)

func TestSimpleIterate(t *testing.T) {
	lastBest := rand.Float64()
	initialSolution := FakeState{objective: lastBest}

	bestCount := 0
	destroyCalled := 0

	destroyOperators := []Operator[float64]{
		func(state State[float64], rnd *rand.Rand) (State[float64], error) {
			destroyCalled++
			current := state.(*FakeState)
			destroyed := current.Clone()
			return destroyed, nil
		},
	}

	repairCalled := 0
	repairOperators := []Operator[float64]{
		func(state State[float64], rnd *rand.Rand) (State[float64], error) {
			repairCalled++
			current := state.(*FakeState)
			current.objective = rand.Float64()
			if current.objective < lastBest {
				lastBest = current.objective
				bestCount++
			}
			return current, nil
		},
	}

	const total = 100

	res, err := SimpleIterate(
		&initialSolution,
		destroyOperators,
		repairOperators,
		[4]float64{3, 2, 1, 0.5},
		0.8,
		total,
	)
	if err != nil {
		t.Fatal(err)
	}
	if res == nil {
		t.Fatal("res is nil")
	}

	if destroyCalled != total {
		t.Errorf("%d destroy calls expected, actual %d calls", total, destroyCalled)
	}
	if repairCalled != total {
		t.Errorf("%d repair calls expected, actual %d calls", total, repairCalled)
	}
}
