package datastruct

import (
	"fmt"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestNewBitMap(t *testing.T) {
	Convey("NewBitMap", t, func() {
		var maxNum = 100
		bm := NewBitMap(maxNum)
		So(bm.bits, ShouldNotBeNil)
		So(bm.bitCount, ShouldEqual, 0)
		So(bm.capacity, ShouldEqual, maxNum)
		So(bm.IsEmpty(), ShouldBeTrue)
		So(bm.IsExist(10), ShouldBeFalse)
		So(bm.Count(), ShouldEqual, 0)
		bm.Set(10)
		So(bm.IsExist(10), ShouldBeTrue)
		So(bm.IsFully(), ShouldBeFalse)
		So(bm.IsEmpty(), ShouldBeFalse)
		So(bm.Count(), ShouldEqual, 1)
		So(bm.Get(), ShouldNotBeNil)
		So(len(bm.Get()), ShouldEqual, 1)
		bm.Del(10)
		So(bm.IsEmpty(), ShouldBeTrue)
		So(bm.IsExist(10), ShouldBeFalse)
		So(bm.Count(), ShouldEqual, 0)
	})
}

func TestNewBitMap3(t *testing.T) {
	Convey("NewBitMap", t, func() {
		var maxNum = 100
		bm := NewBitMap(maxNum)
		bm.Set(10)
		fmt.Println(bm.IsExist(10))
		bm.Del(10)
		fmt.Println(bm.IsExist(10))
	})
}

func BenchmarkBitMap(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	bm := NewBitMap(100000000)
	for i := 0; i < b.N; i++ {
		for j := 0; j < 100000000/8; j++ {
			bm.Set(j + 1)
		}
	}
	for i := 0; i < b.N; i++ {
		bm.Get()
	}
}

func TestBitMap_Count(t *testing.T) {
	var a []int
	a = append(a, 1)
	fmt.Println(len(a))
}
