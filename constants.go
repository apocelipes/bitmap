package bitmap

const (
	bitLength uint32 = 64                   // 一个bucket的长度
	fullMask  uint64 = (1 << bitLength) - 1 // 覆盖bucket所有位的掩码
)
