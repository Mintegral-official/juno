package datastruct

import (
	"fmt"
	"github.com/Mintegral-official/juno/document"
	. "github.com/smartystreets/goconvey/convey"
	"math/rand"
	"testing"
	"time"
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
	//var sl SkipList
	//var el Element
	//fmt.Printf("Structure sizes: SkipList is %v, Element is %v bytes\n", unsafe.Sizeof(sl), unsafe.Sizeof(el))
}

func TestNewSkipList(t *testing.T) {
	Convey("NewSKipList", t, func() {
		s := NewSkipList(DefaultMaxLevel)
		So(s, ShouldNotBeNil)
	})
}

func TestSkipList_Add_Del_Len(t *testing.T) {
	s := NewSkipList(DefaultMaxLevel)
	var arr []int
	arr = GenerateRandomNumber(0, 1500000000, 100)
	for i := 0; i < 100; i++ {
		s.Add(document.DocId(arr[i]), nil)
	}
	Convey("Add & Del & Len & Contains & Get", t, func() {
		So(s.Len(), ShouldEqual, 100)
		s.Del(document.DocId(arr[20]))
		So(s.Len(), ShouldEqual, 99)
		So(s.Contains(document.DocId(arr[90])), ShouldBeTrue)
		_, err := s.Get(document.DocId(arr[34]))
		So(err, ShouldBeNil)
	})
}

func TestSkipList_Get(t *testing.T) {
	s := NewSkipList(DefaultMaxLevel)
	var arr []int
	arr = GenerateRandomNumber(0, 1500000000, 100)
	for i := 0; i < 100; i++ {
		s.Add(document.DocId(arr[i]), nil)
	}
	Convey("findGE & findLT", t, func() {
		_, ok := s.findGE(document.DocId(arr[99]), true, s.previousNodeCache)
		So(ok, ShouldBeTrue)
		_, ok = s.findGE(0, true, s.previousNodeCache)
		So(ok, ShouldBeFalse)
		_, ok = s.findLT(document.DocId(arr[99]))
		So(ok, ShouldBeTrue)
		_, ok = s.findLT(0)
		So(ok, ShouldBeFalse)
	})
}

func add(s *SkipList, arr []int) {
	for i := 0; i < len(arr); i++ {
		s.Add(document.DocId(arr[i]), [1]byte{})
	}
}

func get(s *SkipList, arr []int) {
	for i := 0; i < len(arr)/2; i++ {
		_, _ = s.Get(document.DocId(arr[i]))
	}
}

func BenchmarkNewSkipList_Add(b *testing.B) {
	s := NewSkipList(DefaultMaxLevel)
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		add(s, arr)
	}
}

func BenchmarkSkipList_FindGE(b *testing.B) {
	s := NewSkipList(DefaultMaxLevel)
	add(s, arr)
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		for j := 0; j < len(arr)/2; j++ {
			s.findGE(document.DocId(arr[j]), true, s.previousNodeCache)
		}
	}
}

func BenchmarkSkipList_FindGE_RunParallel(b *testing.B) {
	s := NewSkipList(DefaultMaxLevel)
	add(s, arr)
	b.ResetTimer()
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			for i := 0; i < len(arr)/2; i++ {
				s.findGE(document.DocId(arr[i]), true, s.previousNodeCache)
			}
		}
	})
}

func BenchmarkNewSkipList_FindLT(b *testing.B) {
	s := NewSkipList(DefaultMaxLevel)
	add(s, arr)
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		for i := 0; i < len(arr)/2; i++ {
			s.findLT(document.DocId(arr[i]))
		}
	}
}

func BenchmarkNewSkipList_FindLT_RunParallel(b *testing.B) {
	s := NewSkipList(DefaultMaxLevel)
	add(s, arr)
	b.ResetTimer()
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			for i := 0; i < len(arr)/2; i++ {
				s.findLT(document.DocId(arr[i]))
			}
		}
	})
}

func BenchmarkSkipList_Get(b *testing.B) {
	s := NewSkipList(DefaultMaxLevel)
	add(s, arr)
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		get(s, arr)
	}
}

func BenchmarkSkipList_GetRunParallel(b *testing.B) {
	s := NewSkipList(DefaultMaxLevel)
	add(s, arr)
	b.ResetTimer()
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			get(s, arr)
		}
	})
}
