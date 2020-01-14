package query

import (
	"fmt"
	"github.com/Mintegral-official/juno/datastruct"
	"github.com/Mintegral-official/juno/document"
	"github.com/Mintegral-official/juno/query/check"
	"github.com/Mintegral-official/juno/query/operation"
	. "github.com/smartystreets/goconvey/convey"
	"math/rand"
	"testing"
	"time"
)

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

func add1(s *datastruct.SkipList, a []int) {
	for i := 0; i < 200000; i++ {
		s.Add(document.DocId(a[i]), i%248)
	}
}

func BenchmarkAndQuery_Next(b *testing.B) {

	var arr []int
	var arr1 []int
	var arr2 []int

	t := time.Now()
	arr = GenerateRandomNumber(0, 1500000000, 200000)
	arr1 = GenerateRandomNumber(0, 1500000000, 200000)
	arr2 = GenerateRandomNumber(0, 1500000000, 200000)
	fmt.Println(time.Since(t))

	a := datastruct.NewSkipList(datastruct.DefaultMaxLevel)
	a1 := datastruct.NewSkipList(datastruct.DefaultMaxLevel)
	a2 := datastruct.NewSkipList(datastruct.DefaultMaxLevel)
	add1(a, arr)
	add1(a1, arr1)
	add1(a2, arr2)
	var res = []int{1, 7, 78, 32, 99, 23, 44, 254, 127}
	var r = make([]interface{}, len(res))
	for i, v := range res {
		r[i] = v
	}
	var res1 = []int{1, 97, 123, 32, 99, 2, 127, 254, 23}
	var r1 = make([]interface{}, len(res1))
	for i, v := range res1 {
		r1[i] = v
	}

	s := NewAndQuery([]Query{
		NewAndQuery(
			[]Query{
				NewTermQuery(a.Iterator()),
			},
			[]check.Checker{
				check.NewInChecker(a.Iterator(), r, nil),
			},
		),
		NewTermQuery(a.Iterator()),
		NewOrQuery([]Query{
			NewTermQuery(a.Iterator()),
		}, []check.Checker{
			check.NewChecker(a.Iterator(), 1, operation.GE, nil),
			check.NewChecker(a.Iterator(), 39, operation.LT, nil),
			check.NewChecker(a.Iterator(), 10, operation.EQ, nil),
			check.NewChecker(a.Iterator(), 49, operation.NE, nil),
		}),
	}, []check.Checker{
		check.NewInChecker(a.Iterator(), r1, nil),
	}, )
	c := 0
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, err := s.Next()
		for err == nil {
			c++
			_, err = s.Next()
		}
	}
	b.StopTimer()
	fmt.Println(c)
}
