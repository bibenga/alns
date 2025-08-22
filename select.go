package alns

import (
	"fmt"
	"math/rand/v2"
)

type OperatorSelectionScheme interface {
	Select(rnd *rand.Rand, best, current State) (deleteOpIndx, repairOpIndx int)
	Update(candidate State, deleteOpIndx, repairOpIndx int, outcome Outcome)
}

type RouletteWheel struct {
	scores     [4]float64
	decay      float64
	numDestroy int
	numRepair  int
	opCoupling [][]bool
	dWeights   []float64
	rWeights   []float64
}

var _ OperatorSelectionScheme = &RouletteWheel{}

func NewRouletteWheel(
	scores [4]float64,
	decay float64,
	numDestroy int,
	numRepair int,
	opCoupling [][]bool,
) (RouletteWheel, error) {
	r := RouletteWheel{
		scores:     scores,
		decay:      decay,
		numDestroy: numDestroy,
		numRepair:  numRepair,
		opCoupling: opCoupling,
		dWeights:   make([]float64, numDestroy),
		rWeights:   make([]float64, numRepair),
	}
	for i := range numDestroy {
		r.dWeights[i] = 1
	}
	for i := range numRepair {
		r.rWeights[i] = 1
	}
	if err := r.validate(); err != nil {
		return RouletteWheel{}, err
	}
	return r, nil
}

func (r *RouletteWheel) validate() error {
	if min(r.scores[0], r.scores[1], r.scores[2], r.scores[3]) < 0 {
		return fmt.Errorf("negative scores are not understood")
	}

	if !(0 <= r.decay && r.decay <= 1) {
		return fmt.Errorf("decay outside [0, 1] not understood")
	}

	if r.opCoupling != nil {
		if len(r.opCoupling) == 0 {
			return fmt.Errorf("coupling matrix of shape (%d, %d), expected (%d, %d)",
				0, 0, r.numDestroy, r.numRepair)
		}
		rows := len(r.opCoupling)
		cols := len(r.opCoupling[0])
		for i, row := range r.opCoupling {
			if len(row) != cols {
				return fmt.Errorf("the number of columns in a row %d does not match the expected %d",
					i, cols)
			}
		}
		if rows != r.numDestroy || cols != r.numRepair {
			return fmt.Errorf("coupling matrix of shape (%d, %d), expected (%d, %d)",
				rows, cols, r.numDestroy, r.numRepair)
		}

		for i, row := range r.opCoupling {
			isCoupled := false
			for _, value := range row {
				if value {
					isCoupled = true
					break
				}
			}
			if !isCoupled {
				return fmt.Errorf("destroy operator %d has no coupled repair operators", i)
			}
		}
	}

	return nil
}

func (r *RouletteWheel) Select(rnd *rand.Rand, best State, current State) (int, int) {
	if r.opCoupling != nil {
		dIdx := weightedRandomIndex(rnd, r.dWeights)
		coupledRIdcs := r.flatTrue(r.opCoupling[dIdx])
		rIdx := coupledRIdcs[weightedRandomIndex(rnd, r.extract(r.rWeights, coupledRIdcs))]
		return dIdx, rIdx
	} else {
		dIdx := weightedRandomIndex(rnd, r.dWeights)
		rIdx := weightedRandomIndex(rnd, r.rWeights)
		return dIdx, rIdx
	}
}

func (r *RouletteWheel) Update(candidate State, deleteOpIndx int, repairOpIndx int, outcome Outcome) {
	r.dWeights[deleteOpIndx] *= r.decay
	r.dWeights[deleteOpIndx] += (1 - r.decay) * r.scores[outcome]

	r.rWeights[repairOpIndx] *= r.decay
	r.rWeights[repairOpIndx] += (1 - r.decay) * r.scores[outcome]
}

func (r *RouletteWheel) flatTrue(s []bool) []int {
	res := make([]int, 0, len(s))
	for i, v := range s {
		if v {
			res = append(res, i)
		}
	}
	return res
}

func (r *RouletteWheel) extract(s []float64, indices []int) []float64 {
	res := make([]float64, len(indices))
	for i := range indices {
		res[i] = s[i]
	}
	return res
}
