package bitmap

import (
	"math"
	"reflect"
	"runtime"
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
