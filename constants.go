package bitmap

const (
	bitLength uint32 = 64                   // 一个bucket的长度
	fullMask  uint64 = (1 << bitLength) - 1 // 覆盖bucket所有位的掩码
)

var (
	maxLength = MaxUint64SliceCap() // 平台对应的[]uint64最大长度
)
