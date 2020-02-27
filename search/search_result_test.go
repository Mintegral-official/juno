package search

import (
	"fmt"
	"github.com/Mintegral-official/juno/check"
	"github.com/Mintegral-official/juno/index"
	"github.com/Mintegral-official/juno/operation"
	"github.com/Mintegral-official/juno/query"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

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

		q := query.NewAndQuery([]query.Query{
			query.NewTermQuery(s1.Iterator("fieldName", "1")),
			query.NewTermQuery(s1.Iterator("fieldName", "2")),
			query.NewTermQuery(s1.Iterator("fieldName", "3")),
		}, []check.Checker{
			check.NewChecker(s2.Iterator("fieldName"), 3, operation.EQ, nil, false),
			check.NewChecker(s2.Iterator("fieldName"), 3, operation.EQ, nil, false),
			check.NewAndChecker([]check.Checker{
				check.NewChecker(s2.Iterator("fieldName"), 3, operation.EQ, nil, false),
				check.NewChecker(s2.Iterator("fieldName"), 3, operation.EQ, nil, false),
			}),
		})

		q.SetDebug(1) // 设置debug调试信息

		se := NewSearcher()

		fmt.Println(se.Debug(ss, q).String()) // debug查询
		// se.DebugInfo(ss, q, ids)  ids：指定查找的id列表
	})
}
