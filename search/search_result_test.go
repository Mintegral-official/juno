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
		sq1 := NewSearcher()
		sq1.Search(ss, q)
		fmt.Println(sq1.Docs)
		r, _ := json.Marshal(q.Marshal())
		fmt.Println(string(r))
		sq := NewSearcher()
		sq.Debug(ss, q.Marshal(), nil, []document.DocId{document.DocId(10)})
		bb, _ := json.Marshal(sq.FilterInfo)
		fmt.Println(string(bb))
		fmt.Println(sq.Docs)

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
	var doc1 = &document.DocInfo{
		Id: 10,
		Fields: []*document.Field{
			{
				Name:      "fieldName",
				IndexType: 0,
				Value:     "1",
				ValueType: document.StringFieldType,
			},
			{
				Name:      "fieldNeme",
				IndexType: 0,
				Value:     "2",
				ValueType: document.StringFieldType,
			},
			{
				Name:      "fieldName",
				IndexType: 0,
				Value:     "3",
				ValueType: document.StringFieldType,
			},
			{
				Name:      "fieldName",
				IndexType: 0,
				Value:     "4",
				ValueType: document.StringFieldType,
			},
			{
				Name:      "fieldName",
				IndexType: 0,
				Value:     "5",
				ValueType: document.StringFieldType,
			},
			{
				Name:      "fieldName",
				IndexType: 0,
				Value:     "6",
				ValueType: document.StringFieldType,
			},
			{
				Name:      "fieldName",
				IndexType: 1,
				Value:     3,
				ValueType: document.StringFieldType,
			},
		},
	}
	var doc2 = &document.DocInfo{
		Id: 30,
		Fields: []*document.Field{
			{
				Name:      "fieldName",
				IndexType: 0,
				Value:     "1",
				ValueType: document.StringFieldType,
			},
			{
				Name:      "fieldNeme",
				IndexType: 0,
				Value:     "2",
				ValueType: document.StringFieldType,
			},
			{
				Name:      "fieldName",
				IndexType: 0,
				Value:     "3",
				ValueType: document.StringFieldType,
			},
			{
				Name:      "fieldName",
				IndexType: 0,
				Value:     "5",
				ValueType: document.StringFieldType,
			},
			{
				Name:      "fieldName",
				IndexType: 1,
				Value:     3,
				ValueType: document.StringFieldType,
			},
		},
	}
	var doc3 = &document.DocInfo{
		Id: 40,
		Fields: []*document.Field{
			{
				Name:      "fieldName",
				IndexType: 0,
				Value:     "1",
				ValueType: document.StringFieldType,
			},
			{
				Name:      "fieldNeme",
				IndexType: 0,
				Value:     "2",
				ValueType: document.StringFieldType,
			},
			{
				Name:      "fieldName",
				IndexType: 0,
				Value:     "3",
				ValueType: document.StringFieldType,
			},
			{
				Name:      "fieldName",
				IndexType: 0,
				Value:     "4",
				ValueType: document.StringFieldType,
			},
			{
				Name:      "fieldName",
				IndexType: 0,
				Value:     "5",
				ValueType: document.StringFieldType,
			},
			{
				Name:      "fieldName",
				IndexType: 0,
				Value:     "6",
				ValueType: document.StringFieldType,
			},
			{
				Name:      "fieldName",
				IndexType: 1,
				Value:     3,
				ValueType: document.StringFieldType,
			},
		},
	}
	var doc4 = &document.DocInfo{
		Id: 60,
		Fields: []*document.Field{
			{
				Name:      "fieldName",
				IndexType: 0,
				Value:     "1",
				ValueType: document.StringFieldType,
			},
			{
				Name:      "fieldNeme",
				IndexType: 0,
				Value:     "2",
				ValueType: document.StringFieldType,
			},
			{
				Name:      "fieldName",
				IndexType: 0,
				Value:     "3",
				ValueType: document.StringFieldType,
			},
			{
				Name:      "fieldName",
				IndexType: 0,
				Value:     "4",
				ValueType: document.StringFieldType,
			},
			{
				Name:      "fieldName",
				IndexType: 0,
				Value:     "6",
				ValueType: document.StringFieldType,
			},
			{
				Name:      "fieldName",
				IndexType: 1,
				Value:     3,
				ValueType: document.StringFieldType,
			},
		},
	}
	var doc5 = &document.DocInfo{
		Id: 100,
		Fields: []*document.Field{
			{
				Name:      "fieldName",
				IndexType: 0,
				Value:     "1",
				ValueType: document.StringFieldType,
			},
			{
				Name:      "fieldName",
				IndexType: 1,
				Value:     3,
				ValueType: document.StringFieldType,
			},
		},
	}
	ss := index.NewIndex("")
	Convey("SetDebug", t, func() {
		//So(s1.Add("fieldName\0071", 10), ShouldBeNil)
		//So(s1.Add("fieldName\0071", 30), ShouldBeNil)
		//So(s1.Add("fieldName\0071", 40), ShouldBeNil)
		//So(s1.Add("fieldName\0071", 60), ShouldBeNil)
		//So(s1.Add("fieldName\0071", 100), ShouldBeNil)
		//
		//So(s1.Add("fieldNeme\0072", 30), ShouldBeNil)
		//So(s1.Add("fieldNeme\0072", 10), ShouldBeNil)
		//So(s1.Add("fieldNeme\0072", 40), ShouldBeNil)
		//So(s1.Add("fieldNeme\0072", 60), ShouldBeNil)
		//
		//So(s1.Add("fieldName\0073", 30), ShouldBeNil)
		//So(s1.Add("fieldName\0073", 10), ShouldBeNil)
		//So(s1.Add("fieldName\0073", 40), ShouldBeNil)
		//
		//So(s1.Add("fieldName\0074", 40), ShouldBeNil)
		//So(s1.Add("fieldName\0074", 10), ShouldBeNil)
		//So(s1.Add("fieldName\0074", 60), ShouldBeNil)
		//
		//So(s1.Add("fieldName\0075", 30), ShouldBeNil)
		//So(s1.Add("fieldName\0075", 10), ShouldBeNil)
		//So(s1.Add("fieldName\0075", 40), ShouldBeNil)
		//
		//So(s1.Add("fieldName\0076", 40), ShouldBeNil)
		//So(s1.Add("fieldName\0076", 10), ShouldBeNil)
		//So(s1.Add("fieldName\0076", 60), ShouldBeNil)
		//
		//So(s2.Add("fieldName", 30, 3), ShouldBeNil)
		//So(s2.Add("fieldName", 40, 3), ShouldBeNil)
		//So(s2.Add("fieldName", 60, 3), ShouldBeNil)
		//So(s2.Add("fieldName", 10, 3), ShouldBeNil)
		//So(s2.Add("fieldName", 100, 3), ShouldBeNil)
		ss.Add(doc1)
		ss.Add(doc2)
		ss.Add(doc3)
		ss.Add(doc4)
		ss.Add(doc5)
		s1 := ss.GetInvertedIndex()
		s2 := ss.GetStorageIndex()
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

		//se := NewSearcher()
		//se.Search(ss, q)
		//fmt.Println(se.Docs)
		testCase := []document.DocId{10, 30, 40, 60}
		for _, expect := range testCase {
			v, e := q.Current()
			q.Next()
			vb, ok := ss.GetBitMap().Get(index.DocId(v))
			if ok {
				So(vb, ShouldEqual, expect)
				So(e, ShouldBeNil)
			}
		}
		v, e := q.Current()
		q.Next()
		So(v, ShouldEqual, 0)
		So(e, ShouldNotBeNil)
		fmt.Println(q.DebugInfo().String())

	})

}

