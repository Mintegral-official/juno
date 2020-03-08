package index

import (
	"fmt"
	"github.com/Mintegral-official/juno/document"
	"github.com/Mintegral-official/juno/helpers"
	. "github.com/smartystreets/goconvey/convey"
	"reflect"
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
			Value:     100,
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

		So(index.GetValueById(2), ShouldNotBeNil)
		if1 := index.GetInvertedIndex().Iterator("field1", "1")
		So(index.GetValueById(1), ShouldNotBeNil)
		index.UpdateIds("test\007tt", []document.DocId{1, 2, 100})
		So(index.GetValueById(1), ShouldNotBeNil)
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
		//So(len(*index.GetBitMap()), ShouldEqual, 32768)
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
		//So(len(*index.GetBitMap()), ShouldEqual, 32768)
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
			if sto.Current() != nil {
				c++
			}
			sto.Next()
		}
		So(c, ShouldEqual, 1)
	})
}

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
		{
			Name:      "field3",
			IndexType: 0,
			Value:     "3",
			ValueType: document.StringFieldType,
		},
		{
			Name:      "field4",
			IndexType: 0,
			Value:     "34",
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
		{
			Name:      "field3",
			IndexType: 0,
			Value:     "30",
			ValueType: document.StringFieldType,
		},
		{
			Name:      "field3",
			IndexType: 0,
			Value:     "300",
			ValueType: document.StringFieldType,
		},
	},
}

func TestNewStorageIndexer(t *testing.T) {
	idx := NewIndex("")
	_ = idx.Add(doc4)
	Convey("GetValueById add", t, func() {
		realMap := idx.GetValueById(0)
		expectMap := [2]map[string][]string{
			{
				"field2": []string{"2"},
				"field3": []string{"3"},
				"field4": []string{"34"},
			},
			{
				"field1": []string{"1"},
			},
		}
		So(reflect.DeepEqual(realMap, expectMap), ShouldBeTrue)
	})
	idx.Del(doc5)
	_ = idx.Add(doc5)
	Convey("GetValueById del & add", t, func() {
		realMap := idx.GetValueById(0)
		expectMap := [2]map[string][]string{
			{
				"field2": []string{"20", "200"},
				"field3": []string{"30", "300"},
			},
			{
				"field1": []string{"10"},
			},
		}
		So(reflect.DeepEqual(realMap, expectMap), ShouldBeTrue)
	})
}
