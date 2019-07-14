package bitmap

import (
	"errors"
	"strconv"
	"strings"
)

type Bitmap struct {
	buckets []uint64
	// map中包含的位的总数
	length uint32
}

const (
	bitLength uint32 = 64                   // 一个bucket的长度
	fullMask  uint64 = (1 << bitLength) - 1 // 覆盖bucket所有位的掩码
)

var (
	maxLength = MaxUint64SliceCap() // 平台对应的[]uint64最大长度
)

// NewBitmap 创建有length个位的bitmap
// length不能超过math.MaxInt32，因为slice的长度不能超过MaxInt，
// 所以我们选取32位int的最大值以便同时兼容32位和64位系统
func NewBitmap(length uint32) *Bitmap {
	if length > maxLength {
		return nil
	}

	blocks := length / bitLength
	remainder := length % bitLength
	if remainder != 0 {
		blocks++
	}
	return &Bitmap{
		buckets: make([]uint64, blocks),
		length:  length,
	}
}

func (b *Bitmap) Len() uint32 {
	return b.length
}

// 返回指定索引对应的bucket的index，以及索引在该bucket中的位置
// bucket中位的排列采用大端序，最右边的一位代表bitmap的索引在该bucket中最小。
// 例如对于长度为bitLength的bitmap，buckets[0]的最右边一位代表了bitmap的索引0表示的位
// buckets[0]最左边一位代表了bitmap的索引bitLength-1所表示的位
func (b *Bitmap) getPos(pos uint32) (uint32, uint32, error) {
	if pos >= b.Len() || pos > maxLength {
		return 0, 0, errors.New("out of Bitmap length")
	}

	bucketIndex := pos / bitLength // 确定bucket的索引
	bitIndex := pos % bitLength    // 确定在bucket中的位置

	return bucketIndex, bitIndex, nil
}

// SetOne 将指定索引的位设置为1，索引从0开始
func (b *Bitmap) SetOne(pos uint32) error {
	bucketIndex, bitIndex, err := b.getPos(pos)
	if err != nil {
		return err
	}

	b.buckets[bucketIndex] |= 1 << bitIndex
	return nil
}

// SetZero 将指定位置的位设置为0，索引从0开始
func (b *Bitmap) SetZero(pos uint32) error {
	bucketIndex, bitIndex, err := b.getPos(pos)
	if err != nil {
		return err
	}

	b.buckets[bucketIndex] &= ^(1 << bitIndex)
	return nil
}

// IsOne 检查索引指定的位是否是1，索引从0开始
func (b *Bitmap) IsOne(pos uint32) (bool, error) {
	bucketIndex, bitIndex, err := b.getPos(pos)
	if err != nil {
		return false, err
	}

	return (b.buckets[bucketIndex]>>bitIndex)&1 == 1, nil
}

// Flip 反转索引指定的位的值，1变为0，0变为1，索引从0开始
func (b *Bitmap) Flip(pos uint32) error {
	bucketIndex, bitIndex, err := b.getPos(pos)
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

// 向数字字符串的左侧填充0，使字符串的长度达到fullSize
func paddingLeftZero(data string, fullSize int) string {
	if len(data) >= fullSize {
		return data
	}

	zeros := fullSize - len(data)
	return strings.Repeat("0", zeros) + data
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
