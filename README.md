###LMAX Disruptor
A High Performance Inter-Thread Messaging Library

###Introduction
<a href="https://github.com/LMAX-Exchange/disruptor/wiki/Introduction"  target="_blank">disruptor wiki</a> <br>

###Examples
What we should do: <br>
1. Define our ringBuffer[], BufferSize and so on <br>
2. Define a consumer class to eat data, we should implement function 'Consume' <br>
3. New a 'controller', and then, running the consumers and producers. <br>
4. Everything is done! <br>
#####1. one producer and one consumer
'''
// 1: main
// define our ringBuffer
ringBuffer = [BufferSize]int64{}

// create a controller
// NewController(BufferSize): Define buffer size. Attention: size should be power of 2.
// AddConsumersGroup(AConsumer{}): add a consumer, we should implement class AConsumer{} by ourselves
// BuildDisruptor(): build one producer process
controller := disruptor.
    NewController(BufferSize).
	AddConsumersGroup(AConsumer{}).
	BuildDisruptor()

// consumers start
controller.Start()

// producers start
publish(controller.GetProducer())


// 2: consumer
type AConsumer struct{}

func (this AConsumer) Consume(lower int64, upper int64) {
	// consume data for ringBuffer
	for lower <= upper {
		message := ringBuffer[lower&BufferMask]
		if message != lower {
			panic("Error!!!")
		}
		lower++
	}
}


// 3: producer
func publish(id int, writer *disruptor.MultiProducer) {
	sequence := disruptor.InitialSequenceValue
	for sequence <= Iterations {
		// producer apply for Reservations slots to write data
		sequence = writer.Next(Reservations)
		// write data
		for lower := sequence - Reservations + 1; lower <= sequence; lower++ {
			ringBuffer[lower&BufferMask] = lower
		}
		// Publish these data, then the data can be used by consumers
		writer.PublishBatch(sequence-Reservations+1, sequence)
	}
}
'''

Please find the complete implementation in example folder. <br>
Also, there are some other examples:  one producer and multi consumers,
multi producers and one consumer, multi producers and multi consumers


###Reference
The project's implement is imitate from
<a href="https://github.com/LMAX-Exchange/disruptor"  target="_blank">disruptor 1</a> and
<a href="https://github.com/smartystreets/go-disruptor"  target="_blank">disruptor 2</a>



