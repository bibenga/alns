package alns

type State[O any] interface {
	Objective() O
}
