package disruptor

// writer & producer
type Producer struct {
	// WrittenCursor must be a pointer type
	// case, it will be used by consumer for barrier
	WrittenCursor *Sequence
	// record the min sequence that consumer have eat
	GateSequence *Sequence
	// all consumers' current read sequence
	BarrierSequence Barrier
	BufferSize      int64
}

func NewProducer(written *Sequence, barrier Barrier, bufferSize int64) *Producer {
	// buffer size must be power of 2
	// then we can get the ring buffer index of a sequence by sequence&(size-1)
	AssertPowerOfTwo(bufferSize)

	return &Producer{
		WrittenCursor:   written,
		GateSequence:    NewSequence(),
		BarrierSequence: barrier,
		BufferSize:      bufferSize,
	}
}

// apply for n slots
func (this *Producer) Next(n int64) int64 {
	if n < 1 {
		panic(ErrNotPositiveInteger)
	} else if n > this.BufferSize {
		panic(ErrTooLarge)
	}

	next := this.WrittenCursor.Get() + n
	wrapPoint := next - this.BufferSize
	gate := this.GateSequence.Get()

	// if the buffer is full, producer should be waiting...
	if wrapPoint > gate || gate > this.WrittenCursor.Get() {
		for ; wrapPoint > gate; gate = this.BarrierSequence.GetBarrier(0) {
			this.GateSequence.Set(gate)
		}
	}

	return next
}

// publish the written, then the consumer can consume it
func (this *Producer) Publish(n int64) {
	this.WrittenCursor.Set(n)
}

func (this *Producer) PublishBatch(low int64, high int64) {
	if low > high {
		panic(ErrLowSequenceBiggerThanHigh)
	} else if high-low > this.BufferSize {
		panic(ErrPublishSequencesOutOfBuffer)
	}
	this.Publish(high)
}
