package main

import (
	"../../src"
	"fmt"
	"time"
)

const (
	// 默认buf大小
	BufferSize   = 1024 * 64
	BufferMask   = BufferSize - 1
	Iterations   = 1000000 * 100
	Reservations = 128
)

var ringBuffer = [BufferSize]int64{}

func main() {
	for i := 0; i < 100; i++ {
		ringBuffer = [BufferSize]int64{}

		// four(multi) consumers
		controller := disruptor.
			NewController(BufferSize).
			AddConsumersGroup(AConsumer{}, AConsumer{}, AConsumer{}, AConsumer{}).
			BuildDisruptor()

		controller.Start()

		started := time.Now()
		// 生产者
		publish(controller.GetProducer())
		finished := time.Now()

		// wait the consumer eat all data
		time.Sleep(time.Second)
		controller.Stop()
		fmt.Println(Iterations, finished.Sub(started))
	}
}

// we should define a consumer to eat data,
// the Consume function must be overwritten
//
type AConsumer struct{}

func (this AConsumer) Consume(lower int64, upper int64) {
	for lower <= upper {
		l := ringBuffer[lower&BufferMask]
		if l != lower {
			panic("Error!!!")
		}
		//fmt.Println("CONSUME: -> ", ringBuffer[lower&BufferMask])
		lower++
	}

}

// what the producer do
//
func publish(writer *disruptor.Producer) {
	sequence := disruptor.InitialSequenceValue
	for sequence <= Iterations {
		sequence = writer.Next(Reservations)
		for lower := sequence - Reservations + 1; lower <= sequence; lower++ {
			ringBuffer[lower&BufferMask] = lower
			//fmt.Println("PRODUCE: ", ringBuffer[lower&BufferMask])
		}

		writer.PublishBatch(sequence-Reservations+1, sequence)
	}
}
