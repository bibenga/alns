package alns

import (
	"fmt"
	"math/rand/v2"
)

type OperatorSelectionScheme[O any] interface {
	Select(rnd *rand.Rand, best, current State[O]) (deleteOpIndx, repairOpIndx int)
	Update(candidate State[O], deleteOpIndx, repairOpIndx int, outcome Outcome)
}

type RouletteWheel[O any] struct {
	scores          [4]float64
	decay           float64
	numDestroy      int
	numRepair       int
	opCoupling      [][]bool
	dWeights        []float64
	rWeights        []float64
	coupledRIdcs    []int     // used in Select for caching
	coupledRWeights []float64 // used in Select for caching
}

var _ OperatorSelectionScheme[int] = &RouletteWheel[int]{}

func NewRouletteWheel[O any](
	scores [4]float64,
	decay float64,
	numDestroy int,
	numRepair int,
	opCoupling [][]bool,
) (RouletteWheel[O], error) {
	r := RouletteWheel[O]{
		scores:          scores,
		decay:           decay,
		numDestroy:      numDestroy,
		numRepair:       numRepair,
		opCoupling:      opCoupling,
		dWeights:        make([]float64, numDestroy),
		rWeights:        make([]float64, numRepair),
		coupledRIdcs:    make([]int, 0, numRepair),
		coupledRWeights: make([]float64, 0, numRepair),
	}
	for i := range numDestroy {
		r.dWeights[i] = 1
	}
	for i := range numRepair {
		r.rWeights[i] = 1
	}
	if err := r.validate(); err != nil {
		return RouletteWheel[O]{}, err
	}
	return r, nil
}

func (r *RouletteWheel[O]) validate() error {
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

func (r *RouletteWheel[O]) Select(rnd *rand.Rand, best State[O], current State[O]) (int, int) {
	if r.opCoupling != nil {
		// select destroy operator
		dIdx := weightedRandomIndex(rnd, r.dWeights)

		// extract coupled repair indeces and their weight for selected destroy operator
		r.coupledRIdcs = r.coupledRIdcs[:0]
		r.coupledRWeights = r.coupledRWeights[:0]
		for i, v := range r.opCoupling[dIdx] {
			if v {
				r.coupledRIdcs = append(r.coupledRIdcs, i)
				r.coupledRWeights = append(r.coupledRWeights, r.rWeights[i])
			}
		}

		// select repair operator
		rIdx := r.coupledRIdcs[weightedRandomIndex(rnd, r.coupledRWeights)]

		return dIdx, rIdx
	} else {
		dIdx := weightedRandomIndex(rnd, r.dWeights)
		rIdx := weightedRandomIndex(rnd, r.rWeights)
		return dIdx, rIdx
	}
}

func (r *RouletteWheel[O]) Update(candidate State[O], deleteOpIndx int, repairOpIndx int, outcome Outcome) {
	r.dWeights[deleteOpIndx] *= r.decay
	r.dWeights[deleteOpIndx] += (1 - r.decay) * r.scores[outcome]

	r.rWeights[repairOpIndx] *= r.decay
	r.rWeights[repairOpIndx] += (1 - r.decay) * r.scores[outcome]
}
