package query

import (
	"github.com/Mintegral-official/juno/check"
	"github.com/Mintegral-official/juno/datastruct"
	"github.com/Mintegral-official/juno/document"
	"github.com/Mintegral-official/juno/operation"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestNotAndQuery(t *testing.T) {
	a := NewNotAndQuery([]Query{}, nil)
	Convey("not and query nil", t, func() {
		So(a, ShouldBeNil)
	})
}

func TestNewNotAndQuery_Next(t *testing.T) {
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
	sl2.Add(document.DocId(3), [1]byte{})

	Convey("not and query next (one query)", t, func() {
		a := NewNotAndQuery([]Query{NewTermQuery(sl.Iterator())}, nil)
		testCase := []document.DocId{1, 3, 6, 10}
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

	Convey("not and query next (three query)", t, func() {
		a := NewNotAndQuery([]Query{
			NewTermQuery(sl.Iterator()), NewTermQuery(sl1.Iterator()), NewTermQuery(sl2.Iterator()),
		}, nil)
		So(a, ShouldNotBeNil)
		v, e := a.Current()
		a.Next()
		So(v, ShouldEqual, 10)
		So(e, ShouldBeNil)

		v, e = a.Current()
		a.Next()
		So(v, ShouldEqual, 0)
		So(e, ShouldNotBeNil)
	})
}

func TestNewNotAndQuery_GetGE(t *testing.T) {
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

	Convey("not and query getGE (one query)", t, func() {
		s1 := sl.Iterator()
		a := NewNotAndQuery([]Query{NewTermQuery(s1)}, nil)
		testCase := [][]document.DocId{
			{1, 1}, {2, 3}, {3, 3}, {4, 6}, {5, 6}, {6, 6}, {7, 10}, {8, 10}, {9, 10}, {10, 10},
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

	Convey("not and query getGE (two queries)", t, func() {
		s1 := sl.Iterator()
		s2 := sl1.Iterator()
		a := NewNotAndQuery([]Query{NewTermQuery(s1), NewTermQuery(s2)}, nil)

		testCase := [][]document.DocId{
			{1, 3}, {2, 3}, {3, 3}, {4, 10}, {5, 10}, {6, 10}, {7, 10}, {8, 10}, {9, 10}, {10, 10},
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
}

func TestNewAndQuery_Check(t *testing.T) {
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

	Convey("not and query with check", t, func() {
		a := NewNotAndQuery([]Query{
			NewTermQuery(sl.Iterator()),
			NewTermQuery(sl1.Iterator()),
		}, []check.Checker{
			check.NewChecker(sl2.Iterator(), 2, operation.EQ, nil, false),
		})
		v, e := a.Current()
		a.Next()
		So(v, ShouldEqual, 0)
		So(e, ShouldNotBeNil)
	})

	Convey("not and query with in check", t, func() {
		a := NewNotAndQuery([]Query{
			NewTermQuery(sl.Iterator()),
			NewTermQuery(sl1.Iterator()),
		}, []check.Checker{
			check.NewInChecker(sl2.Iterator(), []int{1, 2, 3}, nil, false),
		})
		v, e := a.Current()
		a.Next()
		So(v, ShouldEqual, 0)
		So(e, ShouldNotBeNil)
	})

	Convey("not and query with or check", t, func() {
		a := NewNotAndQuery([]Query{
			NewTermQuery(sl.Iterator()),
			NewTermQuery(sl1.Iterator()),
		}, []check.Checker{
			check.NewOrChecker([]check.Checker{
				check.NewInChecker(sl2.Iterator(), []int{1, 2, 3}, nil, false),
				check.NewChecker(sl2.Iterator(), 2, operation.EQ, nil, false),
			}),
			check.NewInChecker(sl2.Iterator(), []int{1, 2, 3}, nil, false),
		})

		v, e := a.Current()
		a.Next()
		So(v, ShouldEqual, 0)
		So(e, ShouldNotBeNil)
	})

	Convey("not and query with not and check", t, func() {
		a := NewNotAndQuery([]Query{
			NewTermQuery(sl.Iterator()),
			NewTermQuery(sl1.Iterator()),
		}, []check.Checker{
			check.NewNotAndChecker([]check.Checker{
				check.NewInChecker(sl2.Iterator(), []int{1, 2, 3}, nil, false),
				check.NewChecker(sl2.Iterator(), 2, operation.EQ, nil, false),
			}),
			check.NewInChecker(sl2.Iterator(), []int{1, 2, 3}, nil, false),
		})
		v, e := a.Current()
		a.Next()
		So(v, ShouldEqual, 0)
		So(e, ShouldNotBeNil)
	})
}
