package disruptor

type Controller struct {
	BufferSize           int64
	ConsumerHandlerGroup [][]ConsumerHandler
}

// new a controller
func NewController(bufferSize int64) *Controller {
	AssertPowerOfTwo(bufferSize)
	return &Controller{
		BufferSize: bufferSize,
	}
}

// add a group consumers
// the previous group consumers must consume data before the later group
// so, the later's barriers are the previous group consumers' sequences
func (this *Controller) AddConsumersGroup(c ...ConsumerHandler) *Controller {
	if len(c) > 0 {
		newGroup := make([]ConsumerHandler, len(c))
		copy(newGroup, c)
		// Add a new group
		this.ConsumerHandlerGroup = append(this.ConsumerHandlerGroup, newGroup)
	}
	return this
}

// build a disruptor
func (this *Controller) BuildDisruptor() *Disruptor {
	if len(this.ConsumerHandlerGroup) < 1 || len(this.ConsumerHandlerGroup[0]) < 1 {
		panic(ErrNoConsumer)
	}

	// generate consumer
	writtenCursor := NewSequence()
	consumerBarriers := GeneralBarrier{}
	consumerBarriers.BarrierSequence = append(consumerBarriers.BarrierSequence, writtenCursor)
	allConsumers := []*Consumer{}
	// build all consumers
	for _, g := range this.ConsumerHandlerGroup {
		preConsumerSequence := GeneralBarrier{}
		for _, handler := range g {
			sequence := NewSequence()
			c := NewConsumer(sequence, consumerBarriers, handler)
			allConsumers = append(allConsumers, c)
			preConsumerSequence.BarrierSequence = append(preConsumerSequence.BarrierSequence, sequence)
		}
		consumerBarriers = preConsumerSequence
	}

	// generate producer
	p := NewProducer(writtenCursor, consumerBarriers, this.BufferSize)

	return &Disruptor{producer: p, consumers: allConsumers}
}

// build multi producers' disruptor
func (this *Controller) BuildMultiDisruptor() *MultiDisruptor {
	if len(this.ConsumerHandlerGroup) < 1 || len(this.ConsumerHandlerGroup[0]) < 1 {
		panic(ErrNoConsumer)
	}

	// generate producer
	writtenCursor := NewSequence()
	p := NewMultiProducer(writtenCursor, nil, this.BufferSize)

	// generate consumer
	// we use MultiProducerBarrier here
	// producer barrier just effect for first level consumers
	// the later level consumers depends on the pre-level consumers
	// so we still use GeneralBarrier
	mpBarriers := MultiProducerBarrier{}
	mpBarriers.BufferSizeMask = this.BufferSize - 1
	mpBarriers.WrittenCursor = writtenCursor
	mpBarriers.PublishedSequences = p.PublishedBuffer
	var consumerBarriers Barrier = mpBarriers

	allConsumers := []*Consumer{}
	// build all consumers
	for _, g := range this.ConsumerHandlerGroup {
		// others, we still use GeneralBarrier
		preConsumerSequence := GeneralBarrier{}
		for _, handler := range g {
			sequence := NewSequence()
			c := NewConsumer(sequence, consumerBarriers, handler)
			allConsumers = append(allConsumers, c)
			preConsumerSequence.BarrierSequence = append(preConsumerSequence.BarrierSequence, sequence)
		}
		consumerBarriers = preConsumerSequence
	}

	// Attention: we should set producers' BarrierSequence as consumerBarriers
	p.SetBarrierSequence(consumerBarriers)

	return &MultiDisruptor{producers: p, consumers: allConsumers}
}
