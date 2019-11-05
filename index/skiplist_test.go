package index

import (
	"fmt"
	"github.com/Mintegral-official/juno/helpers"
	"math/rand"
	"testing"
	"time"
	"unsafe"
)

var s *SkipListIterator = NewSkipListIterator(DEFAULT_MAX_LEVEL, helpers.IntCompare)

func random(min, max int) int {
	rand.Seed(time.Now().Unix())
	return rand.Intn(max-min) + min
}

var arr [15000000]int

func init() {

	for i := 0; i < 15000000; i++ {
		arr[i] = random(0, 500000000)
	}

	for i := 0; i < 15000000; i++ {
		s.Add(arr[i], [1]byte{})
	}

	var sl SkipList
	var el Element
	fmt.Printf("Structure sizes: SkipList is %v, Element is %v bytes\n", unsafe.Sizeof(sl), unsafe.Sizeof(el))
}

func TestNewSkipList(t *testing.T) {
	go func() {
		for i := 0; i < 15000000; i++ {
			s.Add(arr[i], [1]byte{})
		}
	}()

	go func() {
		for i := 0; i < 5000000; i++ {
			s.Get(arr[i])
		}
	}()

	go func() {
		for i := 5000000; i < 10000001; i++ {
			s.Get(arr[i])
		}
	}()
	go func() {
		for i := 10000000; i < 15000000; i++ {
			s.Get(arr[i])
		}
	}()


}

func BenchmarkIncSet(b *testing.B) {
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		for j := 0; j < 15000000; j++ {
			s.Add(arr[i], [1]byte{})
		}
	}
	b.SetBytes(int64(b.N))
}

func BenchmarkIncGet(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		go func() {
			for j := 0; j < 15000000; j++ {
				s.Add(arr[i], [1]byte{})
			}
		}()
		go func() {
			for i := 0; i < 5000000; i++ {
				s.Get(arr[i])
			}
		}()

		go func() {
			for i := 5000000; i < 10000001; i++ {
				s.Get(arr[i])
			}
		}()
		go func() {
			for i := 10000000; i < 15000000; i++ {
				s.Get(arr[i])
			}
		}()
	}
	b.SetBytes(int64(b.N))
}

func BenchmarkNewSkipList(b *testing.B) {
	for s.HasNext() {
		fmt.Println(s.Next())
	}
	fmt.Println(s.Len())
}
