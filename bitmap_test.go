package bitmap

import (
	"fmt"
	"math"
	"math/rand"
	"testing"
	"time"
)

func genRandPos(size uint32) uint32 {
	return rand.Uint32() % size
}

func genBitmap() (*Bitmap, uint32) {
	length := rand.Uint32() % uint32(math.MaxInt32)
	bm := NewBitmap(length)
	return bm, length
}

func genBitmapPos() (*Bitmap, uint32) {
	rand.Seed(time.Now().UnixNano())
	bm, length := genBitmap()
	return bm, genRandPos(length)
}

func TestBitmap_Len(t *testing.T) {
	for i := 0; i < 100; i++ {
		bm, length := genBitmap()
		if bm.Len() != length {
			t.Errorf("length error: want: %v\n\thave: %v\n", length, bm.Len())
		}
	}
}

func TestBitmap_WrongLen(t *testing.T) {
	bm := NewBitmap(math.MaxUint32)
	if bm != nil {
		t.Error("NewBitmap out of length but did not return a nil")
	}
}

func TestBitmap_SetOne(t *testing.T) {
	for i := 0; i < 100; i++ {
		bm, pos := genBitmapPos()
		err := bm.SetOne(pos)
		if err != nil {
			t.Error(err)
		}
		ok, err := bm.IsOne(pos)
		if err != nil {
			t.Error(err)
		}
		if !ok {
			t.Errorf("pos: %v should be one\n", pos)
		}
	}
}

func TestBitmap_PosCheck(t *testing.T) {
	testData := []struct {
		length uint32
		pos    []uint32
		failed []bool
	}{
		{
			length: 10,
			pos:    []uint32{0, 9, 10, 20},
			failed: []bool{false, false, true, true},
		},
		{
			length: 0,
			pos:    []uint32{0},
			failed: []bool{true},
		},
		{
			length: 33,
			pos:    []uint32{0, 31, 32, 33, 35},
			failed: []bool{false, false, false, true, true},
		},
	}

	for _, v := range testData {
		bm := NewBitmap(v.length)
		for i := range v.pos {
			err := bm.SetOne(v.pos[i])
			if (err != nil) != v.failed[i] {
				t.Errorf("SetOne check pos [%v] failed\n", v.pos[i])
			}
			err = bm.SetZero(v.pos[i])
			if (err != nil) != v.failed[i] {
				t.Errorf("SetZero check pos [%v] failed\n", v.pos[i])
			}
			_, err = bm.IsOne(v.pos[i])
			if (err != nil) != v.failed[i] {
				t.Errorf("IsOne check pos [%v] failed\n", v.pos[i])
			}
			err = bm.Flip(v.pos[i])
			if (err != nil) != v.failed[i] {
				t.Errorf("Flip check pos [%v] failed\n", v.pos[i])
			}
		}
	}
}

func TestBitmap_SetZero(t *testing.T) {
	for i := 0; i < 100; i++ {
		bm, pos := genBitmapPos()
		err := bm.SetZero(pos)
		if err != nil {
			t.Error(err)
		}
		ok, err := bm.IsOne(pos)
		if err != nil {
			t.Error(err)
		}
		if ok {
			t.Errorf("pos: %v should be zero\n", pos)
		}
	}
}

func TestBitmap_Flip(t *testing.T) {
	for i := 0; i < 100; i++ {
		bm, pos := genBitmapPos()
		err := bm.SetOne(pos)
		if err != nil {
			t.Error(err)
		}
		ok, err := bm.IsOne(pos)
		if err != nil {
			t.Error(err)
		} else if !ok {
			t.Error("set one failed")
		}

		err = bm.Flip(pos)
		if err != nil {
			t.Error(err)
		}

		ok, _ = bm.IsOne(pos)
		if ok {
			t.Error("flip failed")
		}
	}
}

func TestBitmap_FillClearAll(t *testing.T) {
	lengths := []uint32{1, bitLength - 1, bitLength, bitLength + 1}
	for _, length := range lengths {
		bm := NewBitmap(length)
		bm.FillAll()
		for i := uint32(0); i < length; i++ {
			ok, err := bm.IsOne(i)
			if err != nil {
				t.Error(err)
			}
			if !ok {
				t.Errorf("FillAll failed at [%v]\n", i)
			}
		}

		bm.ClearAll()
		for i := uint32(0); i < length; i++ {
			ok, err := bm.IsOne(i)
			if err != nil {
				t.Error(err)
			}
			if ok {
				t.Errorf("ClearAll failed at [%v]\n", i)
			}
		}
	}
}

func ExampleBitmap_String() {
	bm := NewBitmap(65)
	for i := uint32(1); i < 33; i++ {
		_ = bm.SetOne(i)
	}
	fmt.Println(bm)
	// Output:
	// 00000000000000000000000000000001111111111111111111111111111111100
}

func TestBitmap_String(t *testing.T) {
	bm := NewBitmap(65)
	for i := uint32(1); i < 33; i++ {
		_ = bm.SetOne(i)
	}
	//_ = bm.Flip(33)
	//err := bm.Flip(0)
	//t.Log(err)
	t.Log(len(bm.String()))
	bm.FillAll()
	t.Log(bm.String(), bm.buckets[0])
}

func bitSort(arr []uint32) {
	var max uint32
	for i := range arr {
		if max < arr[i] {
			max = arr[i]
		}
	}

	bm := NewBitmap(max + 1)
	for i := range arr {
		err := bm.SetOne(arr[i])
		if err != nil {
			panic(err)
		}
	}

	counter := 0
	for i := uint32(0); i <= max && counter < len(arr); i++ {
		if ok, _ := bm.IsOne(i); ok {
			arr[counter] = i
			counter++
		}
	}
}

