package hillclimbing

import (
	"testing"

	"github.com/bibenga/alns/internal/testutil"
)

func TestHillClimbing(t *testing.T) {
	accept := HillClimbing{}

	best := testutil.NewFakeState(2)
	curr := testutil.NewFakeState(2.1)
	cand := testutil.NewFakeState(1.9)

	accepted, err := accept.Accept(nil, best, curr, cand)
	if err != nil {
		t.Fatal(err)
	}
	if !accepted {
		t.Fatal("expected to be accepted")
	}

	cand = testutil.NewFakeState(2.9)
	accepted, err = accept.Accept(nil, best, curr, cand)
	if err != nil {
		t.Fatal(err)
	}
	if accepted {
		t.Fatal("expected not to be accepted")
	}
}
