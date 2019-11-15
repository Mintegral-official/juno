package query

import (
	"github.com/Mintegral-official/juno/document"
	"github.com/Mintegral-official/juno/helpers"
)

type Heap []Query

func (h Heap) Compare(i, j interface{}) int {
	switch i.(type) {
	case document.DocId:
		return helpers.DocIdFunc(i, j)
	case int:
		return helpers.IntCompare(i, j)
	case string:
		return helpers.StringCompare(i, j)
	case float32:
		return helpers.Float32Compare(i, j)
	case float64:
		return helpers.Float64Compare(i, j)
	default:
		return helpers.IntCompare(i, j)
	}
}

func (h Heap) Len() int {
	return len(h)
}

func (h Heap) Less(i, j int) bool {
	iDocId, err := (h[i]).Current()
	if err != nil {
		return false
	}
	jDocId, err := (h[j]).Current()
	if err != nil {
		return true
	}
	return h.Compare(iDocId, jDocId) < 0
}

func (h Heap) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
}

func (h *Heap) Push(x interface{}) {
	*h = append(*h, x.(Query))
}

func (h *Heap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

func (h *Heap) Top() interface{} {
	if len(*h) == 0 {
		return nil
	}
	return (*h)[0]
}
