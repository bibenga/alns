package alns

func Iterate(
	initial State,
	destroyOperators []Operator,
	repairOperators []Operator,
	scores [4]float64, // scores for RouletteWheel
	decay float64, // decay for RouletteWheel
	maxIterations int, // maxIterations for MaxIterations
) (*Result, error) {
	selector, err := NewRouletteWheel(
		scores,
		decay,
		len(destroyOperators),
		len(repairOperators),
		nil,
	)
	if err != nil {
		return nil, err
	}

	acceptor := HillClimbing{}

	stop := MaxIterations{
		MaxIterations: maxIterations,
	}

	a := ALNS{
		Rnd:               RuntimeRand,
		CollectObjectives: false,
		DestroyOperators:  destroyOperators,
		RepairOperators:   repairOperators,
	}
	return a.Iterate(initial, &selector, &acceptor, &stop)
}
