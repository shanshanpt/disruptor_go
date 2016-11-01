package disruptor

// sequence struct
// padding for avoid false sharing (just for general cpu, 64 Bytes cpu cache line)
type Sequence struct {
	Value   int64
	padding [CpuCacheLinePaddingCount]int64
}

func NewSequence() *Sequence {
	return &Sequence{
		Value: InitialSequenceValue,
	}
}

func (this *Sequence) Get() int64 {
	return this.Value
}

func (this *Sequence) Set(sequence int64) {
	this.Value = sequence
}
