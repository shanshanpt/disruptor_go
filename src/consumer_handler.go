package disruptor

type ConsumerHandler interface {
	Consume(lower int64, upper int64)
}
