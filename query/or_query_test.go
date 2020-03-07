package query

import (
	"github.com/Mintegral-official/juno/check"
	"github.com/Mintegral-official/juno/datastruct"
	"github.com/Mintegral-official/juno/document"
	"github.com/Mintegral-official/juno/operation"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestNewOrQuery(t *testing.T) {
	Convey("or query nil", t, func() {
		oq := NewOrQuery(nil, nil)
		So(oq, ShouldBeNil)
	})
}

func TestNewOrQuery_Next(t *testing.T) {
	sl := datastruct.NewSkipList(datastruct.DefaultMaxLevel)

	sl.Add(document.DocId(1), 1)
	sl.Add(document.DocId(3), 2)
	sl.Add(document.DocId(6), 2)
	sl.Add(document.DocId(10), 1)

	sl1 := datastruct.NewSkipList(datastruct.DefaultMaxLevel)

	sl1.Add(document.DocId(1), [1]byte{})
	sl1.Add(document.DocId(4), [1]byte{})
	sl1.Add(document.DocId(6), [1]byte{})
	sl1.Add(document.DocId(9), [1]byte{})
	sl1.Add(document.DocId(10), [1]byte{})
	sl1.Add(document.DocId(94), [1]byte{})
	sl1.Add(document.DocId(944), [1]byte{})

	sl2 := datastruct.NewSkipList(datastruct.DefaultMaxLevel)

	sl2.Add(document.DocId(2), [1]byte{})
	sl2.Add(document.DocId(5), [1]byte{})
	sl2.Add(document.DocId(6), [1]byte{})
	sl2.Add(document.DocId(8), [1]byte{})

	Convey("or query next(one query)", t, func() {
		a := NewOrQuery([]Query{NewTermQuery(sl.Iterator())}, []check.Checker{
			check.NewChecker(sl.Iterator(), 1, operation.EQ, nil, false),
		})

		testCase := []document.DocId{1, 10}
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

	Convey("or query next (three queries)", t, func() {
		a := NewOrQuery([]Query{NewTermQuery(sl2.Iterator()), NewTermQuery(sl.Iterator()), NewTermQuery(sl1.Iterator())}, nil)

		expectVec := []document.DocId{
			1, 2, 3, 4, 5, 6, 8, 9, 10, 94, 944,
		}

		for _, expect := range expectVec {
			v, e := a.Current()
			So(v, ShouldEqual, expect)
			So(e, ShouldBeNil)
			a.Next()
		}

		a.Next()
		v, e := a.Current()
		So(v, ShouldEqual, 0)
		So(e, ShouldNotBeNil)
	})
}

func TestOrQuery_GetGE(t *testing.T) {
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

	Convey("or query getGE (two queries)", t, func() {
		s1 := sl.Iterator()
		s2 := sl1.Iterator()
		a := NewOrQuery([]Query{NewTermQuery(s1), NewTermQuery(s2)}, nil)

		testCases := [][]document.DocId{
			{1, 1},
			{2, 3},
			{3, 3},
			{4, 4},
			{5, 6},
			{6, 6},
			{7, 9},
			{8, 9},
			{9, 9},
		}

		for _, c := range testCases {
			v, e := a.GetGE(c[0])
			So(v, ShouldEqual, c[1])
			So(e, ShouldBeNil)
		}

		v, e := a.GetGE(10)
		So(v, ShouldEqual, 10)
		So(e, ShouldBeNil)

		v, e = a.GetGE(11)
		So(v, ShouldEqual, 0)
		So(e, ShouldNotBeNil)
	})

	Convey("or query getGE (one query)", t, func() {
		a := NewOrQuery([]Query{NewTermQuery(sl.Iterator())}, nil)
		v, e := a.GetGE(8)
		So(v, ShouldEqual, 10)
		So(e, ShouldBeNil)

		v, e = a.GetGE(2)
		So(v, ShouldEqual, 10)
		So(e, ShouldBeNil)

		v, e = a.GetGE(11)
		So(v, ShouldEqual, 0)
		So(e, ShouldNotBeNil)

	})
}

func TestNewOrQuery2(t *testing.T) {
	sl := datastruct.NewSkipList(datastruct.DefaultMaxLevel)

	sl.Add(document.DocId(1), [1]byte{})
	sl.Add(document.DocId(3), [1]byte{})
	sl.Add(document.DocId(6), [1]byte{})
	sl.Add(document.DocId(10), [1]byte{})

	sl1 := datastruct.NewSkipList(datastruct.DefaultMaxLevel)

	sl1.Add(document.DocId(2), [1]byte{})
	sl1.Add(document.DocId(4), [1]byte{})
	sl1.Add(document.DocId(6), [1]byte{})
	sl1.Add(document.DocId(9), [1]byte{})

	sl2 := datastruct.NewSkipList(datastruct.DefaultMaxLevel)

	sl2.Add(document.DocId(1), [1]byte{})
	sl2.Add(document.DocId(4), [1]byte{})
	sl2.Add(document.DocId(6), [1]byte{})
	sl2.Add(document.DocId(10), [1]byte{})
	sl2.Add(document.DocId(100), [1]byte{})

	Convey("or query next (TermQuery AndQuery)", t, func() {

		q := NewOrQuery([]Query{
			NewTermQuery(sl.Iterator()),
			NewAndQuery([]Query{
				NewTermQuery(sl1.Iterator()),
				NewTermQuery(sl2.Iterator()),
			}, nil),
		}, nil)

		expectVec := []document.DocId{
			1, 3, 4, 6, 10,
		}

		for _, expect := range expectVec {
			v, e := q.Current()
			So(v, ShouldEqual, expect)
			So(e, ShouldBeNil)
			q.Next()
		}

		q.Next()
		v, e := q.Current()
		So(v, ShouldEqual, 0)
		So(e, ShouldNotBeNil)
	})

	Convey("and query next (three TermQuery)", t, func() {

		q := NewAndQuery([]Query{
			NewTermQuery(sl.Iterator()),
			NewTermQuery(sl1.Iterator()),
			NewTermQuery(sl2.Iterator()),
		}, nil)

		expectVec := []document.DocId{6}

		for _, expect := range expectVec {
			v, e := q.Current()
			So(v, ShouldEqual, expect)
			So(e, ShouldBeNil)
			q.Next()
		}

		q.Next()
		v, e := q.Current()
		So(v, ShouldEqual, 0)
		So(e, ShouldNotBeNil)
	})
}

func TestNewOrQuery_Next_check(t *testing.T) {
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

	Convey("or query with check", t, func() {
		a := NewOrQuery([]Query{
			NewTermQuery(sl.Iterator()),
			NewTermQuery(sl1.Iterator()),
		}, []check.Checker{
			check.NewChecker(sl2.Iterator(), 2, operation.EQ, nil, false),
		})
		testCase := []document.DocId{6, 9}

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

	Convey("or query with In checke", t, func() {
		a := NewOrQuery([]Query{
			NewTermQuery(sl.Iterator()),
			NewTermQuery(sl1.Iterator()),
		}, []check.Checker{
			check.NewInChecker(sl2.Iterator(), []int{1, 2, 3}, nil, false),
		})
		testCase := []document.DocId{1, 4, 6, 9}

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

	Convey("or query with or check", t, func() {
		a := NewOrQuery([]Query{
			NewTermQuery(sl.Iterator()),
			NewTermQuery(sl1.Iterator()),
		}, []check.Checker{
			check.NewOrChecker([]check.Checker{
				check.NewInChecker(sl2.Iterator(), []int{1, 2, 3}, nil, false),
				check.NewChecker(sl2.Iterator(), 2, operation.EQ, nil, false),
			}),
			check.NewInChecker(sl2.Iterator(), []int{1, 2, 3}, nil, false),
		})
		testCase := []document.DocId{1, 4, 6, 9}

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

	Convey("or query check4", t, func() {
		a := NewOrQuery([]Query{
			NewTermQuery(sl.Iterator()),
			NewTermQuery(sl1.Iterator()),
		}, []check.Checker{
			check.NewOrChecker([]check.Checker{
				check.NewInChecker(sl2.Iterator(), []int{1, 2, 3}, nil, false),
				check.NewChecker(sl2.Iterator(), 2, operation.EQ, nil, false),
			}),
			check.NewInChecker(sl2.Iterator(), []int{1, 2, 3}, nil, false),
		})
		testCase := []document.DocId{1, 4, 6, 9}
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
		So(a.DebugInfo().String(), ShouldNotBeNil)
	})
}
