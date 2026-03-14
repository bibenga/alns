package simple

import (
	"github.com/bibenga/alns"
	"github.com/bibenga/alns/accept/hillclimbing"
	"github.com/bibenga/alns/select/roulettewheel"
	"github.com/bibenga/alns/stop/maxiterations"
)

func Iterate(
	initial alns.State,
	destroyOperators []alns.Operator,
	repairOperators []alns.Operator,
	scores [4]float64, // scores for RouletteWheel
	decay float64, // decay for RouletteWheel
	maxIterations int, // maxIterations for MaxIterations
) (*alns.Result, error) {
	selector, err := roulettewheel.NewRouletteWheel(
		scores,
		decay,
		len(destroyOperators),
		len(repairOperators),
		nil,
	)
	if err != nil {
		return nil, err
	}

	acceptor := hillclimbing.NewHillClimbing()
	stop := maxiterations.NewMaxIterations(maxIterations)

	a := alns.ALNS{
		Rnd:               alns.RuntimeRand,
		CollectObjectives: false,
		DestroyOperators:  destroyOperators,
		RepairOperators:   repairOperators,
	}
	return a.Iterate(initial, &selector, &acceptor, &stop)
}
