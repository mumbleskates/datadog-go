package statsd

import (
	"math"
	"sync/atomic"
)

// This is just a simple wrapper for 64 bits that we read and write atomically
// to do compare-and-swaps while accumulating.
type atomicFloat struct {
	bits uint64
}

func atomicFloat64(v float64) atomicFloat {
	return atomicFloat{math.Float64bits(v)}
}

func (a *atomicFloat) Float64() float64 {
	return math.Float64frombits(atomic.LoadUint64(&a.bits))
}

func (a *atomicFloat) IncrementBy(incrBy float64) {
	for {
		oldBits := atomic.LoadUint64(&a.bits)
		newBits := math.Float64bits(math.Float64frombits(oldBits) + incrBy)
		if atomic.CompareAndSwapUint64(&a.bits, oldBits, newBits) {
			break
		}
	}
}
