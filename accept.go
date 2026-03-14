package alns

import (
	"math/rand/v2"
)

type AcceptanceCriterion interface {
	Accept(rnd *rand.Rand, best, current, candidate State) (bool, error)
}
