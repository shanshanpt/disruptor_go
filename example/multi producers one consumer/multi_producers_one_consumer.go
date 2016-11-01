package main

import (
	"../../src"
	"fmt"
	"runtime"
	"time"
)

const (
	// 默认buf大小
	BufferSize   = 1024 * 64
	BufferMask   = BufferSize - 1
	Iterations   = 1000000 * 100
	Reservations = 128
)

/*
const (
	// 默认buf大小
	BufferSize   = 16
	BufferMask   = BufferSize - 1
	Iterations   = 100
	Reservations = 4
)*/

var ringBuffer = [BufferSize]int64{}

func main() {
	for i := 0; i < 100; i++ {
		ringBuffer = [BufferSize]int64{}

		// one consumer
		controller := disruptor.
			NewController(BufferSize).
			AddConsumersGroup(AConsumer{}).
			BuildMultiDisruptor()

		controller.Start()

		started := time.Now()

		// multi producers

		// producers count, define as many as you like
		producersCount := 2
		// chan for nice close of producers
		ch := make(chan int)
		for i := 0; i < producersCount; i++ {
			// producer i
			go func(id int, ch chan int) {
				fmt.Println(id, "  start")
				publish(id, controller.GetProducer())
				fmt.Println(id, "  end")
				ch <- 1
			}(i, ch)
		}

		// main goroutine should waiting..
		goroutineStopCount := 0
		for {
			select {
			case <-ch:
				goroutineStopCount++
				fmt.Println("end a goroutine ", goroutineStopCount)
				if goroutineStopCount >= producersCount {
					goto LABEL
				}
				runtime.Gosched()
			default:
				runtime.Gosched()
			}
		}

	LABEL:
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
			panic("Error")
		}
		//fmt.Println("CONSUME: -> ", ringBuffer[lower&BufferMask])
		lower++
	}

}

// what the producer do
//
func publish(id int, writer *disruptor.MultiProducer) {
	sequence := disruptor.InitialSequenceValue
	for sequence <= Iterations {
		sequence = writer.Next(Reservations)
		for lower := sequence - Reservations + 1; lower <= sequence; lower++ {
			ringBuffer[lower&BufferMask] = lower
			//fmt.Println("PRODUCE: ", id, "   ",ringBuffer[lower&BufferMask])
		}

		writer.PublishBatch(sequence-Reservations+1, sequence)
	}
}
