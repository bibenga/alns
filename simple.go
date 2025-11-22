package alns

import "cmp"

func Iterate(
	initial State,
	destroyOperators []Operator,
	repairOperators []Operator,
	scores [4]float64, // scores for RouletteWheel
	decay float64, // decay for RouletteWheel
	maxIterations int, // maxIterations for MaxIterations
) (*Result, error) {
	compare := cmp.Compare[float64]

	selector, err := NewRouletteWheel(
		compare,
		scores,
		decay,
		len(destroyOperators),
		len(repairOperators),
		nil,
	)
	if err != nil {
		return nil, err
	}

	acceptor := HillClimbing{
		Compare: compare,
	}

	stop := MaxIterations{
		MaxIterations: maxIterations,
	}

	a := ALNS{
		Compare:           compare,
		Rnd:               RuntimeRand,
		CollectObjectives: false,
		DestroyOperators:  destroyOperators,
		RepairOperators:   repairOperators,
		Selector:          &selector,
		Acceptor:          &acceptor,
		Stop:              &stop,
		InitialSolution:   initial,
	}
	return a.Iterate()
}
