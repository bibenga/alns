package math

import (
	"math/rand/v2"
	"time"
)

func WeightedRandomIndex(rnd *rand.Rand, weights []float64) int {
	if len(weights) == 0 {
		panic("invalid weights")
	}
	if len(weights) == 1 {
		return 0
	}
	sum := Sum(weights)
	value := rnd.Float64() * sum // adjusted value
	for i, weight := range weights {
		value -= weight
		if value <= 0 {
			return i
		}
	}
	// we will only be here when errors accumulate
	return len(weights) - 1
}

func Sum(weights []float64) float64 {
	sum := 0.0
	for _, w := range weights {
		sum += w
	}
	return sum
}

type Request struct {
	// hot data
	Id                   uint32 // 22 bit
	SrcId                uint16 // 9 bit
	DstId                uint16 // 9 bit
	Vol                  uint32 // 24 bit
	Quant                int16  // 12 bit
	AvailabilityDateUnix int64
	DeadlineDateUnix     int64
	Transit              []uint8
	// cold data
	Sla              time.Duration
	Cid              string
	AvailabilityDate time.Time
	DeadlineDate     time.Time
}
