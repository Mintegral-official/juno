package search

import (
	"encoding/json"
	"fmt"
	"github.com/Mintegral-official/juno/check"
	"github.com/Mintegral-official/juno/document"
	"github.com/Mintegral-official/juno/operation"

	//"github.com/Mintegral-official/juno/document"
	"github.com/Mintegral-official/juno/index"
	"github.com/Mintegral-official/juno/query"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

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

		So(s1.Add("fieldNeme\0072", 3), ShouldBeNil)
		So(s1.Add("fieldNeme\0072", 1), ShouldBeNil)
		So(s1.Add("fieldNeme\0072", 4), ShouldBeNil)
		So(s1.Add("fieldNeme\0072", 6), ShouldBeNil)

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

		q := query.NewAndQuery([]query.Query{
			query.NewTermQuery(s1.Iterator("fieldName", "1")),
			query.NewTermQuery(s1.Iterator("fieldNeme", "2")),
			query.NewTermQuery(s1.Iterator("fieldName", "3")),
			//query.NewTermQuery(s1.Iterator("AAAAAAAA", "s")),
			query.NewAndQuery([]query.Query{
				query.NewTermQuery(s1.Iterator("fieldName", "1")),
				query.NewTermQuery(s1.Iterator("fieldNeme", "2")),
				query.NewTermQuery(s1.Iterator("fieldName", "3")),
				query.NewOrQuery([]query.Query{
					query.NewTermQuery(s1.Iterator("fieldName", "1")),
					query.NewTermQuery(s1.Iterator("fieldNeme", "2")),
					query.NewTermQuery(s1.Iterator("fieldName", "3")),
					query.NewTermQuery(s1.Iterator("fieldName", "4")),
				}, []check.Checker{
					check.NewChecker(s2.Iterator("fieldName"), 2, operation.EQ, nil, false),
					check.NewAndChecker([]check.Checker{
						check.NewChecker(s2.Iterator("fieldName"), 2, operation.EQ, nil, false),
						check.NewChecker(s2.Iterator("fieldName"), 3, operation.EQ, nil, false),
					}),
				}),
			}, []check.Checker{
				check.NewInChecker(s2.Iterator("fieldName"), []int{2, 3, 4, 5}, nil, false),
				check.NewOrChecker([]check.Checker{
					check.NewChecker(s2.Iterator("fieldName"), 2, operation.NE, nil, false),
					check.NewChecker(s2.Iterator("fieldName"), 3, operation.EQ, nil, false),
				}),
			}),
		}, []check.Checker{
			check.NewInChecker(s2.Iterator("fieldName"), []int{2, 3, 4}, nil, false),
			check.NewOrChecker([]check.Checker{
				check.NewChecker(s2.Iterator("fieldName"), 2, operation.GT, nil, false),
				check.NewChecker(s2.Iterator("fieldName"), 3, operation.EQ, nil, false),
			}),
		})
		fmt.Println(q.Current())
		q.Next()
		fmt.Println(q.Current())
		q.Next()
		fmt.Println(q.Current())
		q.Next()
		fmt.Println(q.Current())
		q.Next()
		fmt.Println(q.Current())
		q.Next()
		fmt.Println(q.Current())
		q.Next()
		r, _ := json.Marshal(q.Marshal())
		fmt.Println(string(r))
		sq := NewSearcher()
		sq.Debug(ss, q.Marshal(), nil, []document.DocId{document.DocId(10)})
		bb, _ := json.Marshal(sq.FilterInfo)
		fmt.Println(string(bb))

		//sea := NewSearcher()
		//sea.Debug(ss, q.Marshal(), nil, []document.DocId{document.DocId(10)})
		//bbb, _ := json.Marshal(sea.FilterInfo)
		//fmt.Println(string(bbb))
		//fmt.Println(sea.Docs)
		//
		//sea.Search(ss, q)
		//fmt.Println(sea.Docs)
		//
		//fmt.Println(ss.GetValueById(document.DocId(10)))

		qq := query.NewAndQuery([]query.Query{
			query.NewTermQuery(s1.Iterator("fieldName", "1")),
			query.NewTermQuery(s1.Iterator("fieldNeme", "2")),
			query.NewTermQuery(s1.Iterator("fieldName", "3")),
			query.NewNotAndQuery([]query.Query{
				query.NewTermQuery(s1.Iterator("fieldName", "1")),
				query.NewTermQuery(s1.Iterator("fieldNeme", "2")),
				query.NewTermQuery(s1.Iterator("fieldName", "3")),
				query.NewTermQuery(s1.Iterator("fieldName", "4")),
			}, nil),
		}, nil)

		sea := NewSearcher()
		sea.Search(ss, qq)
		fmt.Println(sea.Docs)
		sea = NewSearcher()

		sea.Debug(ss, qq.Marshal(), nil, []document.DocId{document.DocId(10)})
		bbb, _ := json.Marshal(sea.FilterInfo)
		fmt.Println(string(bbb))
		fmt.Println(sea.Docs)

	})
}

