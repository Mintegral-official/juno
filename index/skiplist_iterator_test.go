package index

import (
	"github.com/Mintegral-official/juno/helpers"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

var s1 = NewSkipListIterator(DEFAULT_MAX_LEVEL, helpers.IntCompare)

//func init() {
//	for i := 0; i < 100; i++ {
//		s1.Add(i, nil)
//	}
//}

func TestNewSkipListIterator(t *testing.T) {
		for i := 0; i < 100; i++ {
			s1.Add(i, nil)
		}
	// slt := NewSkipListIterator(DEFAULT_MAX_LEVEL, helpers.IntCompare)
	Convey("NewSkipListIterator", t, func() {
		So(NewSkipListIterator(DEFAULT_MAX_LEVEL, helpers.IntCompare), ShouldNotBeNil)
	})
}

func TestSkipListIterator_HasNext(t *testing.T) {
		for i := 0; i < 100; i++ {
			s1.Add(i, nil)
		}
	Convey("HasNext", t, func() {
		So(s1.HasNext(), ShouldBeTrue)
		So(s1.Iterator(), ShouldNotBeNil)
	})
}

func TestSkipListIterator_Iterator(t *testing.T) {
		for i := 0; i < 100; i++ {
			s1.Add(i, nil)
		}
	Convey("Next", t, func() {
		So(s1.GetGE(5), ShouldNotBeNil)
		So(s1.GetGE(10), ShouldNotBeNil)
		c := 0
		for s1.HasNext() {
			s1.Next()
			c++
			if c == 10 {
				break
			}
		}
		// So(s1.index, ShouldEqual, 10)
		So(s1.GetGE(5), ShouldBeNil)
		So(s1.GetGE(10), ShouldBeNil)
		So(s1.GetGE(11), ShouldNotBeNil)

	})
}

func add1() {
	for i := 0; i < 200000; i++ {
		s1.Add(i, [1]byte{})
	}
}

func getGE() {
	for i := 0; i < 100000; i++ {
		s1.GetGE(arr[i])
	}
}

func BenchmarkSkipListIterator_GetGE(b *testing.B) {
	add1()
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		getGE()
	}
}

func BenchmarkSkipListIterator_GetGE_RunParallel(b *testing.B) {
	add1()
	b.ResetTimer()
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			getGE()
		}
	})
}

func BenchmarkNewSkipListIterator_Next(b *testing.B) {
	add1()
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		for s1.HasNext() {
			s1.Next()
		}
	}
}

func BenchmarkNewSkipListIterator_Next_RunParallel(b *testing.B) {
	add1()
	b.ResetTimer()
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			for s1.HasNext() {
				s1.Next()
			}
		}
	})
}