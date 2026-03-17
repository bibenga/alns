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

func Zero() FakeState {
	return FakeState{Value: 0}
}

func One() FakeState {
	return FakeState{Value: 1}
}

func Two() FakeState {
	return FakeState{Value: 2}
}
