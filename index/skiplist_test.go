package index
//
//import (
//	"fmt"
//	"github.com/Mintegral-official/juno/helpers"
//	"math/rand"
//	"testing"
//	"time"
//	"unsafe"
//)
//
//var s *SkipListIterator = NewSKipListIterator(DEFAULT_MAX_LEVEL, helpers.IntCompare)
//
//func random(min, max int) int {
//	rand.Seed(time.Now().Unix())
//	return rand.Intn(max-min) + min
//}
//
//var arr [500001]int
//
//func init() {
//
//	for i := 0; i < 15000000; i++ {
//		s.Add(i, [1]byte{})
//	}
//
//	for i := 0; i < 500000; i++ {
//		arr[i] = random(0, 15000000)
//	}
//
//	var sl SkipList
//	var el Element
//	fmt.Printf("Structure sizes: SkipList is %v, Element is %v bytes\n", unsafe.Sizeof(sl), unsafe.Sizeof(el))
//}
//
//func TestNewSkipList(t *testing.T) {
//	go func() {
//		for i := 0; i < 1000000; i++ {
//			s.Add(i, [1]byte{})
//		}
//	}()
//
//	go func() {
//		for i := 0; i < 1000000; i++ {
//			s.Del(i)
//		}
//	}()
//
//}
//
//func BenchmarkIncSet(b *testing.B) {
//	b.ReportAllocs()
//
//	for i := 0; i < b.N; i++ {
//		for j := 0; j < 15000000; j++ {
//			s.Add(i, [1]byte{})
//		}
//	}
//	b.SetBytes(int64(b.N))
//}
//
//func BenchmarkIncGet(b *testing.B) {
//	b.ReportAllocs()
//	// fmt.Println(s.length)
//	for i := 0; i < b.N; i++ {
//		for j := 0; j < 500001; j++ {
//			s.Get(arr[j])
//		}
//	}
//	b.SetBytes(int64(b.N))
//}
//
//func BenchmarkNewSkipList(b *testing.B) {
//	for s.HasNext() {
//		s.Next()
//	}
//}
