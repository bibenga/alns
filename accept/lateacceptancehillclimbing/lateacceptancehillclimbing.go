package lateacceptancehillclimbing

import (
	"errors"
	"math/rand/v2"

	"github.com/bibenga/alns"
)

type LateAcceptanceHillClimbing struct {
	LookbackPeriod int
	Greedy         bool
	BetterHistory  bool

	history []float64
	head    int
	size    int
}

var _ alns.AcceptanceCriterion = &LateAcceptanceHillClimbing{}

func NewLateAcceptanceHillClimbing(
	lookbackPeriod int,
	greedy bool,
	betterHistory bool,
) (*LateAcceptanceHillClimbing, error) {
	if lookbackPeriod < 0 {
		return nil, errors.New("lookbackPeriod must be a non-negative integer")
	}
	var buf []float64
	if lookbackPeriod > 0 {
		buf = make([]float64, lookbackPeriod)
	}
	return &LateAcceptanceHillClimbing{
		LookbackPeriod: lookbackPeriod,
		Greedy:         greedy,
		BetterHistory:  betterHistory,
		history:        buf,
	}, nil
}

func (l *LateAcceptanceHillClimbing) Accept(rnd *rand.Rand, best, current, candidate alns.State) (bool, error) {
	candObj := candidate.Objective()
	currObj := current.Objective()

	if l.size == 0 {
		l.push(currObj)
		return candObj < currObj, nil
	}

	oldest := l.oldest()
	res := candObj < oldest

	if !res && l.Greedy {
		res = candObj < currObj
	}

	if l.BetterHistory {
		l.push(min(currObj, oldest))
	} else {
		l.push(currObj)
	}

	return res, nil
}

func (l *LateAcceptanceHillClimbing) push(v float64) {
	if l.LookbackPeriod == 0 {
		return
	}
	idx := (l.head + l.size) % l.LookbackPeriod
	l.history[idx] = v
	if l.size < l.LookbackPeriod {
		l.size++
	} else {
		l.head = (l.head + 1) % l.LookbackPeriod
	}
}

func (l *LateAcceptanceHillClimbing) oldest() float64 {
	return l.history[l.head]
}
