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
	"testing"
)

func TestNewTermQuery(t *testing.T) {
	ss := index.NewIndex("")
	s := ss.GetInvertedIndex()
	Convey("Add", t, func() {
		So(s.Add("fieldName_1", document.DocId(1)), ShouldBeNil)
		So(s.Add("fieldName_1", document.DocId(5)), ShouldBeNil)
		So(s.Add("fieldName_1", document.DocId(6)), ShouldBeNil)
		So(s.Add("fieldName_1", document.DocId(7)), ShouldBeNil)
		So(s.Add("fieldName_4", document.DocId(2)), ShouldBeNil)
		b := s.Iterator("fieldName", "1")
		tq := NewTermQuery(b)
		res := tq.Marshal(ss)
		fmt.Println(res)

		tq.SetDebug(1)
		fmt.Println(tq.Current())
		fmt.Println(tq.Next())
		fmt.Println(tq.Next())
		fmt.Println(tq.Next())
		fmt.Println(tq.Next())
		fmt.Println(tq.DebugInfo())
		//jf := &JSONFormatter{}
		//str, _ := jf.Marshal(res) // 转换成json的形式
		//fmt.Println(str)
		//rr1, _ := jf.Unmarshal(str) // 反序列化
		//sss := tq.Unmarshal(ss, rr1, nil)

		sss := tq.Unmarshal(ss, res, nil)
		fmt.Println(sss.Current())
		fmt.Println(sss.Next())
		fmt.Println(sss.Next())
		fmt.Println(sss.Next())
		fmt.Println(sss.DebugInfo())
		//aq := NewAndQuery([]Query{
		//	NewTermQuery(s.Iterator("fieldName", "10")),
		//	NewTermQuery(s.Iterator("fieldName", "1")),
		//	NewTermQuery(s.Iterator("fieldName", "100")),
		//	NewOrQuery([]Query{
		//		NewTermQuery(s.Iterator("fieldName", "111")),
		//		NewTermQuery(s.Iterator("fieldName", "123")),
		//	}, nil),
		//	NewNotAndQuery([]Query{
		//		NewTermQuery(s.Iterator("fieldName", "111456")),
		//		NewTermQuery(s.Iterator("fieldName", "123111")),
		//	}, nil),
		//}, nil)
		////fmt.Println(aq.Marshal(ss))
		//res, _ := json.Marshal(aq.Marshal(ss))
		//fmt.Println(string(res))
	})
}

func TestTermQuery_Current(t *testing.T) {
	ss := index.NewIndex("")
	s := ss.GetInvertedIndex()
	Convey("Add", t, func() {
		So(s.Add("fieldName_1", document.DocId(1)), ShouldBeNil)
		So(s.Add("fieldName_1", document.DocId(5)), ShouldBeNil)
		So(s.Add("fieldName_1", document.DocId(6)), ShouldBeNil)
		So(s.Add("fieldName_1", document.DocId(7)), ShouldBeNil)

		So(s.Add("fieldName_2", document.DocId(1)), ShouldBeNil)
		So(s.Add("fieldName_2", document.DocId(5)), ShouldBeNil)
		So(s.Add("fieldName_2", document.DocId(6)), ShouldBeNil)
		So(s.Add("fieldName_2", document.DocId(7)), ShouldBeNil)

		So(s.Add("fieldName_3", document.DocId(2)), ShouldBeNil)
		So(s.Add("fieldName_3", document.DocId(7)), ShouldBeNil)

		So(s.Add("fieldName_4", document.DocId(2)), ShouldBeNil)
		So(s.Add("fieldName_4", document.DocId(7)), ShouldBeNil)
		b := NewNotAndQuery([]Query{
			NewTermQuery(s.Iterator("fieldName", "1")),
			NewTermQuery(s.Iterator("fieldName", "2")),
			NewTermQuery(s.Iterator("fieldName", "3")),
			NewTermQuery(s.Iterator("fieldName", "4")),
		}, nil)
		tq := b
		res := tq.Marshal(ss)
		fmt.Println(res)

		tq.SetDebug(1)
		fmt.Println(tq.Next())
		fmt.Println(tq.Next())
		fmt.Println(tq.Next())
		fmt.Println(tq.Next())
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
		So(s1.Add("fieldName_1", 1), ShouldBeNil)
		So(s1.Add("fieldName_1", 3), ShouldBeNil)
		So(s1.Add("fieldName_1", 4), ShouldBeNil)
		So(s1.Add("fieldName_1", 6), ShouldBeNil)
		So(s1.Add("fieldName_1", 10), ShouldBeNil)

		So(s1.Add("fieldName_2", 3), ShouldBeNil)
		So(s1.Add("fieldName_2", 1), ShouldBeNil)
		So(s1.Add("fieldName_2", 4), ShouldBeNil)
		So(s1.Add("fieldName_2", 6), ShouldBeNil)

		So(s1.Add("fieldName_3", 3), ShouldBeNil)
		So(s1.Add("fieldName_3", 1), ShouldBeNil)
		So(s1.Add("fieldName_3", 4), ShouldBeNil)

		So(s1.Add("fieldName_4", 4), ShouldBeNil)
		So(s1.Add("fieldName_4", 1), ShouldBeNil)
		So(s1.Add("fieldName_4", 6), ShouldBeNil)

		So(s1.Add("fieldName_5", 3), ShouldBeNil)
		So(s1.Add("fieldName_5", 1), ShouldBeNil)
		So(s1.Add("fieldName_5", 4), ShouldBeNil)

		So(s1.Add("fieldName_6", 4), ShouldBeNil)
		So(s1.Add("fieldName_6", 1), ShouldBeNil)
		So(s1.Add("fieldName_6", 6), ShouldBeNil)

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
		fmt.Println(q.Next())
		fmt.Println(q.Next())
		fmt.Println(q.Next())
		fmt.Println(q.Next())

		fmt.Println(q.DebugInfo())

		res := q.Marshal(ss) // query marshal params: index
		fmt.Println(res)
		//jf := &JSONFormatter{}
		//str, _ := jf.Marshal(res) // 转换成json的形式
		//fmt.Println("\n", str)
		//rr1, _ := jf.Unmarshal(str)     // 反序列化
		//rr := q.Unmarshal(ss, rr1, nil) // unmarshal query  params:   1. index   2. query marshal结果  3. operation

		rr := q.Unmarshal(ss, res, nil)
		fmt.Println(rr.Current())
		fmt.Println(rr.Next())
		fmt.Println(rr.Next())
		fmt.Println(rr.Next())
		fmt.Println(rr.Next())
		fmt.Println(rr.DebugInfo())
	})
}

func TestNewExpression(t *testing.T) {
	var a = [][]string{{"field:fieldName_6", "reason: found id"}}
	var b = [][]string{{"field:fieldName_6", "reason: found id"}}
	fmt.Println(reflect.DeepEqual(a[0], b[0]))
	fmt.Println(helpers.CompareSlice(a, b))
}