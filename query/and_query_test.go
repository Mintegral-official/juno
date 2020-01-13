package query

import (
	"fmt"
	"github.com/Mintegral-official/juno/datastruct"
	"github.com/Mintegral-official/juno/document"
	"github.com/Mintegral-official/juno/query/check"
	. "github.com/smartystreets/goconvey/convey"
	"math/rand"
	"testing"
	"time"
)

var arr []int
var arr1 []int
var arr2 []int

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
	arr1 = GenerateRandomNumber(0, 1500000000, 200000)
	arr2 = GenerateRandomNumber(0, 1500000000, 200000)
	fmt.Println(time.Since(t))
	fmt.Println(len(arr))
	//var sl SkipList
	//var el Element
	//fmt.Printf("Structure sizes: SkipList is %v, Element is %v bytes\n", unsafe.Sizeof(sl), unsafe.Sizeof(el))
}

func TestAndQuery(t *testing.T) {
	a := NewAndQuery(nil, nil)
	Convey("and query", t, func() {
		So(a, ShouldNotBeNil)
	})
}

func TestAndQuery_GetGE(t *testing.T) {
	sl := datastruct.NewSkipList(datastruct.DefaultMaxLevel)

	sl.Add(document.DocId(1), [1]byte{})
	sl.Add(document.DocId(3), [1]byte{})
	sl.Add(document.DocId(6), [1]byte{})
	sl.Add(document.DocId(10), [1]byte{})

	sl1 := datastruct.NewSkipList(datastruct.DefaultMaxLevel)

	sl1.Add(document.DocId(1), [1]byte{})
	sl1.Add(document.DocId(4), [1]byte{})
	sl1.Add(document.DocId(6), [1]byte{})
	sl1.Add(document.DocId(9), [1]byte{})

	Convey("and query get1", t, func() {
		a := NewAndQuery([]Query{NewTermQuery(sl.Iterator())}, nil)

		v, e := a.GetGE(document.DocId(1))
		So(v, ShouldEqual, 1)
		So(e, ShouldBeNil)

		v, e = a.GetGE(document.DocId(1))
		So(v, ShouldEqual, 1)
		So(e, ShouldBeNil)

		v, e = a.GetGE(document.DocId(1))
		So(v, ShouldEqual, 1)
		So(e, ShouldBeNil)

		v, e = a.GetGE(document.DocId(3))
		So(v, ShouldEqual, 3)
		So(e, ShouldBeNil)

		v, e = a.GetGE(document.DocId(4))
		So(v, ShouldEqual, 6)
		So(e, ShouldBeNil)

		v, e = a.GetGE(document.DocId(6))
		So(v, ShouldEqual, 6)
		So(e, ShouldBeNil)

		v, e = a.GetGE(document.DocId(7))
		So(v, ShouldEqual, 10)
		So(e, ShouldBeNil)
	})

	Convey("and query get2", t, func() {
		a := NewAndQuery([]Query{NewTermQuery(sl.Iterator()), NewTermQuery(sl1.Iterator())}, nil)
		v, e := a.GetGE(document.DocId(1))
		So(v, ShouldEqual, 1)
		So(e, ShouldBeNil)

		v, e = a.GetGE(document.DocId(1))
		So(v, ShouldEqual, 1)
		So(e, ShouldBeNil)

		v, e = a.GetGE(document.DocId(1))
		So(v, ShouldEqual, 1)
		So(e, ShouldBeNil)

		v, e = a.GetGE(document.DocId(1))
		So(v, ShouldEqual, 1)
		So(e, ShouldBeNil)

		v, e = a.GetGE(document.DocId(3))
		So(v, ShouldEqual, 6)
		So(e, ShouldBeNil)

		v, e = a.GetGE(document.DocId(3))
		So(v, ShouldEqual, 6)
		So(e, ShouldBeNil)

		v, e = a.GetGE(document.DocId(4))
		So(v, ShouldEqual, 6)
		So(e, ShouldBeNil)

		v, e = a.GetGE(document.DocId(6))
		So(v, ShouldEqual, 6)
		So(e, ShouldBeNil)
		fmt.Println()

		v, e = a.GetGE(document.DocId(10))
		So(v, ShouldEqual, 0)
		So(e, ShouldNotBeNil)
	})

}

