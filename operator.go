package alns

import "math/rand/v2"

type Operator[O any] func(state State[O], rnd *rand.Rand) (State[O], error)
