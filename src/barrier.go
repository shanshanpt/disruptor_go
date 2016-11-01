package disruptor

type Barrier interface {
	GetBarrier(int64) int64
}
