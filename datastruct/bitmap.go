package datastruct

import "github.com/MintegralTech/juno/document"

type BitSet []document.DocId

const (
	AddressBitsPerWord uint8          = 6
	WordsPerSize       document.DocId = 64
)

func NewBitMap() *BitSet {
	wordsLen := (2<<20 - 1) >> AddressBitsPerWord
	temp := BitSet(make([]document.DocId, wordsLen+1, wordsLen+1))
	return &temp
}

func (bs *BitSet) Set(bitIndex document.DocId) {
	wIndex := bs.wordIndex(bitIndex)
	bs.expandTo(wIndex)
	(*bs)[wIndex] |= document.DocId(0x01) << (bitIndex % WordsPerSize)
}

func (bs *BitSet) Del(bitIndex document.DocId) {
	wIndex := bs.wordIndex(bitIndex)
	if wIndex < len(*bs) {
		(*bs)[wIndex] &^= document.DocId(0x01) << (bitIndex % WordsPerSize)
	}
}

func (bs *BitSet) IsExist(bitIndex document.DocId) bool {
	wIndex := bs.wordIndex(bitIndex)
	return (wIndex < len(*bs)) && ((*bs)[wIndex]&(document.DocId(0x01)<<(bitIndex%WordsPerSize)) != 0)
}

func (bs *BitSet) IsFully() bool {
	return len(*bs) == cap(*bs)
}

func (bs BitSet) wordIndex(bitIndex document.DocId) int {
	return int(bitIndex >> AddressBitsPerWord)
}

func (bs *BitSet) expandTo(wordIndex int) {
	wordsRequired := wordIndex + 1
	if len(*bs) < wordsRequired {
		if wordsRequired < 2*len(*bs) {
			wordsRequired = 2 * len(*bs)
		}
		newCap := make([]document.DocId, wordsRequired, wordsRequired)
		copy(newCap, *bs)
		*bs = newCap
	}
}
