package query

import (
	"fmt"
	"github.com/Mintegral-official/juno/check"
	"github.com/Mintegral-official/juno/document"
	"github.com/Mintegral-official/juno/helpers"
	"github.com/Mintegral-official/juno/index"
	"github.com/Mintegral-official/juno/operation"
	. "github.com/smartystreets/goconvey/convey"
	"reflect"
	"strings"
	"testing"
)

func TestNewTermQuery(t *testing.T) {
	ss := index.NewIndex("")
	s := ss.GetInvertedIndex()
	Convey("Add", t, func() {
		So(s.Add("fieldName\0071", document.DocId(1)), ShouldBeNil)
		So(s.Add("fieldName\0071", document.DocId(5)), ShouldBeNil)
		So(s.Add("fieldName\0071", document.DocId(6)), ShouldBeNil)
		So(s.Add("fieldName\0071", document.DocId(7)), ShouldBeNil)
		So(s.Add("fieldName\0074", document.DocId(2)), ShouldBeNil)
		b := s.Iterator("fieldName", "1")
		tq := NewTermQuery(b)
		res := tq.Marshal()
		fmt.Println(res)

		tq.SetDebug(1)
		fmt.Println(tq.Current())
		for i := 0; i < 4; i++ {
			tq.Next()
			fmt.Println(tq.Current())
		}
		fmt.Println(tq.DebugInfo())

		sss := tq.Unmarshal(ss, res, nil)
		fmt.Println(sss.Current())
		for i := 0; i < 4; i++ {
			tq.Next()
			fmt.Println(tq.Current())
		}
		fmt.Println(sss.DebugInfo())

	})
}

func TestTermQuery_Current(t *testing.T) {
	ss := index.NewIndex("")
	s := ss.GetInvertedIndex()
	Convey("Add", t, func() {
		So(s.Add("fieldName\0071", document.DocId(1)), ShouldBeNil)
		So(s.Add("fieldName\0071", document.DocId(5)), ShouldBeNil)
		So(s.Add("fieldName\0071", document.DocId(6)), ShouldBeNil)
		So(s.Add("fieldName\0071", document.DocId(7)), ShouldBeNil)

		So(s.Add("fieldName\0072", document.DocId(1)), ShouldBeNil)
		So(s.Add("fieldName\0072", document.DocId(5)), ShouldBeNil)
		So(s.Add("fieldName\0072", document.DocId(6)), ShouldBeNil)
		So(s.Add("fieldName\0072", document.DocId(7)), ShouldBeNil)

		So(s.Add("fieldName\0073", document.DocId(2)), ShouldBeNil)
		So(s.Add("fieldName\0073", document.DocId(7)), ShouldBeNil)

		So(s.Add("fieldName\0074", document.DocId(2)), ShouldBeNil)
		So(s.Add("fieldName\0074", document.DocId(7)), ShouldBeNil)
		b := NewNotAndQuery([]Query{
			NewTermQuery(s.Iterator("fieldName", "1")),
			NewTermQuery(s.Iterator("fieldName", "2")),
			NewTermQuery(s.Iterator("fieldName", "3")),
			NewTermQuery(s.Iterator("fieldName", "4")),
		}, nil)
		tq := b
		res := tq.Marshal()
		fmt.Println(res)

		tq.SetDebug(1)
		//fmt.Println(tq.Next())
		//fmt.Println(tq.Next())
		//fmt.Println(tq.Next())
		//fmt.Println(tq.Next())
		fmt.Println(s.GetValueById(document.DocId(5)))
		fmt.Println(s.GetValueById(document.DocId(6)))
		fmt.Println(s.GetValueById(document.DocId(7)))
		fmt.Println(tq.DebugInfo())
	})
}

