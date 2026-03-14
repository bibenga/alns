package maxruntime

import (
	"math/rand/v2"
	"time"

	"github.com/bibenga/alns"
)

type MaxRuntime struct {
	MaxRuntime time.Duration
	started    time.Time
}

var _ alns.StoppingCriterion = &MaxRuntime{}

func NewMaxRuntime(maxRuntime time.Duration) MaxRuntime {
	return MaxRuntime{
		MaxRuntime: maxRuntime,
	}
}

func (s *MaxRuntime) IsDone(rnd *rand.Rand, best, current alns.State) (bool, error) {
	if s.started.IsZero() {
		s.started = time.Now()
		return false, nil
	}
	return time.Since(s.started) > s.MaxRuntime, nil
}
