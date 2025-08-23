package alns

import (
	"fmt"
)

type Outcome int

const (
	Best   Outcome = iota // Candidate solution is a new global best
	Better                // Candidate solution is better than the current incumbent
	Accept                // Candidate solution is accepted
	Reject                // Candidate solution is rejected
)

func (o Outcome) String() string {
	switch o {
	case Best:
		return "Best"
	case Better:
		return "Better"
	case Accept:
		return "Accept"
	case Reject:
		return "Reject"
	default:
		return fmt.Sprintf("%%!Outcome(%d)", o)
	}
}
