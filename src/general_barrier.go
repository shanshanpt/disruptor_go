package disruptor

// GeneralBarrier depends on the BarrierSequence Param
// the BarrierSequence Param can be producers' current write cursors
// or consumers' current read cursors
type GeneralBarrier struct {
	BarrierSequence []*Sequence
}

func (this GeneralBarrier) GetBarrier(int64) int64 {
	min := MaxSequenceValue

	for _, b := range this.BarrierSequence {
		if min > b.Get() {
			min = b.Get()
		}
	}
	return min
}
