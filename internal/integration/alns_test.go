package integration

import (
	"math/rand/v2"
	"testing"

	"github.com/bibenga/alns"
	"github.com/bibenga/alns/accept/hillclimbing"
	"github.com/bibenga/alns/internal/testutil"
	"github.com/bibenga/alns/select/roulettewheel"
	"github.com/bibenga/alns/stop/maxiterations"
)

func TestAlns(t *testing.T) {
	const total = 10000

	opSelect, err := roulettewheel.NewRouletteWheel([4]float64{3, 2, 1, 0.5}, 0.8, 1, 1, nil)
	if err != nil {
		t.Fatal(err)
	}
	accept := hillclimbing.NewHillClimbing()
	stop := maxiterations.NewMaxIterations(total)

	lastBest := rand.Float64()
	initSol := testutil.NewFakeState(lastBest)

	a := alns.ALNS{
		Rnd:               alns.RuntimeRand,
		CollectObjectives: true,
	}

	bestCount := uint(0)
	destroyCalled := uint(0)
	a.DestroyOperators = append(a.DestroyOperators,
		func(state alns.State, rnd *rand.Rand) (alns.State, error) {
			destroyCalled++
			current := state.(*testutil.FakeState)
			destroyed := current.Clone()
			return destroyed, nil
		},
	)

	repairCalled := 0
	a.RepairOperators = append(a.RepairOperators,
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
	)

	res, err := a.Iterate(&initSol, &opSelect, &accept, &stop)
	if err != nil {
		t.Fatal(err)
	}

	if destroyCalled != total {
		t.Errorf("%d destroy calls expected, actual %d calls", total, destroyCalled)
	}
	if repairCalled != total {
		t.Errorf("%d repair calls expected, actual %d calls", total, repairCalled)
	}
	// if stop.currentIteration != total+1 {
	// 	t.Errorf("%d iterations expected, actual %d calls", total+1, stop.currentIteration)
	// }
	if len(res.Statistics.Objectives) != total+1 {
		t.Errorf("%d objectives expected, actual %d objectives", total+1, len(res.Statistics.Objectives))
	}
	repairOperatorCounts := alns.OperatorCounts{bestCount, 0, 0, total - bestCount}
	if res.Statistics.RepairOperatorCounts[0] != repairOperatorCounts {
		t.Errorf("expected repair opeator statistics %v, actual %v",
			repairOperatorCounts, res.Statistics.RepairOperatorCounts[0])
	}
	rejectOperatorCounts := alns.OperatorCounts{uint(bestCount), 0, 0, total - bestCount}
	if res.Statistics.DestroyOperatorCounts[0] != rejectOperatorCounts {
		t.Errorf("expected destory opeator statistics %v, actual %v",
			rejectOperatorCounts, res.Statistics.DestroyOperatorCounts[0])
	}
}

func TestAlnsCollectObjectives(t *testing.T) {
	solve := func(collectObjectives bool) *alns.Result {
		opSelect, _ := roulettewheel.NewRouletteWheel([4]float64{3, 2, 1, 0.5}, 0.8, 1, 1, nil)
		accept := hillclimbing.NewHillClimbing()
		stop := maxiterations.MaxIterations{MaxIterations: 10}
		initSol := testutil.FakeState{Value: 1}
		a := alns.ALNS{
			Rnd:               rand.New(rand.NewPCG(1, 2)),
			CollectObjectives: collectObjectives,
			DestroyOperators: []alns.Operator{
				func(state alns.State, rnd *rand.Rand) (alns.State, error) { return state, nil },
			},
			RepairOperators: []alns.Operator{
				func(state alns.State, rnd *rand.Rand) (alns.State, error) { return state, nil },
			},
		}
		res, err := a.Iterate(initSol, &opSelect, &accept, &stop)
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
