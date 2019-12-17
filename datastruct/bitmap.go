package datastruct

type BitSet []uint64

const (
	AddressBitsPerWord uint8  = 6
	WordsPerSize       uint64 = 64
)

func NewBitMap() *BitSet {
	wordsLen := (2<<20 - 1) >> AddressBitsPerWord
	temp := BitSet(make([]uint64, wordsLen+1, wordsLen+1))
	return &temp
}

func (bs *BitSet) Set(bitIndex uint64) {
	wIndex := bs.wordIndex(bitIndex)
	bs.expandTo(wIndex)
	(*bs)[wIndex] |= uint64(0x01) << (bitIndex % WordsPerSize)
}

func (bs *BitSet) Del(bitIndex uint64) {
	wIndex := bs.wordIndex(bitIndex)
	if wIndex < len(*bs) {
		(*bs)[wIndex] &^= uint64(0x01) << (bitIndex % WordsPerSize)
	}
}

func (bs *BitSet) IsExist(bitIndex uint64) bool {
	wIndex := bs.wordIndex(bitIndex)
	return (wIndex < len(*bs)) && ((*bs)[wIndex]&(uint64(0x01)<<(bitIndex%WordsPerSize)) != 0)
}

func (bs *BitSet) IsFully() bool {
	return len(*bs) == cap(*bs)
}

func (bs BitSet) wordIndex(bitIndex uint64) int {
	return int(bitIndex >> AddressBitsPerWord)
}

func (bs *BitSet) expandTo(wordIndex int) {
	wordsRequired := wordIndex + 1
	if len(*bs) < wordsRequired {
		if wordsRequired < 2*len(*bs) {
			wordsRequired = 2 * len(*bs)
		}
		newCap := make([]uint64, wordsRequired, wordsRequired)
		copy(newCap, *bs)
		*bs = newCap
	}
}
