package disruptor

import (
	"runtime"
	//"time"
	"fmt"
)

// reader & consumer
type Consumer struct {
	// ReadSequence must be a pointer type
	// case, it will be used by producer for barrier
	ReadSequence *Sequence
	// you know, barrier may be the producer's writer cursor
	// in addition, if there are two group consumers, A and B,
	// B should consume after A, so B's BarrierSequence is
	// all consumers' read sequence in A
	BarrierSequence Barrier
	// consume handler
	ReadHandler ConsumerHandler
	// running state
	IsRunning bool
}

func NewConsumer(read *Sequence, barrier Barrier, handler ConsumerHandler) *Consumer {
	return &Consumer{
		ReadSequence:    read,
		BarrierSequence: barrier,
		ReadHandler:     handler,
		IsRunning:       false,
	}
}

// consumer can consume as many as it can
// batch consume
func (this *Consumer) consuming() {
	// record the consumed sequence
	haveConsumed := this.ReadSequence.Get()

	for this.IsRunning {
		next := haveConsumed + 1

		// Attention:
		// this.BarrierSequence may include the producer's cursor and
		// the pre-group consumer's sequences
		high := this.BarrierSequence.GetBarrier(next)

		// have data to consume
		if next <= high {
			// consume
			this.ReadHandler.Consume(next, high)
			// change the read sequence
			this.ReadSequence.Set(high)
			haveConsumed = high
		}

		runtime.Gosched()
		//time.Sleep(time.Millisecond)
	}
}

// start to consume
func (this *Consumer) Start() {
	fmt.Println("start...")
	this.IsRunning = true
	go this.consuming()
}

func (this *Consumer) Stop() {
	this.IsRunning = false
}