func TestSearcher_Debug(t *testing.T) {
	ss := index.NewIndex("")
	s1 := ss.GetInvertedIndex()
	s2 := ss.GetStorageIndex()
	Convey("Add", t, func() {
		So(s1.Add("fieldName\0071", 1), ShouldBeNil)
		So(s1.Add("fieldName\0071", 3), ShouldBeNil)
		So(s1.Add("fieldName\0071", 4), ShouldBeNil)
		So(s1.Add("fieldName\0071", 6), ShouldBeNil)
		So(s1.Add("fieldName\0071", 10), ShouldBeNil)

		So(s1.Add("fieldNeme\0072", 3), ShouldBeNil)
		So(s1.Add("fieldNeme\0072", 1), ShouldBeNil)
		So(s1.Add("fieldNeme\0072", 4), ShouldBeNil)
		So(s1.Add("fieldNeme\0072", 6), ShouldBeNil)

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

		q := query.NewAndQuery([]query.Query{
			query.NewTermQuery(s1.Iterator("fieldName", "1")),
			query.NewTermQuery(s1.Iterator("fieldNeme", "2")),
			query.NewTermQuery(s1.Iterator("fieldName", "3")),
			query.NewAndQuery([]query.Query{
				query.NewTermQuery(s1.Iterator("fieldName", "1")),
				query.NewTermQuery(s1.Iterator("fieldNeme", "2")),
				query.NewTermQuery(s1.Iterator("fieldName", "3")),
				query.NewOrQuery([]query.Query{
					query.NewTermQuery(s1.Iterator("fieldName", "1")),
					query.NewTermQuery(s1.Iterator("fieldNeme", "2")),
					query.NewTermQuery(s1.Iterator("fieldName", "3")),
					query.NewTermQuery(s1.Iterator("fieldName", "4")),
				}, []check.Checker{
					check.NewChecker(s2.Iterator("fieldName"), 2, operation.NE, nil, false),
					check.NewAndChecker([]check.Checker{
						check.NewChecker(s2.Iterator("fieldName"), 2, operation.EQ, nil, false),
						check.NewChecker(s2.Iterator("fieldName"), 3, operation.EQ, nil, false),
					}),
				}),
			}, []check.Checker{
				check.NewInChecker(s2.Iterator("fieldName"), []int{2, 3, 4, 5}, nil, false),
				check.NewOrChecker([]check.Checker{
					check.NewChecker(s2.Iterator("fieldName"), 2, operation.NE, nil, false),
					check.NewChecker(s2.Iterator("fieldName"), 3, operation.EQ, nil, false),
				}),
			}),
		}, []check.Checker{
			check.NewInChecker(s2.Iterator("fieldName"), []int{2, 3, 4}, nil, false),
			check.NewOrChecker([]check.Checker{
				check.NewChecker(s2.Iterator("fieldName"), 2, operation.GT, nil, false),
				check.NewChecker(s2.Iterator("fieldName"), 3, operation.EQ, nil, false),
			}),
		})
		q.SetDebug(1)
		testCase := []document.DocId{1, 3, 4}
		for _, expect := range testCase {
			v, e := q.Current()
			q.Next()
			So(v, ShouldEqual, expect)
			So(e, ShouldBeNil)
		}
		v, e := q.Current()
		q.Next()
		So(v, ShouldEqual, 0)
		So(e, ShouldNotBeNil)
		fmt.Println(q.DebugInfo().String())
	})

}
