package alns

import (
	"fmt"
	"math/rand/v2"
)

type OperatorSelectionScheme interface {
	Select(rnd *rand.Rand, best, current State) (deleteOpIndx, repairOpIndx int, err error)
	Update(candidate State, deleteOpIndx, repairOpIndx int, outcome Outcome) error
}

// The `RouletteWheel` scheme updates operator weights as a convex combination of the current weight, and the new score.
type RouletteWheel struct {
	scores          [4]float64 // representing the weight updates when the candidate solution results in a new global
	decay           float64    // operator decay parameter :math:`\theta \in [0, 1]`
	numDestroy      int        // number of destroy operators
	numRepair       int        // number of repair operators
	opCoupling      [][]bool   // boolean matrix that indicates coupling between destroy and repair operators
	dWeights        []float64  // the weights of the destroy operators
	rWeights        []float64  // the weights of the repair operators
	coupledRIdcs    []int      // used in Select for caching
	coupledRWeights []float64  // used in Select for caching
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
	if opCoupling != nil {
		r.coupledRIdcs = make([]int, 0, numRepair)
		r.coupledRWeights = make([]float64, 0, numRepair)
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

func (r *RouletteWheel) Select(rnd *rand.Rand, best State, current State) (int, int, error) {
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

		return dIdx, rIdx, nil
	} else {
		dIdx := weightedRandomIndex(rnd, r.dWeights)
		rIdx := weightedRandomIndex(rnd, r.rWeights)
		return dIdx, rIdx, nil
	}
}

func (r *RouletteWheel) Update(candidate State, deleteOpIndx int, repairOpIndx int, outcome Outcome) error {
	r.dWeights[deleteOpIndx] *= r.decay
	r.dWeights[deleteOpIndx] += (1 - r.decay) * r.scores[outcome]

	r.rWeights[repairOpIndx] *= r.decay
	r.rWeights[repairOpIndx] += (1 - r.decay) * r.scores[outcome]

	return nil
}
