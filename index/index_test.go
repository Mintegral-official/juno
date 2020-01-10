package index

import (
	"fmt"
	"github.com/Mintegral-official/juno/datastruct"
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
			Value:     1,
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
			IndexType: 2,
			Value:     "1",
			ValueType: document.StringFieldType,
		},
	},
}

var doc2 = &document.DocInfo{
	Id: 1,
	Fields: []*document.Field{
		{
			Name:      "field1",
			IndexType: 0,
			Value:     "1",
			ValueType: document.StringFieldType,
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
			Value:     "1",
			ValueType: document.StringFieldType,
		},
	},
}

var doc3 = &document.DocInfo{
	Id: 2,
	Fields: []*document.Field{
		{
			Name:      "field1",
			IndexType: 0,
			Value:     "1",
			ValueType: document.StringFieldType,
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
			Value:     1,
			ValueType: document.IntFieldType,
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
		if1 := index.GetInvertedIndex().Iterator("field1", "1")
		c := 0
		for if1.HasNext() {
			if if1.Current() != nil {
				c++
			}
			if1.Next()
		}
		So(c, ShouldEqual, 3)

		if2 := index.invertedIndex.Iterator("field2", "2")
		c = 0
		for if2.HasNext() {
			if if2.Current() != nil {
				c++
			}
			if2.Next()
		}
		So(c, ShouldEqual, 2)
		sf1 := index.GetStorageIndex().Iterator("field1")
		c = 0
		for sf1.HasNext() {
			if sf1.Current() != nil {
				c++
			}
			sf1.Next()
		}
		So(c, ShouldEqual, 2)
		sf2 := index.storageIndex.Iterator("field2")
		c = 0
		for sf2.HasNext() {
			if sf2.Current() != nil {
				c++
			}
			sf2.Next()
		}
		So(c, ShouldEqual, 1)
		So(len(*index.GetBitMap()), ShouldEqual, 32768)
		So(index.GetCampaignMap(), ShouldNotBeNil)
		So(index.GetDataType("field1"), ShouldEqual, 1)
		So(index.GetDataType("field2"), ShouldEqual, 3)
	})

	Convey("Del", t, func() {
		index := NewIndex("index")
		So(index.Add(nil), ShouldEqual, helpers.DocumentError)
		So(index.Add(doc1), ShouldBeNil)
		So(index.Add(doc2), ShouldBeNil)
		So(index.Add(doc3), ShouldBeNil)
		index.Del(doc1)
		if1 := index.GetInvertedIndex().Iterator("field1", "1")
		c := 0
		for if1.HasNext() {
			if if1.Current() != nil {
				c++
			}
			if1.Next()
		}
		So(c, ShouldEqual, 2)

		if2 := index.invertedIndex.Iterator("field2", "2")
		c = 0
		for if2.HasNext() {
			if if2.Current() != nil {
				c++
			}
			if2.Next()
		}
		So(c, ShouldEqual, 1)
		sf1 := index.GetStorageIndex().Iterator("field1")
		c = 0
		for sf1.HasNext() {
			if sf1.Current() != nil {
				c++
			}
			sf1.Next()
		}
		So(c, ShouldEqual, 1)
		sf2 := index.storageIndex.Iterator("field2")
		c = 0
		for sf2.HasNext() {
			if sf2.Current() != nil {
				c++
			}
			sf2.Next()
		}
		So(c, ShouldEqual, 1)
		So(len(*index.GetBitMap()), ShouldEqual, 32768)
		So(index.GetCampaignMap(), ShouldNotBeNil)
		So(index.GetDataType("field1"), ShouldEqual, 1)
		So(index.GetDataType("field2"), ShouldEqual, 3)
	})
}

func f1(a interface{}) interface{} {
	return a.(bool)
}

func TestInterface(t *testing.T) {
	var a interface{} = true
	fmt.Println(f1(a))
}

func TestStorageIndexer_Add(t *testing.T) {
	var a = &document.DocInfo{
		Id: 0,
		Fields: []*document.Field{
			{
				Name:      "f1",
				IndexType: 0,
				Value:     []int64{1, 2, 3},
				ValueType: document.SliceFieldType,
			},
			{
				Name:      "f2",
				IndexType: 1,
				Value:     []float64{1.1, 2.2, 3.4},
				ValueType: document.SliceFieldType,
			},
			{
				Name:      "f1",
				IndexType: 0,
				Value:     []int64{1, 22, 33},
				ValueType: document.SliceFieldType,
			},
		},
	}
	Convey("add", t, func() {
		index := NewIndex("index")
		_ = index.Add(a)
		idx := index.invertedIndex.Iterator("f1", "1")
		So(idx.HasNext(), ShouldBeTrue)
		c := 0
		for idx.HasNext() {
			if idx.Current() != nil {
				c++
			}
			idx.Next()
		}
		So(c, ShouldEqual, 1)
		sto := index.storageIndex.Iterator("f2")
		So(sto.HasNext(), ShouldBeTrue)
		c = 0
		for sto.HasNext() {
			if sto.Current().(*datastruct.Element).Value() != nil {
				c++
			}
			sto.Next()
		}
		So(c, ShouldEqual, 1)
	})
}
