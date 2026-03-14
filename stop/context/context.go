package context

import (
	"context"
	"math/rand/v2"

	"github.com/bibenga/alns"
)

type Context struct {
	Context context.Context
}

var _ alns.StoppingCriterion = &Context{}

func NewContext(context context.Context) Context {
	return Context{
		Context: context,
	}
}

func (s *Context) IsDone(rnd *rand.Rand, best, current alns.State) (bool, error) {
	select {
	case <-s.Context.Done():
		return true, nil
	default:
		return false, nil
	}
}
