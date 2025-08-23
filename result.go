package alns

type Result[O any] struct {
	BestState  State[O]
	Statistics Statistics[O]
}