func testEq(a, b []uint32) bool {

	// If one is nil, the other must also be nil.
	if (a == nil) != (b == nil) {
		return false
	}

	if len(a) != len(b) {
		return false
	}

	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}

	return true
}

func TestBitmapSort(t *testing.T) {
	testData := []struct {
		src []uint32
		res []uint32
	}{
		{
			src: []uint32{1, 2, 4, 7, 0, 6, 5},
			res: []uint32{0, 1, 2, 4, 5, 6, 7},
		},
		{
			src: []uint32{11, 3, 0, 6, 5},
			res: []uint32{0, 3, 5, 6, 11},
		},
		{
			src: []uint32{1000},
			res: []uint32{1000},
		},
		{
			src: []uint32{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
			res: []uint32{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
		},
	}

	for _, v := range testData {
		bitSort(v.src)
		if !testEq(v.src, v.res) {
			t.Error("not equal", v.src, v.res)
		}
	}
}

// benchmarks
func benchmarkClear(b *testing.B, size uint32) {
	bm := NewBitmap(size)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		bm.ClearAll()
	}
}

func BenchmarkBitmap_ClearAll_bitLength(b *testing.B) {
	benchmarkClear(b, bitLength)
}

func BenchmarkBitmap_ClearAll_1K(b *testing.B) {
	benchmarkClear(b, 1000)
}

func BenchmarkBitmap_ClearAll_10K(b *testing.B) {
	benchmarkClear(b, 10000)
}

func benchmarkFill(b *testing.B, size uint32) {
	bm := NewBitmap(size)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		bm.FillAll()
	}
}

func BenchmarkBitmap_FillAll_bitLength(b *testing.B) {
	benchmarkFill(b, bitLength)
}

func BenchmarkBitmap_FillAll_1K(b *testing.B) {
	benchmarkFill(b, 1000)
}

func BenchmarkBitmap_FillAll_10K(b *testing.B) {
	benchmarkFill(b, 10000)
}

func benchmarkFlip(b *testing.B, size uint32) {
	bm := NewBitmap(size)
	rand.Seed(time.Now().UnixNano())
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		pos := genRandPos(size)
		_ = bm.Flip(pos)
	}
}

func BenchmarkBitmap_Flip_bitLength(b *testing.B) {
	benchmarkFlip(b, bitLength)
}

func BenchmarkBitmap_Flip_1k(b *testing.B) {
	benchmarkFlip(b, 1000)
}


func BenchmarkBitmap_Flip_10k(b *testing.B) {
	benchmarkFlip(b, 10000)
}

type fillFlag int

const (
	allFill fillFlag = iota
	noFill
	halfFill
)

func dealWithFillFlag(bm *Bitmap, flag fillFlag) {
	switch flag {
	case allFill:
		bm.FillAll()
	case halfFill:
		for i := uint32(0); i < bm.Len()/2; i++ {
			_ = bm.SetOne(i)
		}
	case noFill:
		fallthrough
	default:
		// nop
	}
}

func benchmarkIsOne(b *testing.B, size uint32, flag fillFlag) {
	bm := NewBitmap(size)
	dealWithFillFlag(bm, flag)
	rand.Seed(time.Now().UnixNano())
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		pos := genRandPos(size)
		_, _ = bm.IsOne(pos)
	}
}

func BenchmarkBitmap_IsOne_AllFill_bitLength(b *testing.B) {
	benchmarkIsOne(b, bitLength, allFill)
}

func BenchmarkBitmap_IsOne_AllFill_1K(b *testing.B) {
	benchmarkIsOne(b, 1000, allFill)
}

func BenchmarkBitmap_IsOne_AllFill_10K(b *testing.B) {
	benchmarkIsOne(b, 10000, allFill)
}

func BenchmarkBitmap_IsOne_HalfFill_bitLength(b *testing.B) {
	benchmarkIsOne(b, bitLength, halfFill)
}

func BenchmarkBitmap_IsOne_HalfFill_1K(b *testing.B) {
	benchmarkIsOne(b, 1000, halfFill)
}


func BenchmarkBitmap_IsOne_HalfFill_10K(b *testing.B) {
	benchmarkIsOne(b, 10000, halfFill)
}

func BenchmarkBitmap_IsOne_NoFill_bitLength(b *testing.B) {
	benchmarkIsOne(b, bitLength, noFill)
}

func BenchmarkBitmap_IsOne_NoFill_1K(b *testing.B) {
	benchmarkIsOne(b, 1000, noFill)
}


func BenchmarkBitmap_IsOne_NoFill_10K(b *testing.B) {
	benchmarkIsOne(b, 10000, noFill)
}

func benchmarkSet(b *testing.B, size uint32, value int) {
	bm := NewBitmap(size)
	var f func(uint32) error
	switch value {
	case 1:
		f = bm.SetOne
	case 0:
		f = bm.SetZero
	default:
		return
	}
	rand.Seed(time.Now().UnixNano())
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		pos := genRandPos(size)
		_ = f(pos)
	}
}

func BenchmarkBitmap_SetOne_bitLength(b *testing.B) {
	benchmarkSet(b, bitLength, 1)
}

func BenchmarkBitmap_SetOne_1K(b *testing.B) {
	benchmarkSet(b, 1000, 1)
}

func BenchmarkBitmap_SetOne_10K(b *testing.B) {
	benchmarkSet(b, 10000, 1)
}

func BenchmarkBitmap_SetZero_bitLength(b *testing.B) {
	benchmarkSet(b, bitLength, 0)
}

func BenchmarkBitmap_SetZero_1K(b *testing.B) {
	benchmarkSet(b, 1000, 0)
}

func BenchmarkBitmap_SetZero_10K(b *testing.B) {
	benchmarkSet(b, 10000, 0)
}