func TestAndQuery_Next(t *testing.T) {
	sl := datastruct.NewSkipList(datastruct.DefaultMaxLevel)

	sl.Add(document.DocId(1), [1]byte{})
	sl.Add(document.DocId(3), [1]byte{})
	sl.Add(document.DocId(6), [1]byte{})
	sl.Add(document.DocId(10), [1]byte{})

	sl1 := datastruct.NewSkipList(datastruct.DefaultMaxLevel)

	sl1.Add(document.DocId(1), [1]byte{})
	sl1.Add(document.DocId(4), [1]byte{})
	sl1.Add(document.DocId(6), [1]byte{})
	sl1.Add(document.DocId(9), [1]byte{})

	Convey("and query next1", t, func() {
		a := NewAndQuery([]Query{NewTermQuery(sl.Iterator())}, nil)
		v, e := a.Current()
		So(v, ShouldEqual, 1)
		So(e, ShouldBeNil)

		v, e = a.Next()
		So(v, ShouldEqual, 3)
		So(e, ShouldBeNil)

		v, e = a.Next()
		So(v, ShouldEqual, 6)
		So(e, ShouldBeNil)

		v, e = a.Next()
		So(v, ShouldEqual, 10)
		So(e, ShouldBeNil)
	})

	Convey("and query next2", t, func() {
		a := NewAndQuery([]Query{NewTermQuery(sl.Iterator()), NewTermQuery(sl1.Iterator())}, nil)
		v, e := a.Current()
		So(v, ShouldEqual, 1)
		So(e, ShouldBeNil)

		v, e = a.Next()
		So(v, ShouldEqual, 6)
		So(e, ShouldBeNil)

	})
}

func TestAndQuery_Current(t *testing.T) {
	sl := datastruct.NewSkipList(datastruct.DefaultMaxLevel)

	sl.Add(document.DocId(1), [1]byte{})
	sl.Add(document.DocId(3), [1]byte{})
	sl.Add(document.DocId(6), [1]byte{})
	sl.Add(document.DocId(10), [1]byte{})

	Convey("and query current", t, func() {
		a := NewAndQuery([]Query{NewTermQuery(sl.Iterator())}, nil)
		v, e := a.Current()
		So(v, ShouldEqual, 1)
		So(e, ShouldBeNil)

	})
}

func add1(s *datastruct.SkipList) {
	for i := 0; i < 200000; i++ {
		s.Add(document.DocId(arr[i]), i%248)
	}
}

func BenchmarkAndQuery_Next(b *testing.B) {
	a := datastruct.NewSkipList(datastruct.DefaultMaxLevel)
	a1 := datastruct.NewSkipList(datastruct.DefaultMaxLevel)
	a2 := datastruct.NewSkipList(datastruct.DefaultMaxLevel)
	add1(a)
	add1(a1)
	add1(a2)
	var res = []int{1, 7, 78, 32, 99, 23}
	var r = make([]interface{}, len(res))
	for i, v := range res {
		r[i] = v
	}

	var res1 = []int{1, 97, 123, 346, 32, 99, 2,}
	var r1 = make([]interface{}, len(res1))
	for i, v := range res1 {
		r1[i] = v
	}
	s := NewAndQuery(
		[]Query{
			NewTermQuery(a.Iterator()),
			NewTermQuery(a1.Iterator()),
			NewTermQuery(a2.Iterator()),
		},
		[]check.Checker{
			check.NewInChecker(a.Iterator(), r, nil),
			check.NewInChecker(a1.Iterator(), r1, nil),
		})
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, err := s.Next()
		for err == nil {
			_, err = s.Next()
		}
	}
}
