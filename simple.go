package alns

import (
	"cmp"
)

func SimpleIterate[O cmp.Ordered](
	initial State[O],
	destroyOperators []Operator[O],
	repairOperators []Operator[O],
	scores [4]float64, // scores for RouletteWheel
	decay float64, // decay for RouletteWheel
	maxIterations int, // maxIterations for MaxIterations
) (*Result[O], error) {
	selector, err := NewRouletteWheel[O](
		scores,
		decay,
		len(destroyOperators),
		len(repairOperators),
		nil,
	)
	if err != nil {
		return nil, err
	}

	acceptor := HillClimbing[O]{Compare: cmp.Compare[O]}
	stop := MaxIterations[O]{MaxIterations: maxIterations}

	a := ALNS[O]{
		Rnd:               RuntimeRand,
		Compare:           cmp.Compare[O],
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
