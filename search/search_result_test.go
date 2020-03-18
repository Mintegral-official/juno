package search

import (
	"encoding/json"
	"fmt"
	"github.com/MintegralTech/juno/check"
	"github.com/MintegralTech/juno/document"
	"github.com/MintegralTech/juno/index"
	"github.com/MintegralTech/juno/operation"
	"github.com/MintegralTech/juno/query"
	. "github.com/smartystreets/goconvey/convey"
	"reflect"
	"testing"
)

func TestSearcher(t *testing.T) {
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
	Convey("Search Test", t, func() {
		_ = ss.Add(doc1)
		_ = ss.Add(doc2)
		_ = ss.Add(doc3)
		_ = ss.Add(doc4)
		_ = ss.Add(doc5)
		s1 := ss.GetInvertedIndex()
		s2 := ss.GetStorageIndex()
		q1 := query.NewAndQuery([]query.Query{
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
		})
		q := query.NewAndQuery([]query.Query{
			query.NewTermQuery(s1.Iterator("fieldName", "1")),
			query.NewTermQuery(s1.Iterator("fieldNeme", "2")),
			query.NewTermQuery(s1.Iterator("fieldName", "3")),
			q1,
		}, []check.Checker{
			check.NewInChecker(s2.Iterator("fieldName"), []int{2, 3, 4}, nil, false),
			check.NewOrChecker([]check.Checker{
				check.NewChecker(s2.Iterator("fieldName"), 2, operation.GT, nil, false),
				check.NewChecker(s2.Iterator("fieldName"), 3, operation.EQ, nil, false),
			}),
		})

		se := Search(ss, q)
		testCase := []document.DocId{10, 30, 40, 60}
		for i, expect := range testCase {
			So(se.Docs[i], ShouldEqual, expect)
		}
		v, e := q.Current()
		q.Next()
		So(v, ShouldEqual, 0)
		So(e, ShouldNotBeNil)
		res, _ := json.Marshal(q.Marshal())
		So(res, ShouldNotBeNil)

		q.SetLabel("test label q")
		q1.SetLabel("test label q1")
		a := Replay(ss, q.Marshal(), []document.DocId{1, 10})
		res, _ = json.Marshal(a)
		fmt.Println(string(res))
	})

}

func TestNewSearcher_Inc_Index(t *testing.T) {

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
	Convey("search inc index", t, func() {
		idx := index.NewIndex("")
		_ = idx.Add(doc4)
		q := query.NewTermQuery(idx.GetInvertedIndex().Iterator("field2", "2"))
		expectMap := &index.IndexDebugInfo{
			InvertIndex: map[string][]string{
				"field2": []string{"2"},
			},
			StorageIndex: map[string][]string{
				"field1": []string{"1"},
			},
		}
		realMap := idx.GetIndexDebugInfoById(0)
		So(reflect.DeepEqual(realMap, expectMap), ShouldBeTrue)

		s1 := Search(idx, q)
		So(s1.Docs[0], ShouldEqual, 0)

		idx.Del(doc5)
		_ = idx.Add(doc5)
		expectMap = &index.IndexDebugInfo{
			InvertIndex: map[string][]string{
				"field2": []string{"20", "200"},
			},
			StorageIndex: map[string][]string{
				"field1": []string{"10"},
			},
		}
		realMap = idx.GetIndexDebugInfoById(0)
		So(realMap, ShouldNotBeNil)
		So(reflect.DeepEqual(realMap, expectMap), ShouldBeTrue)

		q = query.NewTermQuery(idx.GetInvertedIndex().Iterator("field2", "20"))
		s1 = Search(idx, q)
		So(s1.Docs[0], ShouldEqual, 0)
	})

}

func TestIndexV2(t *testing.T) {
	var doc1 = &document.DocInfo{
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

	Convey("testIndex", t, func() {
		idx1 := index.NewIndexV2("abc")
		So(idx1.Add(doc1), ShouldBeNil)
		So(idx1.Add(doc2), ShouldBeNil)
		q := query.NewAndQuery([]query.Query{
			query.NewTermQuery(idx1.GetInvertedIndex().Iterator("fieldName", "3")),
		}, nil)
		s := Search(idx1, q)
		So(len(s.Docs), ShouldEqual, 1)
		So(s.Docs[0], ShouldEqual, 30)

		Convey("testMerge", func() {
			idx2 := index.NewIndexV2("merged")
			So(idx2.MergeIndex(idx1), ShouldBeNil)
			q := query.NewAndQuery([]query.Query{
				query.NewTermQuery(idx1.GetInvertedIndex().Iterator("fieldName", "3")),
			}, nil)
			s = Search(idx2, q)
			So(len(s.Docs), ShouldEqual, 1)
			So(s.Docs[0], ShouldEqual, 30)
		})

		Convey("test debug", func() {
			q := query.NewAndQuery([]query.Query{
				query.NewTermQuery(idx1.GetInvertedIndex().Iterator("fieldName", "3")),
			}, nil)
			filter := Replay(idx1, q.Marshal(), []document.DocId{30})
			fmt.Println(filter)
			fmt.Println("xxxxxxxxxxxxxxxxxxxxxxxxxx")
			fmt.Println(q.Marshal())
			fmt.Println("111111111111111111", idx1.GetIndexDebugInfoById(0))
		})

	})

}
