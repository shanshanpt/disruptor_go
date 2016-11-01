package disruptor

import "sync/atomic"

// writer & multi producer
type MultiProducer struct {
	// these fields are the same as the Producer struct{}
	WrittenCursor   *Sequence
	GateSequence    *Sequence
	BarrierSequence Barrier
	BufferSize      int64
	// An extra field to record these data that producers have written
	// for example: a producer apply N slots, and writing data process
	// spend M time. During the time, other producers can also apply slots, and
	// write data. And consumers can eat data that have been written.
	//
	// the WrittenCursor can be use by all producers,
	// then, consumers can not eat data by WrittenCursor,
	// so we add PublishedBuffer for consumer to eat data
	PublishedBuffer []int64
}

func NewMultiProducer(written *Sequence, barrier Barrier, bufferSize int64) *MultiProducer {
	// buffer size must be power of 2
	// then we can get the ring buffer index of a sequence by sequence&(size-1)
	AssertPowerOfTwo(bufferSize)

	// Init: -1 represent init value, sequence start from 0
	publishedBuf := make([]int64, bufferSize)
	for i := int64(0); i < bufferSize; i++ {
		publishedBuf[i] = InitialSequenceValue
	}
	return &MultiProducer{
		WrittenCursor:   written,
		GateSequence:    NewSequence(),
		BarrierSequence: barrier,
		BufferSize:      bufferSize,
		PublishedBuffer: publishedBuf,
	}
}

func (this *MultiProducer) SetBarrierSequence(b Barrier) {
	this.BarrierSequence = b
}

// apply for n slots
func (this *MultiProducer) Next(n int64) int64 {
	if n < 1 {
		panic(ErrNotPositiveInteger)
	} else if n > this.BufferSize {
		panic(ErrTooLarge)
	}

	// producer will execute the for circle until CAS success
	for {
		// the below codes are the same as Producer.Next()
		currentCursor := this.WrittenCursor.Get()
		next := currentCursor + n
		wrapPoint := next - this.BufferSize
		gate := this.GateSequence.Get()

		// if the buffer is full, producer should be waiting...
		if wrapPoint > gate || gate > this.WrittenCursor.Get() {
			for ; wrapPoint > gate; gate = this.BarrierSequence.GetBarrier(0) {
				this.GateSequence.Set(gate)
			}
		}

		// ATTENTION: we use CAS instead of lock
		// If CAS successful, the this.WrittenCursor.Value will be set to 'next'
		if atomic.CompareAndSwapInt64(&this.WrittenCursor.Value, currentCursor, next) {
			return next
		}
	}
}

// publish the written, then the consumer can consume it
func (this *MultiProducer) Publish(n int64) {
	// means the nth data have been written successful
	// mask = BufferSize-1
	this.PublishedBuffer[n&(this.BufferSize-1)] = n
}

func (this *MultiProducer) PublishBatch(low int64, high int64) {
	if low > high {
		panic(ErrLowSequenceBiggerThanHigh)
	} else if high-low > this.BufferSize {
		panic(ErrPublishSequencesOutOfBuffer)
	}

	// we should write low~high data one by one from high to low
	for i := high; i >= low; i-- {
		this.PublishedBuffer[i&(this.BufferSize-1)] = i
	}
}
