package bitmap

import (
	"strconv"
	"strings"
)

// Bitmap 每个Bitmap包含了数个容纳uint64类型对应bit数的空间
// 这些空间的索引从0开始
type Bitmap struct {
	buckets []uint64
	// map中包含的位的总数
	length uint32
}

// NewBitmap 创建有length个位的bitmap
// length不能超过math.MaxInt32，因为slice的长度不能超过MaxInt，
// 所以我们选取32位int的最大值以便同时兼容32位和64位系统
func NewBitmap(length uint32) *Bitmap {
	if length > maxLength {
		return nil
	}

	blocks := calcBlocks(length)
	return &Bitmap{
		buckets: make([]uint64, blocks),
		length:  length,
	}
}

// Len 返回Bitmap当前含有buckets的数量
func (b *Bitmap) Len() uint32 {
	return b.length
}

// SetOne 将指定索引的位设置为1，索引从0开始
func (b *Bitmap) SetOne(pos uint32) error {
	bucketIndex, bitIndex, err := getPos(b.Len(), pos)
	if err != nil {
		return err
	}

	b.buckets[bucketIndex] |= 1 << bitIndex
	return nil
}

// SetZero 将指定位置的位设置为0，索引从0开始
func (b *Bitmap) SetZero(pos uint32) error {
	bucketIndex, bitIndex, err := getPos(b.Len(), pos)
	if err != nil {
		return err
	}

	b.buckets[bucketIndex] &= ^(1 << bitIndex)
	return nil
}

// IsOne 检查索引指定的位是否是1，索引从0开始
func (b *Bitmap) IsOne(pos uint32) (bool, error) {
	bucketIndex, bitIndex, err := getPos(b.Len(), pos)
	if err != nil {
		return false, err
	}

	return (b.buckets[bucketIndex]>>bitIndex)&1 == 1, nil
}

// Flip 反转索引指定的位的值，1变为0，0变为1，索引从0开始
func (b *Bitmap) Flip(pos uint32) error {
	bucketIndex, bitIndex, err := getPos(b.Len(), pos)
	if err != nil {
		return err
	}

	b.buckets[bucketIndex] ^= 1 << bitIndex
	return nil
}

// ClearAll 将所有位清零
func (b *Bitmap) ClearAll() {
	for i := range b.buckets {
		b.buckets[i] = 0
	}
}

// FillAll 将所有位设置为1
func (b *Bitmap) FillAll() {
	for i := range b.buckets {
		b.buckets[i] = fullMask
	}
}

// String 将每一位的数据按照buckets的组织顺序打印
func (b *Bitmap) String() string {
	buff := strings.Builder{}
	length := b.Len()
	index := 0
	for length != 0 {
		data := ""
		if length > bitLength {
			data = strconv.FormatUint(uint64(b.buckets[index]), 2)
			data = paddingLeftZero(data, int(bitLength))
			length -= bitLength
		} else {
			// 该bucket中没有使用全部的位，过滤出已经使用的位
			mask := fullMask
			mask >>= bitLength - length
			data = strconv.FormatUint(b.buckets[index]&mask, 2)
			data = paddingLeftZero(data, int(length))
			length = 0
		}
		buff.WriteString(data)
		index++
	}

	return buff.String()
}
