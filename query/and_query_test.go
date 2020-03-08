package query

import (
	"fmt"
	"github.com/Mintegral-official/juno/check"
	"github.com/Mintegral-official/juno/datastruct"
	"github.com/Mintegral-official/juno/document"
	"github.com/Mintegral-official/juno/operation"
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
	Convey("and query nil", t, func() {
		So(a, ShouldBeNil)
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

	Convey("and query get one query", t, func() {
		a := NewAndQuery([]Query{NewTermQuery(sl.Iterator())}, nil)

		testCase := [][]document.DocId{
			{1, 1}, {1, 1}, {1, 1}, {3, 3}, {4, 6}, {6, 6}, {7, 10}, {10, 10},
		}

		for _, expect := range testCase {
			v, e := a.GetGE(expect[0])
			So(v, ShouldEqual, expect[1])
			So(e, ShouldBeNil)
		}
		v, e := a.GetGE(11)
		So(v, ShouldEqual, 0)
		So(e, ShouldNotBeNil)
	})

	Convey("and query get two queries", t, func() {
		a := NewAndQuery([]Query{NewTermQuery(sl.Iterator()), NewTermQuery(sl1.Iterator())}, nil)
		testCase := [][]document.DocId{
			{1, 1}, {1, 1}, {1, 1}, {3, 6}, {3, 6}, {4, 6}, {6, 6},
		}
		for _, expect := range testCase {
			v, e := a.GetGE(expect[0])
			So(v, ShouldEqual, expect[1])
			So(e, ShouldBeNil)
		}
		v, e := a.GetGE(7)
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

	sl2 := datastruct.NewSkipList(datastruct.DefaultMaxLevel)

	sl2.Add(document.DocId(2), [1]byte{})
	sl2.Add(document.DocId(3), [1]byte{})
	sl2.Add(document.DocId(6), [1]byte{})
	sl2.Add(document.DocId(9), [1]byte{})

	Convey("and query next one query", t, func() {
		a := NewAndQuery([]Query{NewTermQuery(sl.Iterator())}, nil)

		testCase := []document.DocId{
			1, 3, 6, 10,
		}
		for _, expect := range testCase {
			v, e := a.Current()
			a.Next()
			So(v, ShouldEqual, expect)
			So(e, ShouldBeNil)
		}
		v, e := a.Current()
		a.Next()
		So(v, ShouldEqual, 0)
		So(e, ShouldNotBeNil)
	})

	Convey("and query next two queries", t, func() {
		a := NewAndQuery([]Query{NewTermQuery(sl.Iterator()), NewTermQuery(sl1.Iterator())}, nil)
		testCase := []document.DocId{
			1, 6,
		}
		for _, expect := range testCase {
			v, e := a.Current()
			a.Next()
			So(v, ShouldEqual, expect)
			So(e, ShouldBeNil)
		}
		v, e := a.Current()
		a.Next()
		So(v, ShouldEqual, 0)
		So(e, ShouldNotBeNil)

	})

	Convey("and query current three queries", t, func() {
		a := NewAndQuery([]Query{NewTermQuery(sl.Iterator()), NewTermQuery(sl1.Iterator()), NewTermQuery(sl2.Iterator())}, nil)

		v, e := a.Current()
		So(v, ShouldEqual, 6)
		So(e, ShouldEqual, nil)

		a.Next()

		v, e = a.Current()
		So(v, ShouldEqual, 0)
		So(e, ShouldNotBeNil)

	})
}

func TestAndQuery_Current(t *testing.T) {
	sl := datastruct.NewSkipList(datastruct.DefaultMaxLevel)

	sl.Add(document.DocId(1), [1]byte{})
	sl.Add(document.DocId(3), [1]byte{})
	sl.Add(document.DocId(6), [1]byte{})
	sl.Add(document.DocId(10), [1]byte{})

	Convey("and query current one query", t, func() {
		a := NewAndQuery([]Query{NewTermQuery(sl.Iterator())}, nil)
		v, e := a.Current()
		So(v, ShouldEqual, 1)
		So(e, ShouldBeNil)

		testCase := []document.DocId{1, 3, 6, 10}
		for _, expect := range testCase {
			v, e = a.Current()
			a.Next()
			So(v, ShouldEqual, expect)
			So(e, ShouldBeNil)
		}

		v, e = a.Current()
		a.Next()
		So(v, ShouldEqual, 0)
		So(e, ShouldNotBeNil)

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
				check.NewInChecker(a.Iterator(), r, nil, false),
			},
		),
		NewTermQuery(a.Iterator()),
		NewTermQuery(a.Iterator()),
		NewTermQuery(a.Iterator()),
		NewTermQuery(a.Iterator()),
		NewTermQuery(a.Iterator()),
		NewOrQuery([]Query{
			NewTermQuery(a.Iterator()),
		}, []check.Checker{
			check.NewChecker(a.Iterator(), 1, operation.GE, nil, false),
			check.NewChecker(a.Iterator(), 39, operation.LT, nil, false),
			check.NewChecker(a.Iterator(), 10, operation.EQ, nil, false),
			check.NewChecker(a.Iterator(), 49, operation.NE, nil, false),
		}),
	}, []check.Checker{
		check.NewInChecker(a.Iterator(), r1, nil, false),
	})
	c := 0
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, err := s.Current()
		s.Next()
		for err == nil {
			c++
			_, err = s.Current()
			s.Next()
		}
	}
	b.StopTimer()
	fmt.Println(c)
}

func TestNewAndQuery_check(t *testing.T) {
	sl := datastruct.NewSkipList(datastruct.DefaultMaxLevel)

	sl.Add(document.DocId(1), 1)
	sl.Add(document.DocId(3), 1)
	sl.Add(document.DocId(6), 2)
	sl.Add(document.DocId(10), 2)

	sl1 := datastruct.NewSkipList(datastruct.DefaultMaxLevel)

	sl1.Add(document.DocId(1), 1)
	sl1.Add(document.DocId(4), 1)
	sl1.Add(document.DocId(6), 2)
	sl1.Add(document.DocId(9), 2)

	sl2 := datastruct.NewSkipList(datastruct.DefaultMaxLevel)

	sl2.Add(document.DocId(1), 1)
	sl2.Add(document.DocId(4), 1)
	sl2.Add(document.DocId(6), 2)
	sl2.Add(document.DocId(9), 2)

	Convey("and query with check", t, func() {
		a := NewAndQuery([]Query{
			NewTermQuery(sl.Iterator()),
			NewTermQuery(sl1.Iterator()),
		}, []check.Checker{
			check.NewChecker(sl2.Iterator(), 2, operation.EQ, nil, false),
		})
		testCase := []document.DocId{6}

		for _, expect := range testCase {
			v, e := a.Current()
			a.Next()
			So(v, ShouldEqual, expect)
			So(e, ShouldBeNil)
		}
		v, e := a.Current()
		a.Next()
		So(v, ShouldEqual, 0)
		So(e, ShouldNotBeNil)
	})

	Convey("and query with In check", t, func() {
		a := NewAndQuery([]Query{
			NewTermQuery(sl.Iterator()),
			NewTermQuery(sl1.Iterator()),
		}, []check.Checker{
			check.NewInChecker(sl2.Iterator(), []int{1, 2, 3}, nil, false),
		})
		testCase := []document.DocId{1, 6}

		for _, expect := range testCase {
			v, e := a.Current()
			a.Next()
			So(v, ShouldEqual, expect)
			So(e, ShouldBeNil)
		}
		v, e := a.Current()
		a.Next()
		So(v, ShouldEqual, 0)
		So(e, ShouldNotBeNil)
	})

	Convey("and query with and check", t, func() {
		a := NewAndQuery([]Query{
			NewTermQuery(sl.Iterator()),
			NewTermQuery(sl1.Iterator()),
		}, []check.Checker{
			check.NewAndChecker([]check.Checker{
				check.NewInChecker(sl2.Iterator(), []int{1, 2, 3}, nil, false),
				check.NewChecker(sl2.Iterator(), 2, operation.EQ, nil, false),
			}),
			check.NewInChecker(sl2.Iterator(), []int{1, 2, 3}, nil, false),
		})
		testCase := []document.DocId{6}

		for _, expect := range testCase {
			v, e := a.Current()
			a.Next()
			So(v, ShouldEqual, expect)
			So(e, ShouldBeNil)
		}
		v, e := a.Current()
		a.Next()
		So(v, ShouldEqual, 0)
		So(e, ShouldNotBeNil)
	})

	Convey("and query debug", t, func() {
		a := NewAndQuery([]Query{
			NewTermQuery(sl.Iterator()),
			NewTermQuery(sl1.Iterator()),
		}, []check.Checker{
			check.NewOrChecker([]check.Checker{
				check.NewInChecker(sl2.Iterator(), []int{1, 2, 3}, nil, false),
				check.NewChecker(sl2.Iterator(), 2, operation.EQ, nil, false),
			}),
			check.NewInChecker(sl2.Iterator(), []int{1, 2, 3}, nil, false),
		})
		testCase := []document.DocId{1, 6}
		a.SetDebug(1)
		for _, expect := range testCase {
			v, e := a.Current()
			a.Next()
			So(v, ShouldEqual, expect)
			So(e, ShouldBeNil)
		}
		v, e := a.Current()
		a.Next()
		So(v, ShouldEqual, 0)
		So(e, ShouldNotBeNil)
		So(a.DebugInfo().String(), ShouldNotEqual, "")
	})
}
