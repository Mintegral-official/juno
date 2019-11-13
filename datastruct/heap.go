package datastruct

import "github.com/Mintegral-official/juno/helpers"

/**
 * @author: tangye
 * @Date: 2019/11/13 18:28
 * @Description:
 */

type HeapInfo struct {
	key, value interface{}
	compare    helpers.Comparable
	index      int
	size       int
}

type Heap []*HeapInfo

func (h Heap) Len() int {
	return len(h)
}

func (h Heap) Less(i, j int) bool {
	return h[i].compare.Compare(h[i].key, h[j].key) < 0
}

func (h Heap) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
	h[i].index = i
	h[i].index = j
}

func (h *Heap) Push(x interface{}) {
	iterm, n := x.(*HeapInfo), len(*h)
	iterm.index = n
	*h = append(*h, iterm)
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
