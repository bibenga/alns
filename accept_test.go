package alns

import (
	"cmp"
	"testing"
)

func TestHillClimbing(t *testing.T) {
	accept := HillClimbing{
		Compare: cmp.Compare[float64],
	}

	best := FakeState{objective: 2}
	curr := FakeState{objective: 2.1}
	cand := FakeState{objective: 1.9}

	accepted, err := accept.Accept(nil, best, curr, cand)
	if err != nil {
		t.Fatal(err)
	}
	if !accepted {
		t.Fatal("expected to be accepted")
	}

	cand = FakeState{objective: 2.9}
	accepted, err = accept.Accept(nil, best, curr, cand)
	if err != nil {
		t.Fatal(err)
	}
	if accepted {
		t.Fatal("expected not to be accepted")
	}
}
