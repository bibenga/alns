package alns

import "math/rand/v2"

type Operator func(state State, rnd *rand.Rand) State

type operator struct {
	call Operator
	name string
}
