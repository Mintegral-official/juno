package query

import (
	"github.com/MintegralTech/juno/document"
	"github.com/MintegralTech/juno/index"
	"testing"
)

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

//func TestSqlQuery_LRD(t *testing.T) {
//	s := "field1>= 1 and   (field2 !=1 or ( field2 = 1 and field1=1)) | (field1 in [1,2] and field1 !in [2,3])"
//	sq := NewSqlQuery(s, nil, false)
//	Convey("sql query", t, func() {
//		node := sq.exp2Tree()
//		n := node.To()
//		//node.Print()
//		So(n.Len(), ShouldEqual, 11)
//		idx := index.NewIndex("index")
//		So(idx.Add(doc1), ShouldBeNil)
//		So(idx.Add(doc2), ShouldBeNil)
//		So(idx.Add(doc3), ShouldBeNil)
//		q := sq.LRD(idx)
//		if _, err := q.Current(); err != nil {
//			So(err, ShouldNotBeNil)
//		}
//		id, err := q.Next()
//		for err == nil {
//			So(id, ShouldNotEqual, 0)
//			id, err = q.Next()
//		}
//	})
//}

func BenchmarkSqlQuery_LRD(b *testing.B) {
	s := "field1>= 1 and   (field2 !=1 or ( field2 = 1 and field1=1)) | (field1 in [1,2] and field1 !in [2,3])"
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sq := NewSqlQuery(s, nil, false)
		idx := index.NewIndex("index")
		_ = idx.Add(doc1)
		_ = idx.Add(doc2)
		_ = idx.Add(doc3)
		_ = sq.LRD(idx)
	}
}
