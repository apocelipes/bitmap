package bitmap

import (
	"math"
	"reflect"
	"runtime"
	"strings"
)

// maxSliceCap 计算元素类型和i相同的slice最大可以容纳多少元素。
// 不同平台slice的最大长度不同，
// slice的最大长度是用elements*elementSize计算的，
// 不同类型的元素占用的空间不同，导致了slice最大可容纳元素数的不同
func maxSliceCap(i interface{}) int {
	_64bit := uintptr(1 << (^uintptr(0) >> 63) / 2)

	var goosWindows, goosDarwin, goarchArm64 uintptr
	switch runtime.GOOS {
	case "darwin":
		goosDarwin = 1
	case "windows":
		goosWindows = 1
	}

	switch runtime.GOARCH {
	case "arm64":
		goarchArm64 = 1
	}

	heapMapBits := (_64bit*goosWindows)*35 + (_64bit*(1-goosWindows)*(1-goosDarwin*goarchArm64))*39 + goosDarwin*goarchArm64*31 + (1-_64bit)*32
	maxMem := uintptr(1<<heapMapBits - 1)

	elemSize := reflect.ValueOf(i).Type().Size()
	max := maxMem / elemSize

	if int(max) < 0 {
		return 1<<31 - 1
	}

	return int(max)
}

// MaxUint64SliceCap 将[]uint64的最大长度限制在MaxInt32，
// 因为Linux/Darwin amd64上允许的最大长度大于MaxInt32，
// Windows amd64等于这个值，Darwin arm64上小于该值，
// 因此超过MaxInt32的值默认返回MaxInt32
func MaxUint64SliceCap() uint32 {
	max := int64(maxSliceCap(uint64(0)))
	if max > math.MaxInt32 {
		return math.MaxInt32
	}

	return uint32(max)
}

// 向数字字符串的左侧填充0，使字符串的长度达到fullSize
func paddingLeftZero(data string, fullSize int) string {
	if len(data) >= fullSize {
		return data
	}

	zeros := fullSize - len(data)
	return strings.Repeat("0", zeros) + data
}

// 返回指定索引对应的bucket的index，以及索引在该bucket中的位置
// bucket中位的排列采用大端序，最右边的一位代表bitmap的索引在该bucket中最小。
// 例如对于长度为bitLength的bitmap，buckets[0]的最右边一位代表了bitmap的索引0表示的位
// buckets[0]最左边一位代表了bitmap的索引bitLength-1所表示的位
func getPos(mapLength, pos uint32) (uint32, uint32, error) {
	if pos >= mapLength || pos > maxLength {
		return 0, 0, errOutOfLength
	}

	bucketIndex := pos / bitLength // 确定bucket的索引
	bitIndex := pos % bitLength    // 确定在bucket中的位置

	return bucketIndex, bitIndex, nil
}

// calcBlocks 根据需要的位数计算bitmap中包含的bucket的个数
func calcBlocks(length uint32) uint32 {
	blocks := length / bitLength
	remainder := length % bitLength
	if remainder != 0 {
		blocks++
	}

	return blocks
}
