package lateacceptancehillclimbing

import (
	"fmt"
	"testing"

	"github.com/bibenga/alns/internal/testutil"
)

func TestRaisesInvalidLookbackPeriod(t *testing.T) {
	cases := []int{-1}
	for _, lp := range cases {
		t.Run(fmt.Sprintf("lookback_period=%d", lp), func(t *testing.T) {
			_, err := NewLateAcceptanceHillClimbing(lp, false, false)
			if err == nil {
				t.Errorf("expected error for lookback_period=%d, got nil", lp)
			}
		})
	}
}

func TestProperties(t *testing.T) {
	cases := []struct {
		lookbackPeriod int
		greedy         bool
		betterHistory  bool
	}{
		{0, true, true},
		{1, false, true},
		{10, true, false},
		{100, false, false},
	}
	for _, tc := range cases {
		t.Run(
			fmt.Sprintf("lookback=%d/greedy=%v/betterHistory=%v", tc.lookbackPeriod, tc.greedy, tc.betterHistory),
			func(t *testing.T) {
				lahc, err := NewLateAcceptanceHillClimbing(tc.lookbackPeriod, tc.greedy, tc.betterHistory)
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				if lahc.LookbackPeriod != tc.lookbackPeriod {
					t.Errorf("LookbackPeriod: got %d, want %d", lahc.LookbackPeriod, tc.lookbackPeriod)
				}
				if lahc.Greedy != tc.greedy {
					t.Errorf("Greedy: got %v, want %v", lahc.Greedy, tc.greedy)
				}
				if lahc.BetterHistory != tc.betterHistory {
					t.Errorf("BetterHistory: got %v, want %v", lahc.BetterHistory, tc.betterHistory)
				}
			},
		)
	}
}

func TestZeroLookbackPeriod(t *testing.T) {
	cases := []struct {
		greedy        bool
		betterHistory bool
	}{
		{false, false}, {true, false}, {false, true}, {true, true},
	}
	for _, tc := range cases {
		t.Run(
			fmt.Sprintf("greedy=%v/betterHistory=%v", tc.greedy, tc.betterHistory),
			func(t *testing.T) {
				lahc, _ := NewLateAcceptanceHillClimbing(0, tc.greedy, tc.betterHistory)

				accept, err := lahc.Accept(nil, testutil.Zero(), testutil.Two(), testutil.One())
				if err != nil || !accept {
					t.Error("expected accept: Zero best, Two current, One candidate")
				}
				accept, err = lahc.Accept(nil, testutil.Zero(), testutil.One(), testutil.One())
				if err != nil || accept {
					t.Error("expected reject: Zero best, One current, One candidate")
				}
				accept, err = lahc.Accept(nil, testutil.Zero(), testutil.Zero(), testutil.Zero())
				if err != nil || accept {
					t.Error("expected reject: all Zero")
				}
				accept, err = lahc.Accept(nil, testutil.Zero(), testutil.Two(), testutil.One())
				if err != nil || !accept {
					t.Error("expected accept: Zero best, Two current, One candidate")
				}
			},
		)
	}
}

func TestAccept(t *testing.T) {
	for _, lp := range []int{0, 3, 10, 50} {
		t.Run(
			fmt.Sprintf("lookback=%d", lp),
			func(t *testing.T) {
				lahc, _ := NewLateAcceptanceHillClimbing(lp, false, false)

				ok, err := lahc.Accept(nil, testutil.Zero(), testutil.Two(), testutil.One())
				if err != nil || !ok {
					t.Error("expected initial accept")
				}
				for i := range lp {
					// Historical current is 2, candidate is 1 → should accept.
					ok, err = lahc.Accept(nil, testutil.Zero(), testutil.One(), testutil.One())
					if err != nil || !ok {
						t.Errorf("iter=%d: expected accept", i)
					}
				}
			},
		)
	}
}

func TestReject(t *testing.T) {
	for _, lp := range []int{0, 3, 10, 50} {
		t.Run(
			fmt.Sprintf("lookback=%d", lp),
			func(t *testing.T) {
				lahc, _ := NewLateAcceptanceHillClimbing(lp, false, false)

				for i := range lp {
					ok, err := lahc.Accept(nil, testutil.Zero(), testutil.One(), testutil.Zero())
					if err != nil || !ok {
						t.Errorf("iter=%d: expected accept", i)
					}
				}

				// Historical current is 1, candidate is 1 → reject (not strictly better).
				ok, err := lahc.Accept(nil, testutil.Zero(), testutil.Zero(), testutil.One())
				if err != nil || ok {
					t.Error("expected reject")
				}
			},
		)
	}
}

