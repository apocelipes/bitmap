package bitmap

import (
	"errors"
	"math"
	"strconv"
	"strings"
)

type Bitmap struct {
	buckets []uint64
	length  uint32
}

const (
	bitLength uint32 = 64
	fullMask uint64 = (1 << bitLength) - 1
)

func NewBitmap(length uint32) *Bitmap {
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

func (b *Bitmap) getPos(pos uint32) (uint32, uint32, error) {
	if pos >= b.Len() || pos > uint32(math.MaxInt32) {
		return 0, 0, errors.New("out of Bitmap length")
	}

	bucketIndex := pos / bitLength
	bitIndex := pos % bitLength

	return bucketIndex, bitIndex, nil
}

func (b *Bitmap) SetOne(pos uint32) error {
	bucketIndex, bitIndex, err := b.getPos(pos)
	if err != nil {
		return err
	}

	b.buckets[bucketIndex] |= 1 << bitIndex
	return nil
}

func (b *Bitmap) SetZero(pos uint32) error {
	bucketIndex, bitIndex, err := b.getPos(pos)
	if err != nil {
		return err
	}

	b.buckets[bucketIndex] &= ^(1 << bitIndex)
	return nil
}

func (b *Bitmap) IsOne(pos uint32) (bool, error) {
	bucketIndex, bitIndex, err := b.getPos(pos)
	if err != nil {
		return false, err
	}

	return (b.buckets[bucketIndex]>>bitIndex)&1 == 1, nil
}

func (b *Bitmap) Flip(pos uint32) error {
	bucketIndex, bitIndex, err := b.getPos(pos)
	if err != nil {
		return err
	}

	b.buckets[bucketIndex] ^= 1 << bitIndex
	return nil
}

func (b *Bitmap) ClearAll() {
	for i := range b.buckets {
		b.buckets[i] = 0
	}
}

func (b *Bitmap) FillAll() {
	for i := range b.buckets {
		b.buckets[i] = fullMask
	}
}

func paddingLeftZero(data string, fullSize int) string {
	if len(data) >= fullSize {
		return data
	}

	zeros := fullSize - len(data)
	return strings.Repeat("0", zeros) + data
}

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
