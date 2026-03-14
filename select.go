package alns

import (
	"math/rand/v2"
)

type OperatorSelectionScheme interface {
	Select(rnd *rand.Rand, best, current State) (deleteOpIndx, repairOpIndx int, err error)
	Update(candidate State, deleteOpIndx, repairOpIndx int, outcome Outcome) error
}
