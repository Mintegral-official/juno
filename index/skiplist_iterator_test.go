package index

import (
	"fmt"
	"github.com/Mintegral-official/juno/helpers"
	"testing"
)

var s1 = NewSkipListIterator(DEFAULT_MAX_LEVEL, helpers.IntCompare)

func init() {
	for i := 0; i < 100; i++ {
		s1.Add(i, nil)
	}
}

func Test(t *testing.T) {
	fmt.Println(s1.GetGE(5))
	fmt.Println(s1.GetGE(10))
	c := 0
	for s1.HasNext() {
		s1.Next()
		c++
		if c == 10 {
			break
		}
	}
	fmt.Println(s1.index)

	fmt.Println(s.findGE(5, false, s1.previousNodeCache))
}