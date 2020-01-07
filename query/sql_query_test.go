package query

import (
	"github.com/Mintegral-official/juno/document"
	"github.com/Mintegral-official/juno/index"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

//func TestNewSqlQuery(t *testing.T) {
//	s := "country= CN &   (a =1 | ( b = 1 & a!=0)) | (c in [1,2,3] & d not in [2,4,5])"
//	s = strings.Replace(strings.Replace(s, " not in ", " # ", -1), " in ", "@", -1)
//	sq := NewSqlQuery(s)
//	Convey("sql query", t, func() {
//		node := sq.exp2Tree()
//		node.Print()
//	})
//}

func TestSqlQuery_LRD(t *testing.T) {
	var doc1 = &document.DocInfo{
		Id: 0,
		Fields: []*document.Field{
			{
				Name:      "field1",
				IndexType: 1,
				Value:     int64(1),
				ValueType: document.IntFieldType,
			},
			{
				Name:      "field2",
				IndexType: 0,
				Value:     "2",
				ValueType: document.StringFieldType,
			},
			{
				Name:      "field1",
				IndexType: 0,
				Value:     int64(1),
				ValueType: document.IntFieldType,
			},
		},
	}

	var doc2 = &document.DocInfo{
		Id: 1,
		Fields: []*document.Field{
			{
				Name:      "field1",
				IndexType: 0,
				Value:     int64(1),
				ValueType: document.IntFieldType,
			},
			{
				Name:      "field2",
				IndexType: 1,
				Value:     "2",
				ValueType: document.StringFieldType,
			},
			{
				Name:      "field1",
				IndexType: 0,
				Value:     int64(1),
				ValueType: document.IntFieldType,
			},
		},
	}

	var doc3 = &document.DocInfo{
		Id: 2,
		Fields: []*document.Field{
			{
				Name:      "field1",
				IndexType: 0,
				Value:     int64(1),
				ValueType: document.IntFieldType,
			},
			{
				Name:      "field2",
				IndexType: 0,
				Value:     "2",
				ValueType: document.StringFieldType,
			},
			{
				Name:      "field1",
				IndexType: 1,
				Value:     int64(1),
				ValueType: document.IntFieldType,
			},
		},
	}
	s := "field1>= 1 and   (field2 !=1 or ( field2 = 1 and field1=1)) | (field1 in [1,2] and field1 !in [2,3])"
	sq := NewSqlQuery(s)
	Convey("sql query", t, func() {
		node := sq.exp2Tree()
		n := node.To()
		//node.Print()
		So(n.Len(), ShouldEqual, 11)
		//fmt.Println(n.Len())
		idx := index.NewIndex("index")
		So(idx.Add(doc1), ShouldBeNil)
		So(idx.Add(doc2), ShouldBeNil)
		So(idx.Add(doc3), ShouldBeNil)
		q := sq.LRD(idx)
		if _, err := q.Current(); err != nil {
			So(err, ShouldNotBeNil)
			//	fmt.Println(err)
		}
		id, err := q.Next()
		for err == nil {
			So(id, ShouldNotEqual, 0)
			//fmt.Println(id)
			id, err = q.Next()
		}
		//So(idx.DebugInfo(), ShouldNotBeNil)
		//So(q.DebugInfo(), ShouldNotBeNil)
		//fmt.Println(idx.DebugInfo())
		//fmt.Println(q.DebugInfo())
	})
}

//func TestSqlQuery_Next(t *testing.T) {
//	Convey("sql", t, func() {
//		s := "advertiserId=457 | platform=1 | (price @ [20.0, 1.4, 3.6, 5.7, 2.5] & price >= 1.4)"
//		sq := NewSqlQuery(s)
//		sq.exp2Tree().Print()
//		idx := index.NewIndex("index")
//		_ = idx.Add(doca)
//		_ = idx.Add(docb)
//		_ = idx.Add(docc)
//		q := sq.LRD(idx)
//		if id, err := q.Current(); err != nil {
//			So(err, ShouldNotBeNil)
//			//fmt.Println(err)
//		} else {
//			//fmt.Println(id)
//			So(id, ShouldNotEqual, 0)
//		}
//		id, err := q.Next()
//		for err == nil {
//			fmt.Println(id)
//			id, err = q.Next()
//		}
//		So(idx.DebugInfo(), ShouldNotBeNil)
//		So(q.DebugInfo(), ShouldNotBeNil)
//		//fmt.Println(idx.DebugInfo())
//		//fmt.Println(q.DebugInfo())
//	})
//}

