package query

import (
	"fmt"
	"github.com/Mintegral-official/juno/document"
	"github.com/Mintegral-official/juno/index"
	. "github.com/smartystreets/goconvey/convey"
	"strings"
	"testing"
)

var doca = &document.DocInfo{
	Id: 0,
	Fields: []*document.Field{
		{
			Name:      "field1",
			IndexType: 2,
			Value:     1,
		},
		{
			Name:      "field2",
			IndexType: 2,
			Value:     1,
		},
		{
			Name:      "field1",
			IndexType: 2,
			Value:     2,
		},
	},
}

var docb = &document.DocInfo{
	Id: 1,
	Fields: []*document.Field{
		{
			Name:      "field1",
			IndexType: 2,
			Value:     3,
		},
		{
			Name:      "field2",
			IndexType: 2,
			Value:     2,
		},
		{
			Name:      "field1",
			IndexType: 2,
			Value:     4,
		},
	},
}

var docc = &document.DocInfo{
	Id: 2,
	Fields: []*document.Field{
		{
			Name:      "field1",
			IndexType: 2,
			Value:     5,
		},
		{
			Name:      "field2",
			IndexType: 2,
			Value:     3,
		},
		{
			Name:      "field1",
			IndexType: 2,
			Value:     6,
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

func TestSqlQuery_LRD(t *testing.T) {
	s := "field1>= 1 &   (filed2 !=1 | ( filed2 = 1 & field1=1)) | (filed1 @ [1,2] & field1 # [2,3])"
	sq := NewSqlQuery(s)
	Convey("sql query", t, func() {
		node := sq.exp2Tree()
		//node.Print()
		n := node.To()
		fmt.Println(n.Len())
		//for !n.Empty() {
		//	fmt.Println(n.Pop())
		//}
		idx := index.NewIndex("index")
		So(idx.Add(doca), ShouldBeNil)
		So(idx.Add(docb), ShouldBeNil)
		So(idx.Add(docc), ShouldBeNil)
		q := sq.LRD(idx)

		if _, err := q.Current(); err != nil {
			fmt.Println(err)
		}
		id, err := q.Next()
		for err == nil {
			fmt.Println(id)
			id, err = q.Next()
		}
		fmt.Println(idx.String())
		fmt.Println(q.String())
	})
}

func TestNotAndQuery_Next(t *testing.T) {
	s := "advertiserId=457 | platform=1 | (price @ [20.0, 1.4, 3.6, 5.7, 2.5] & price >= 1.4)"
	sq := NewSqlQuery(s)
	sq.exp2Tree().Print()
	idx := index.NewIndex("index")
	_ = idx.Add(doca)
	_ = idx.Add(docb)
	_ = idx.Add(docc)
	q := sq.LRD(idx)
	if id, err := q.Current(); err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(id)
	}
	id, err := q.Next()
	for err == nil {
		fmt.Println(id)
		id, err = q.Next()
	}
	fmt.Println(idx.String())
	fmt.Println(q.String())
}
