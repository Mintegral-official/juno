package juno

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
			IndexType: 0,
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
		So(index.Add(nil), ShouldEqual, helpers.DOCUMENT_ERROR)
		So(index.Add(doc1), ShouldBeNil)
		So(index.Add(doc2), ShouldBeNil)
		So(index.Add(doc3), ShouldBeNil)
		if1 := index.invertedIndex.Iterator("field1")
		for if1.HasNext() {
			fmt.Println(if1.Next())
		}
		fmt.Println("************************************")
		if2 := index.invertedIndex.Iterator("field2")
		for if2.HasNext() {
			fmt.Println(if2.Next())
		}

		fmt.Println("************************************")
		sf1 := index.storageIndex.Iterator("field2")
		for sf1.HasNext() {
			fmt.Println(sf1.Next())
		}

		fmt.Println("************************************")
		sf2 := index.storageIndex.Iterator("field2")
		for sf2.HasNext() {
			fmt.Println(sf2.Next())
		}

		So(index.Del(doc1), ShouldBeNil)

	})
}
