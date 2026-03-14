package simple

import (
	"math/rand/v2"
	"testing"

	"github.com/bibenga/alns"
	"github.com/bibenga/alns/internal/testutil"
)

func TestIterate(t *testing.T) {
	lastBest := rand.Float64()
	initialSolution := testutil.FakeState{Value: lastBest}

	bestCount := 0
	destroyCalled := 0

	destroyOperators := []alns.Operator{
		func(state alns.State, rnd *rand.Rand) (alns.State, error) {
			destroyCalled++
			current := state.(*testutil.FakeState)
			destroyed := current.Clone()
			return destroyed, nil
		},
	}

	repairCalled := 0
	repairOperators := []alns.Operator{
		func(state alns.State, rnd *rand.Rand) (alns.State, error) {
			repairCalled++
			current := state.(*testutil.FakeState)
			current.Value = rand.Float64()
			if current.Value < lastBest {
				lastBest = current.Value
				bestCount++
			}
			return current, nil
		},
	}

	const total = 100

	res, err := Iterate(
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
