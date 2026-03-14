package multi

import (
	"math/rand/v2"

	"github.com/bibenga/alns"
)

type StoppingCriterions []alns.StoppingCriterion

var _ alns.StoppingCriterion = StoppingCriterions{}

func NewStoppingCriterions(criterions ...alns.StoppingCriterion) StoppingCriterions {
	return criterions
}

func (s StoppingCriterions) IsDone(rnd *rand.Rand, best, current alns.State) (bool, error) {
	if len(s) == 0 {
		panic("no criterias were specified")
	}
	for _, c := range s {
		if done, err := c.IsDone(rnd, best, current); err != nil {
			return true, err
		} else if done {
			return true, nil
		}
	}
	return false, nil
}
