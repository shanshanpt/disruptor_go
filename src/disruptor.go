package disruptor

// One producer & multiple consumer
type Disruptor struct {
	producer  *Producer
	consumers []*Consumer
}

// Multiple producer & multiple consumer
type MultiDisruptor struct {
	producers *MultiProducer
	consumers []*Consumer
}

func (this *Disruptor) Start() {
	for _, c := range this.consumers {
		c.Start()
	}
}

func (this *Disruptor) Stop() {
	for _, c := range this.consumers {
		c.Stop()
	}
}

func (this *Disruptor) GetProducer() *Producer {
	return this.producer
}

func (this *Disruptor) GetConsumer() []*Consumer {
	return this.consumers
}

func (this *MultiDisruptor) Start() {
	for _, c := range this.consumers {
		c.Start()
	}
}

func (this *MultiDisruptor) Stop() {
	for _, c := range this.consumers {
		c.Stop()
	}
}

func (this *MultiDisruptor) GetProducer() *MultiProducer {
	return this.producers
}

func (this *MultiDisruptor) GetConsumer() []*Consumer {
	return this.consumers
}
