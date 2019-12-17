package query

import (
	"github.com/Mintegral-official/juno/document"
	. "github.com/smartystreets/goconvey/convey"
	"strings"
	"testing"
)

var doc1 = &document.DocInfo{
	Id: 0,
	Fields: []*document.Field{
		{
			Name:      "field1",
			IndexType: 1,
			Value:     1,
		},
		{
			Name:      "field2",
			IndexType: 0,
			Value:     2,
		},
		{
			Name:      "field1",
			IndexType: 2,
			Value:     1,
		},
	},
}

var doc2 = &document.DocInfo{
	Id: 1,
	Fields: []*document.Field{
		{
			Name:      "field1",
			IndexType: 0,
			Value:     1,
		},
		{
			Name:      "field2",
			IndexType: 1,
			Value:     2,
		},
		{
			Name:      "field1",
			IndexType: 0,
			Value:     1,
		},
	},
}

var doc3 = &document.DocInfo{
	Id: 2,
	Fields: []*document.Field{
		{
			Name:      "field1",
			IndexType: 0,
			Value:     1,
		},
		{
			Name:      "field2",
			IndexType: 0,
			Value:     2,
		},
		{
			Name:      "field1",
			IndexType: 1,
			Value:     1,
		},
	},
}

func TestNewSqlQuery(t *testing.T) {
	s := "country= CN &   (a =1 | ( b = 1 & a!=0)) | (c in [1,2,3] & d not in [2,4,5])"
	s = strings.Replace(strings.Replace(s, " not in ", " # ", -1), " in ", "@", -1)
	sq := NewSqlQuery(s)
	Convey("sql query", t, func() {
		node := sq.exp2Tree()
		node.Print()
	})
}

//func TestSqlQuery_LRD(t *testing.T) {
//	s := "field1> 1 &   (filed2 !=1 | ( filed2 = 1 & field1=1)) | (filed1 @ [1,2] & field1 # [2,3])"
//	sq := NewSqlQuery(s)
//	Convey("sql query", t, func() {
//		node := sq.exp2Tree()
//		node.Print()
//		idx := index.NewIndex("index")
//		So(idx.Add(doc1), ShouldBeNil)
//		So(idx.Add(doc2), ShouldBeNil)
//		So(idx.Add(doc3), ShouldBeNil)
//		q := sq.LRD(idx)
//		fmt.Println(q)
//		//if _, err := q.Current(); err != nil {
//		//	fmt.Println(err)
//		//}
//		//id, err := q.Next()
//		//for err == nil {
//		//	fmt.Println(id)
//		//	id, err = q.Next()
//		//}
//	})
//}
