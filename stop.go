package alns

import (
	"math/rand/v2"
)

type StoppingCriterion interface {
	IsDone(rnd *rand.Rand, best, current State) (bool, error)
}