func TestGreedyAccept(t *testing.T) {
	for _, lp := range []int{0, 3, 10, 50} {
		t.Run(
			fmt.Sprintf("lookback=%d", lp),
			func(t *testing.T) {
				lahc, _ := NewLateAcceptanceHillClimbing(lp, true, false)

				ok, err := lahc.Accept(nil, testutil.Zero(), testutil.Zero(), testutil.Two())
				if err != nil || ok {
					t.Error("expected reject: candidate(2) worse than current(0)")
				}
				for i := range lp {
					// Candidate(1) < current(2), accepted greedily despite historical=1.
					ok, err = lahc.Accept(nil, testutil.Zero(), testutil.Two(), testutil.One())
					if err != nil || !ok {
						t.Errorf("iter=%d: expected greedy accept", i)
					}
				}
			},
		)
	}
}

func TestBetterHistorySmallExample(t *testing.T) {
	lahc, _ := NewLateAcceptanceHillClimbing(1, false, true)

	ok, err := lahc.Accept(nil, testutil.Zero(), testutil.One(), testutil.Zero())
	if err != nil || !ok {
		t.Error("step1: expected accept")
	}
	ok, err = lahc.Accept(nil, testutil.Zero(), testutil.Two(), testutil.Zero())
	if err != nil || !ok {
		t.Error("step2: expected accept")
	}

	// Previous current stays at 1 because 2 was not better than historical.
	ok, err = lahc.Accept(nil, testutil.Zero(), testutil.Zero(), testutil.One())
	if err != nil || ok {
		t.Error("step3: expected reject")
	}

	// Previous current is updated to Zero.
	ok, err = lahc.Accept(nil, testutil.Zero(), testutil.Zero(), testutil.Zero())
	if err != nil || ok {
		t.Error("step4: expected reject")
	}
}

func TestBetterHistoryReject(t *testing.T) {
	for _, lp := range []int{3, 10, 50} {
		t.Run(
			fmt.Sprintf("lookback=%d", lp),
			func(t *testing.T) {
				lahc, _ := NewLateAcceptanceHillClimbing(lp, false, true)

				for i := range lp {
					ok, err := lahc.Accept(nil, testutil.Zero(), testutil.One(), testutil.Two())
					if err != nil || ok {
						t.Errorf("phase1 iter=%d: expected reject", i)
					}
				}
				for i := range lp {
					// Current solutions not stored because they are worse than historical.
					ok, err := lahc.Accept(nil, testutil.Zero(), testutil.Two(), testutil.Two())
					if err != nil || ok {
						t.Errorf("phase2 iter=%d: expected reject", i)
					}
				}
				for i := range lp {
					// Candidates(1) do not improve historical solutions(1).
					ok, err := lahc.Accept(nil, testutil.Zero(), testutil.Two(), testutil.One())
					if err != nil || ok {
						t.Errorf("phase3 iter=%d: expected reject", i)
					}
				}
			},
		)
	}
}

func TestFullExample(t *testing.T) {
	lahc, _ := NewLateAcceptanceHillClimbing(1, true, true)

	// Iteration 1: candidate(1) vs current(2); store current in history.
	ok, err := lahc.Accept(nil, testutil.Zero(), testutil.Two(), testutil.One())
	if err != nil || !ok {
		t.Error("iter1: expected accept")
	}
	// Iteration 2: candidate(1) accepted via late-acceptance(2); history → 1.
	ok, err = lahc.Accept(nil, testutil.Zero(), testutil.One(), testutil.One())
	if err != nil || !ok {
		t.Error("iter2: expected accept")
	}
	// Iteration 3: candidate(1) accepted greedily vs current(2);
	// history not updated because current(2) does not improve historical(1).
	ok, err = lahc.Accept(nil, testutil.Zero(), testutil.Two(), testutil.One())
	if err != nil || !ok {
		t.Error("iter3: expected accept")
	}
	// Iteration 4: candidate(1) not accepted — doesn't beat current(1)
	// nor historical(1).
	ok, err = lahc.Accept(nil, testutil.Zero(), testutil.One(), testutil.One())
	if err != nil || ok {
		t.Error("iter4: expected reject")
	}
}
