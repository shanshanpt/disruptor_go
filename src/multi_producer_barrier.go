package disruptor

// MultiProducerBarrier depends on producers' shared write cursor
// and the published sequences by producers
// the struct just for consumers in multi-producers scene
type MultiProducerBarrier struct {
	WrittenCursor      *Sequence
	BufferSizeMask     int64
	PublishedSequences []int64
}

func (this MultiProducerBarrier) GetBarrier(lowSequence int64) int64 {
	high := this.WrittenCursor.Get()
	for ; lowSequence <= high; lowSequence++ {
		// sequence have not been written yet, just return those sequence less than lowSequence
		if lowSequence != this.PublishedSequences[lowSequence&this.BufferSizeMask] {
			return lowSequence - 1
		}
	}
	return high
}
