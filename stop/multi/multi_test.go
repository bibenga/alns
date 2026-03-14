package multi

import (
	"testing"

	"github.com/bibenga/alns/stop/maxiterations"
)

func TestMultipleStoppingCriterions(t *testing.T) {
	stop := MultipleStoppingCriterions{
		&maxiterations.MaxIterations{MaxIterations: 10},
		&maxiterations.MaxIterations{MaxIterations: 20},
	}

	i := 0
	for {
		if done, err := stop.IsDone(nil, nil, nil); err != nil {
			t.Fatal(err)
		} else if done {
			break
		}
		i++
	}

	if i != 10 {
		t.Fatalf("10 iterations expected, actual %d iterations", i)
	}
}
