package datastruct

import (
	"fmt"
	"github.com/MintegralTech/juno/document"
	. "github.com/smartystreets/goconvey/convey"
	"math/rand"
	"strconv"
	"testing"
)

func TestNewBitMap(t *testing.T) {
	Convey("NewBitMap", t, func() {
		bm := NewBitMap()
		So(bm, ShouldNotBeNil)
		So(len(*bm), ShouldEqual, 32768)
		So(cap(*bm), ShouldEqual, 32768)
		So(bm.IsExist(10), ShouldBeFalse)
		bm.Set(10)
		So(bm.IsExist(10), ShouldBeTrue)
		So(bm.IsFully(), ShouldBeTrue)
		bm.Del(10)
		So(bm.IsExist(10), ShouldBeFalse)
	})
}

func TestNewBitMap3(t *testing.T) {
	Convey("bitmap exist", t, func() {
		bm := NewBitMap()
		bm.Set(10)
		So(bm.IsExist(10), ShouldBeTrue)
		bm.Del(10)
		So(bm.IsExist(10), ShouldBeFalse)
	})
}

func BenchmarkBitMap(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	bm := NewBitMap()
	for i := 0; i < b.N; i++ {
		for j := 0; j < 100000000/8; j++ {
			bm.Set(document.DocId(j + 1))
		}
	}
}

func BenchmarkBitSet_Del(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = fmt.Sprint(rand.Int())
	}
}

func BenchmarkBitSet_IsExist(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = strconv.Itoa(rand.Int())
	}
}
