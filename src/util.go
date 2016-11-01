package disruptor

import "errors"

const (
	MaxSequenceValue     int64 = (1 << 63) - 1
	InitialSequenceValue int64 = -1
	// solve the false sharing problem
	CpuCacheLinePaddingCount int = 7
)

var (
	ErrTooLarge                    = errors.New("Apply too many slots")
	ErrNotPositiveInteger          = errors.New("Sequece must be a positive integer value")
	ErrNotPowerOfTwo               = errors.New("Value must be power of 2")
	ErrNoConsumer                  = errors.New("No consumers!")
	ErrLowSequenceBiggerThanHigh   = errors.New("Low sequence is bigger than the high one")
	ErrPublishSequencesOutOfBuffer = errors.New("Too many publish sequences!")
)

// if value is power of 2, then value&(value-1)=0
// example, 8 binary cod is 1000, 8-1=7 binary cod is 0111
// so: 1000&0111=0
func AssertPowerOfTwo(value int64) {
	if value > 0 && (value&(value-1)) != 0 {
		panic(ErrNotPowerOfTwo)
	}
}
