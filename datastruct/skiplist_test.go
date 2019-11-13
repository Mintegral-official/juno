package datastruct

import (
	"fmt"
	"github.com/Mintegral-official/juno/helpers"
	. "github.com/smartystreets/goconvey/convey"
	"math/rand"
	"testing"
	"time"
	"unsafe"
)

var arr []int

//生成count个[start,end)结束的不重复的随机数
func GenerateRandomNumber(start int, end int, count int) []int {
	//范围检查
	if end < start || (end-start) < count {
		return []int{0}
	}

	//存放结果的slice
	nums := []int{}
	i := 0
	//随机数生成器，加入时间戳保证每次生成的随机数不一样
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i < count {
		//生成随机数
		num := r.Intn((end - start)) + start

		//查重
		exist := false
		for _, v := range nums {
			if v == num {
				exist = true
				break
			}
		}

		if !exist {
			nums = append(nums, num)
			i++
		}
	}
	return nums
}

func init() {
	t := time.Now()
	arr = GenerateRandomNumber(0, 1500000000, 200000)
	fmt.Println(time.Since(t))

	fmt.Println(len(arr))

	var sl SkipList
	var el Element
	fmt.Printf("Structure sizes: SkipList is %v, Element is %v bytes\n", unsafe.Sizeof(sl), unsafe.Sizeof(el))
}

func TestNewSkipList(t *testing.T) {
	Convey("NewSKipList", t, func() {
		So(NewSkipList(DEFAULT_MAX_LEVEL, helpers.IntCompare), ShouldNotBeNil)
	})
}

func TestSkipList_Add_Del_Len(t *testing.T) {
	var s = NewSkipList(DEFAULT_MAX_LEVEL, helpers.IntCompare)
	var arr []int
	arr = GenerateRandomNumber(0, 1500000000, 100)
	for i := 0; i < 100; i++ {
		s.Add(arr[i], nil)
	}
	Convey("Add & Del & Len & Contains & Get", t, func() {
		So(s.Len(), ShouldEqual, 100)
		s.Del(arr[20])
		So(s.Len(), ShouldEqual, 99)
		So(s.Contains(arr[90]), ShouldBeTrue)
		_, err := s.Get(arr[34])
		So(err, ShouldBeNil)
	})
}

func TestSkipList_Get(t *testing.T) {
	//fmt.Println(s.findGE(-1, true, s.previousNodeCache))
	s := NewSkipList(DEFAULT_MAX_LEVEL, helpers.IntCompare)
	var arr []int
	arr = GenerateRandomNumber(0, 1500000000, 100)
	for i := 0; i < 100; i++ {
		s.Add(arr[i], nil)
	}
	Convey("findGE & findLT", t, func() {
		// 找到 ==  返回 true
		_, ok := s.findGE(arr[99], true, s.previousNodeCache)
		So(ok, ShouldBeTrue)
		// 找到 > 返回false
		_, ok = s.findGE(-1, true, s.previousNodeCache)
		So(ok, ShouldBeFalse)
		_, ok = s.findLT(arr[99])
		So(ok, ShouldBeTrue)
		_, ok = s.findLT(-1)
		So(ok, ShouldBeFalse)
	})
}

func add(s *SkipList, arr []int) {
	for i := 0; i < 200000; i++ {
		s.Add(arr[i], [1]byte{})
	}
}

func get(s *SkipList, arr []int) {
	for i := 0; i < 100000; i++ {
		_, _ = s.Get(arr[i])
	}
}

func BenchmarkNewSkipList_Add(b *testing.B) {
	var s = NewSkipList(DEFAULT_MAX_LEVEL, helpers.IntCompare)
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		add(s, arr) // BenchmarkNewSkipList_Add-8   	       3	 452658412 ns/op	18214261 B/op	  800000 allocs/op
	}
}

func BenchmarkSkipList_FindGE(b *testing.B) {
	var s = NewSkipList(DEFAULT_MAX_LEVEL, helpers.IntCompare)
	add(s, arr)
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		for j := 0; j < 100000; j++ {
			s.findGE(arr[j], true, s.previousNodeCache)
		}
	}
}

func BenchmarkSkipList_FindGE_RunParallel(b *testing.B) {
	var s = NewSkipList(DEFAULT_MAX_LEVEL, helpers.IntCompare)
	add(s, arr)
	b.ResetTimer()
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		// BenchmarkSkipListIterator_FindGE_RunParallel-8   	     300	   4641216 ns/op	   80010 B/op	   10000 allocs/op
		for pb.Next() {
			for i := 0; i < 100000; i++ {
				s.findGE(arr[i], true, s.previousNodeCache)
			}
		}
	})
}

func BenchmarkNewSkipList_FindLT(b *testing.B) {
	var s = NewSkipList(DEFAULT_MAX_LEVEL, helpers.IntCompare)
	add(s, arr)
	b.ResetTimer()
	b.ReportAllocs()
	// BenchmarkNewSkipList_FindLT-8   	2000000000	         0.01 ns/op	       0 B/op	       0 allocs/op
	for i := 0; i < b.N; i++ {
		for i := 0; i < 100000; i++ {
			s.findLT(arr[i])
		}
	}
}

func BenchmarkNewSkipList_FindLT_RunParallel(b *testing.B) {
	var s = NewSkipList(DEFAULT_MAX_LEVEL, helpers.IntCompare)
	add(s, arr)
	b.ResetTimer()
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			for i := 0; i < 100000; i++ {
				s.findLT(arr[i])
			}
		}
	})
}

func BenchmarkSkipList_Get(b *testing.B) {
	var s = NewSkipList(DEFAULT_MAX_LEVEL, helpers.IntCompare)
	add(s, arr)
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		get(s, arr)
	}
}

func BenchmarkSkipList_GetRunParallel(b *testing.B) {
	var s = NewSkipList(DEFAULT_MAX_LEVEL, helpers.IntCompare)
	add(s, arr)
	b.ResetTimer()
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		//BenchmarkSkipList_GetRunParallel-8   	     500	   3026856 ns/op	   80007 B/op	   10000 allocs/op
		for pb.Next() {
			get(s, arr)
		}
	})
}
