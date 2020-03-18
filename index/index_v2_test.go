package index

import (
	"fmt"
	"github.com/MintegralTech/juno/document"
	"github.com/MintegralTech/juno/helpers"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestNewIndexV2(t *testing.T) {
	Convey("NewIndex", t, func() {
		So(NewIndexV2("index"), ShouldNotBeNil)
	})

	Convey("Add", t, func() {
		index := NewIndexV2("index")
		So(index.Add(nil), ShouldEqual, helpers.DocumentError)
		So(index.Add(doc1), ShouldBeNil)
		So(index.Add(doc2), ShouldBeNil)
		So(index.Add(doc3), ShouldBeNil)

		So(index.GetIndexDebugInfoById(2), ShouldNotBeNil)
		if1 := index.GetInvertedIndex().Iterator("field1", "1")
		So(index.GetIndexDebugInfoById(1), ShouldNotBeNil)
		c := 0
		for if1.HasNext() {
			if if1.Current() != nil {
				c++
			}
			if1.Next()
		}
		So(c, ShouldEqual, 3)

		if2 := index.GetInvertedIndex().Iterator("field2", "2")
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
		sf2 := index.GetStorageIndex().Iterator("field2")
		c = 0
		for sf2.HasNext() {
			if sf2.Current() != nil {
				c++
			}
			sf2.Next()
		}
		So(c, ShouldEqual, 1)
		//So(len(*index.GetBitMap()), ShouldEqual, 32768)
		So(index.GetDataType("field1"), ShouldEqual, 1)
		So(index.GetDataType("field2"), ShouldEqual, 3)
	})
}

func TestStorageIndexerV2_Add(t *testing.T) {
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
		index := NewIndexV2("index")
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

func TestMergeIndexV2(t *testing.T) {
	Convey("TestMergeIndexV2", t, func() {
		idx := NewIndexV2("")
		So(idx.Add(doc1), ShouldBeNil)
		So(idx.Add(doc2), ShouldBeNil)
		So(idx.GetInvertedIndex().Count(), ShouldEqual, 2)
		So(idx.GetStorageIndex().Count(), ShouldEqual, 2)
		Convey("GetValueById add", func() {
			realMap := idx.GetIndexDebugInfoById(0)
			//expectMap := [2]map[string][]string{
			//	{
			//		"field1": []string{"1"},
			//		"field2": []string{"2"},
			//	},
			//	{
			//		"field1": []string{"1"},
			//	},
			//}
			fmt.Println(realMap)
		})
		idx2 := NewIndexV2("1234")
		Convey("merge nil index", func() {
			idx2.MergeIndex(idx)
			So(idx.GetInvertedIndex().Count(), ShouldEqual, 2)
			So(idx.GetStorageIndex().Count(), ShouldEqual, 2)
		})

		So(idx2.Add(doc5), ShouldBeNil)
		So(idx2.Add(doc3), ShouldBeNil)

		Convey("GetValueById del & add", func() {
			realMap := idx2.GetIndexDebugInfoById(0)
			//expectMap := &IndexDebugInfo{
			//	map[string][]string{
			//		"field2": []string{"20", "200"},
			//		"field3": []string{"30", "300"},
			//	},
			//	map[string][]string{
			//		"field1": []string{"10"},
			//	},
			//}
			fmt.Println(realMap)
		})

		idx2.MergeIndex(idx)
		Convey("GetValueById merge", func() {
			realMap := idx2.GetIndexDebugInfoById(0)
			//expectMap := &IndexDebugInfo{
			//	map[string][]string{
			//		"field2": []string{"20", "200"},
			//		"field3": []string{"30", "300"},
			//	},
			//	map[string][]string{
			//		"field1": []string{"10"},
			//	},
			//}
			fmt.Println(realMap)
			id, ok := idx2.campaignMapping.Get(DocId(2))
			So(id, ShouldEqual, 1)
			So(ok, ShouldBeTrue)

			id, ok = idx2.campaignMapping.Get(DocId(0))
			So(id, ShouldEqual, 0)
			So(ok, ShouldBeTrue)

			id, ok = idx2.campaignMapping.Get(DocId(1))
			So(id, ShouldEqual, 2)
			So(ok, ShouldBeTrue)
		})
	})
}
