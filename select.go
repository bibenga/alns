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
	compare         CompareFunc
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
	compare CompareFunc,
	scores [4]float64,
	decay float64,
	numDestroy int,
	numRepair int,
	opCoupling [][]bool,
) (RouletteWheel, error) {
	r := RouletteWheel{
		compare:    compare,
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

func (s *RouletteWheel) validate() error {
	if min(s.scores[0], s.scores[1], s.scores[2], s.scores[3]) < 0 {
		return fmt.Errorf("negative scores are not understood")
	}

	if !(0 <= s.decay && s.decay <= 1) {
		return fmt.Errorf("decay outside [0, 1] not understood")
	}

	if s.opCoupling != nil {
		if len(s.opCoupling) == 0 {
			return fmt.Errorf("coupling matrix of shape (%d, %d), expected (%d, %d)",
				0, 0, s.numDestroy, s.numRepair)
		}
		rows := len(s.opCoupling)
		cols := len(s.opCoupling[0])
		for i, row := range s.opCoupling {
			if len(row) != cols {
				return fmt.Errorf("the number of columns in a row %d does not match the expected %d",
					i, cols)
			}
		}
		if rows != s.numDestroy || cols != s.numRepair {
			return fmt.Errorf("coupling matrix of shape (%d, %d), expected (%d, %d)",
				rows, cols, s.numDestroy, s.numRepair)
		}

		for i, row := range s.opCoupling {
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

func (s *RouletteWheel) Select(rnd *rand.Rand, best State, current State) (int, int, error) {
	if s.opCoupling != nil {
		// select destroy operator
		dIdx := weightedRandomIndex(s.compare, rnd, s.dWeights)

		// extract coupled repair indeces and their weight for selected destroy operator
		s.coupledRIdcs = s.coupledRIdcs[:0]
		s.coupledRWeights = s.coupledRWeights[:0]
		for i, v := range s.opCoupling[dIdx] {
			if v {
				s.coupledRIdcs = append(s.coupledRIdcs, i)
				s.coupledRWeights = append(s.coupledRWeights, s.rWeights[i])
			}
		}

		// select repair operator
		rIdx := s.coupledRIdcs[weightedRandomIndex(s.compare, rnd, s.coupledRWeights)]

		return dIdx, rIdx, nil
	} else {
		dIdx := weightedRandomIndex(s.compare, rnd, s.dWeights)
		rIdx := weightedRandomIndex(s.compare, rnd, s.rWeights)
		return dIdx, rIdx, nil
	}
}

func (s *RouletteWheel) Update(candidate State, deleteOpIndx int, repairOpIndx int, outcome Outcome) error {
	s.dWeights[deleteOpIndx] *= s.decay
	s.dWeights[deleteOpIndx] += (1 - s.decay) * s.scores[outcome]

	s.rWeights[repairOpIndx] *= s.decay
	s.rWeights[repairOpIndx] += (1 - s.decay) * s.scores[outcome]

	return nil
}