func TestNewTermQuery1(t *testing.T) {
	ss := index.NewIndex("")
	s1 := ss.GetInvertedIndex()
	s2 := ss.GetStorageIndex()
	Convey("Add", t, func() {
		So(s1.Add("fieldName\0071", 1), ShouldBeNil)
		So(s1.Add("fieldName\0071", 3), ShouldBeNil)
		So(s1.Add("fieldName\0071", 4), ShouldBeNil)
		So(s1.Add("fieldName\0071", 6), ShouldBeNil)
		So(s1.Add("fieldName\0071", 10), ShouldBeNil)

		So(s1.Add("fieldName\0072", 3), ShouldBeNil)
		So(s1.Add("fieldName\0072", 1), ShouldBeNil)
		So(s1.Add("fieldName\0072", 4), ShouldBeNil)
		So(s1.Add("fieldName\0072", 6), ShouldBeNil)

		So(s1.Add("fieldName\0073", 3), ShouldBeNil)
		So(s1.Add("fieldName\0073", 1), ShouldBeNil)
		So(s1.Add("fieldName\0073", 4), ShouldBeNil)

		So(s1.Add("fieldName\0074", 4), ShouldBeNil)
		So(s1.Add("fieldName\0074", 1), ShouldBeNil)
		So(s1.Add("fieldName\0074", 6), ShouldBeNil)

		So(s1.Add("fieldName\0075", 3), ShouldBeNil)
		So(s1.Add("fieldName\0075", 1), ShouldBeNil)
		So(s1.Add("fieldName\0075", 4), ShouldBeNil)

		So(s1.Add("fieldName\0076", 4), ShouldBeNil)
		So(s1.Add("fieldName\0076", 1), ShouldBeNil)
		So(s1.Add("fieldName\0076", 6), ShouldBeNil)

		So(s2.Add("fieldName", 3, 3), ShouldBeNil)
		So(s2.Add("fieldName", 4, 3), ShouldBeNil)
		So(s2.Add("fieldName", 6, 3), ShouldBeNil)
		So(s2.Add("fieldName", 1, 3), ShouldBeNil)
		So(s2.Add("fieldName", 10, 3), ShouldBeNil)

		q := NewAndQuery([]Query{
			NewTermQuery(s1.Iterator("fieldName", "1")),
			NewTermQuery(s1.Iterator("fieldName", "2")),
			NewTermQuery(s1.Iterator("fieldName", "3")),
		}, []check.Checker{
			check.NewChecker(s2.Iterator("fieldName"), 3, operation.EQ, nil, false),
			check.NewChecker(s2.Iterator("fieldName"), 3, operation.EQ, nil, false),
			check.NewAndChecker([]check.Checker{
				check.NewChecker(s2.Iterator("fieldName"), 3, operation.EQ, nil, false),
				check.NewChecker(s2.Iterator("fieldName"), 3, operation.EQ, nil, false),
			}),
		})

		q.SetDebug(1)

		fmt.Println(q.Current())
		//fmt.Println(q.Next())
		//fmt.Println(q.Next())
		//fmt.Println(q.Next())
		//fmt.Println(q.Next())

		fmt.Println(q.DebugInfo())

		res := q.Marshal() // query marshal params: index
		fmt.Println(res)
		//r, \007 := json.Marshal(res)
		//fmt.Println(string(r))
		//rr := q.Unmarshal(ss, rr1, nil) // unmarshal query  params:   1. index   2. query marshal结果  3. operation

		rr := q.Unmarshal(ss, res, nil)
		fmt.Println(rr.Current())
		//fmt.Println(rr.Next())
		//fmt.Println(rr.Next())
		//fmt.Println(rr.Next())
		//fmt.Println(rr.Next())
		fmt.Println(rr.DebugInfo())
	})
}

func TestNewExpression(t *testing.T) {
	var a = [][]string{{"field:fieldName\0076", "reason: found id"}}
	var b = [][]string{{"field:fieldName\0076", "reason: found id"}}
	fmt.Println(reflect.DeepEqual(a[0], b[0]))
	fmt.Println(helpers.CompareSlice(a, b))
	r := "hello" + "\007" + "world"
	fmt.Println(len(r))
	fmt.Println("----" + "\007" + "---")
	fmt.Println("----" + " " + "---")
	fmt.Println(r[5])
	fmt.Println(strings.Split(r, "\007"))
}
