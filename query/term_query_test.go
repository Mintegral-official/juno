package query

import (
	"fmt"
	"github.com/Mintegral-official/juno/check"
	"github.com/Mintegral-official/juno/document"
	"github.com/Mintegral-official/juno/index"
	"github.com/Mintegral-official/juno/operation"
	. "github.com/smartystreets/goconvey/convey"
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
		So(s.Del("fieldName", 1), ShouldBeFalse)
		So(s.Del("fieldName_1", document.DocId(1)), ShouldBeTrue)
		b := s.Iterator("fieldName", "1")
		tq := NewTermQuery(b)
		res := tq.Marshal(ss)
		fmt.Println(res)
		sss := tq.Unmarshal(ss, tq.Marshal(ss), nil)
		fmt.Println(sss.Current())
		fmt.Println(sss.Next())
		fmt.Println(sss.Next())
		fmt.Println(sss.Next())
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

func TestNewTermQuery1(t *testing.T) {
	ss := index.NewIndex("")
	s1 := ss.GetInvertedIndex()
	s2 := ss.GetStorageIndex()
	Convey("Add", t, func() {
		So(s1.Add("fieldName_1", 1), ShouldBeNil)
		So(s1.Add("fieldName_1", 3), ShouldBeNil)
		So(s1.Add("fieldName_1", 4), ShouldBeNil)
		So(s1.Add("fieldName_1", 6), ShouldBeNil)
		So(s1.Add("fieldName_2", 3), ShouldBeNil)
		So(s1.Add("fieldName_2", 4), ShouldBeNil)
		So(s1.Add("fieldName_2", 6), ShouldBeNil)
		So(s1.Add("fieldName_1", 10), ShouldBeNil)
		So(s2.Add("fieldName", 3, 3), ShouldBeNil)
		So(s2.Add("fieldName", 4, 3), ShouldBeNil)
		So(s2.Add("fieldName", 6, 3), ShouldBeNil)
		So(s2.Add("fieldName", 8, 4), ShouldBeNil)
		So(s2.Add("fieldName", 10, 3), ShouldBeNil)

		q := NewOrQuery([]Query{
			NewTermQuery(s1.Iterator("fieldName", "1")),
			NewAndQuery([]Query{
				NewTermQuery(s1.Iterator("fieldName", "1")),
				NewTermQuery(s1.Iterator("fieldName", "2")),
			}, nil),
		}, []check.Checker{
			check.NewChecker(s2.Iterator("fieldName"), 3, operation.EQ, nil, false),
		})

		fmt.Println(q.Current())
		fmt.Println(q.Next())
		fmt.Println(q.Next())
		fmt.Println(q.Next())
		fmt.Println(q.Next())

		res := q.Marshal(ss)            // query marshal params: index
		rr := q.Unmarshal(ss, res, nil) // unmarshal query  params:   1. index   2. query marshal结果  3. operation
		fmt.Println(rr.Current())
		fmt.Println(rr.Next())
		fmt.Println(rr.Next())
		fmt.Println(rr.Next())
		fmt.Println(rr.Next())
	})
}
