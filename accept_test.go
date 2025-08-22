package alns

import (
	"testing"
)

func TestHillClimbing(t *testing.T) {
	accept := HillClimbing{}

	best := FakeState{objective: 2}
	curr := FakeState{objective: 2.1}
	cand := FakeState{objective: 1.9}

	accepted := accept.Accept(nil, best, curr, cand)
	if !accepted {
		t.Fatal("expected to be accepted")
	}

	cand = FakeState{objective: 2.9}
	accepted = accept.Accept(nil, best, curr, cand)
	if accepted {
		t.Fatal("expected not to be accepted")
	}
}
