package alns

import (
	"fmt"
)

type Outcome int

const (
	BEST Outcome = iota
	BETTER
	ACCEPT
	REJECT
)

func (o Outcome) String() string {
	switch o {
	case BEST:
		return "BEST"
	case BETTER:
		return "BETTER"
	case ACCEPT:
		return "ACCEPT"
	case REJECT:
		return "REJECT"
	default:
		return fmt.Sprintf("%%!Outcome(%d)", o)
	}
}
