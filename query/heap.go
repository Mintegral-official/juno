package query

import "github.com/Mintegral-official/juno/helpers"

type Heap []Query

func (h Heap) Len() int {
	return len(h)
}

func (h Heap) Less(i, j int) bool {
	iDocId, iErr := (h[i]).Current()
	jDocId, jErr := (h[j]).Current()

	if iErr != nil && jErr != nil {
		return true
	}
	if iErr != nil {
		return false
	}
	if jErr != nil {
		return true
	}
	return helpers.Compare(iDocId, jDocId) < 0
}

func (h Heap) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
}

func (h *Heap) Push(x interface{}) {
	if x != nil {
		*h = append(*h, x.(Query))
	}
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
