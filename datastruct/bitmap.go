package datastruct

import (
	"fmt"
)

const bitSize = 8

var bitInit = []byte{1, 1 << 1, 1 << 2, 1 << 3, 1 << 4, 1 << 5, 1 << 6, 1 << 7}

type BitMap struct {
	bits     []byte
	bitCount int
	capacity int
}

func NewBitMap(maxNum int) *BitMap {
	return &BitMap{
		bits:     make([]byte, (maxNum+7)/bitSize),
		bitCount: 0,
		capacity: maxNum,
	}
}

func (bm *BitMap) Set(num int) {
	byteIndex, bitPos := bm.offset(num)
	bm.bits[byteIndex] |= bitInit[bitPos]
	bm.bitCount++
}

func (bm *BitMap) Del(num int) {
	byteIndex, bitPos := bm.offset(num)
	bm.bits[byteIndex] &= ^bitInit[bitPos]
	bm.bitCount--
}

func (bm *BitMap) IsExist(num int) bool {
	byteIndex := num / bitSize
	if byteIndex >= len(bm.bits) {
		return false
	}
	bitPos := num % bitSize
	return bm.bits[byteIndex]&bitInit[bitPos] != 0
}

func (bm *BitMap) offset(num int) (byteIndex int, bitPos byte) {
	byteIndex = num / bitSize
	if byteIndex >= len(bm.bits) {
		panic(fmt.Sprintf(" error: index value %d is out of range ", byteIndex))
		return
	}
	bitPos = byte(num % bitSize)
	return byteIndex, bitPos
}

func (bm *BitMap) Size() int {
	return len(bm.bits) * bitSize
}

func (bm *BitMap) IsEmpty() bool {
	return bm.bitCount == 0
}

func (bm *BitMap) IsFully() bool {
	return bm.bitCount == bm.capacity
}

func (bm *BitMap) Count() int {
	return bm.bitCount
}

func (bm *BitMap) Get() []int {
	var data []int
	count := bm.Size()
	for index := 0; index < count; index++ {
		if bm.IsExist(index) {
			data = append(data, index)
		}
	}
	return data
}