func TestNewSearcher(t *testing.T) {

	var doc4 = &document.DocInfo{
		Id: 0,
		Fields: []*document.Field{
			{
				Name:      "field1",
				IndexType: 1,
				Value:     1,
				ValueType: document.IntFieldType,
			},
			{
				Name:      "field2",
				IndexType: 0,
				Value:     "2",
				ValueType: document.StringFieldType,
			},
		},
	}

	var doc5 = &document.DocInfo{
		Id: 0,
		Fields: []*document.Field{
			{
				Name:      "field1",
				IndexType: 1,
				Value:     10,
				ValueType: document.IntFieldType,
			},
			{
				Name:      "field2",
				IndexType: 0,
				Value:     "20",
				ValueType: document.StringFieldType,
			},
			{
				Name:      "field2",
				IndexType: 0,
				Value:     "200",
				ValueType: document.StringFieldType,
			},
		},
	}

	idx := index.NewIndex("")
	_ = idx.Add(doc4)
	q := query.NewTermQuery(idx.GetInvertedIndex().Iterator("field2", "2"))
	fmt.Println(idx.GetValueById(0))

	s1 := NewSearcher()
	s1.Search(idx, q)
	fmt.Println(s1.Docs)

	idx.Del(doc5)
	_ = idx.Add(doc5)
	fmt.Println(idx.GetValueById(0))

	q = query.NewTermQuery(idx.GetInvertedIndex().Iterator("field2", "20"))
	s1 = NewSearcher()
	s1.Search(idx, q)
	fmt.Println(s1.Docs)
}
