package testutil

type FakeState struct {
	Value float64
}

func NewFakeState(objective float64) FakeState {
	return FakeState{
		Value: objective,
	}
}

func (s FakeState) Clone() *FakeState {
	return &FakeState{}
}

func (s FakeState) Objective() float64 {
	return s.Value
}
