package index

import (
	"fmt"
	"github.com/Mintegral-official/juno/document"
	"github.com/Mintegral-official/juno/helpers"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

var doc1 = &document.DocInfo{
	Id: 0,
	Fields: []*document.Field{
		{
			Name:      "field1",
			IndexType: 1,
			Value:     nil,
		},
		{
			Name:      "field2",
			IndexType: 0,
			Value:     nil,
		},
		{
			Name:      "field1",
			IndexType: 2,
			Value:     nil,
		},
	},
}

var doc2 = &document.DocInfo{
	Id: 1,
	Fields: []*document.Field{
		{
			Name:      "field1",
			IndexType: 0,
			Value:     nil,
		},
		{
			Name:      "field2",
			IndexType: 1,
			Value:     nil,
		},
		{
			Name:      "field1",
			IndexType: 0,
			Value:     nil,
		},
	},
}

var doc3 = &document.DocInfo{
	Id: 2,
	Fields: []*document.Field{
		{
			Name:      "field1",
			IndexType: 0,
			Value:     nil,
		},
		{
			Name:      "field2",
			IndexType: 0,
			Value:     nil,
		},
		{
			Name:      "field1",
			IndexType: 1,
			Value:     nil,
		},
	},
}

func TestNewIndex(t *testing.T) {
	Convey("NewIndex", t, func() {
		So(NewIndex("index"), ShouldNotBeNil)
	})

	Convey("Add", t, func() {
		index := NewIndex("index")
		So(index.Add(nil), ShouldEqual, helpers.DocumentError)
		So(index.Add(doc1), ShouldBeNil)
		So(index.Add(doc2), ShouldBeNil)
		So(index.Add(doc3), ShouldBeNil)
		if1 := index.invertedIndex.Iterator("field1")
		c := 0
		for if1.HasNext() {
			if if1.Next() != nil {
				c++
			}
		}
		So(c, ShouldEqual, 3)

		if2 := index.invertedIndex.Iterator("field2")
		c = 0
		for if2.HasNext() {
			if if2.Next() != nil {
				c++
			}
		}
		So(c, ShouldEqual, 2)
		sf1 := index.storageIndex.Iterator("field1")
		c = 0
		for sf1.HasNext() {
			if sf1.Next() != nil {
				c++
			}
		}
		So(c, ShouldEqual, 2)
		sf2 := index.storageIndex.Iterator("field2")
		c = 0
		for sf2.HasNext() {
			if sf2.Next() != nil {
				c++
			}
		}
		So(c, ShouldEqual, 1)
		So(index.Del(doc1), ShouldBeNil)
	})
}

func TestInterface(t *testing.T) {
	var a, b interface{}
	i := int64(1)
	a = &i
	b = &i

	if b == a {
		fmt.Println("xxxxxxxxxxxxxxx")
	}

}
