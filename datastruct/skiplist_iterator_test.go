package datastruct

import (
	"fmt"
	"github.com/Mintegral-official/juno/helpers"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestNewSkipListIterator(t *testing.T) {
	sl, _ := NewSkipList(DefaultMaxLevel, helpers.IntCompare)
	sl.Add(1, nil)
	sl.Add(3, nil)
	// slt := NewSkipListIterator(DEFAULT_MAX_LEVEL, helpers.IntCompare)
	Convey("NewSkipListIterator", t, func() {
		iter := sl.Iterator()

		v := iter.Current().(*Element)
		So(v, ShouldNotBeNil)
		So(v.Key(), ShouldEqual, 1)

		So(iter.HasNext(), ShouldBeTrue)
		iter.Next()
		v = iter.Current().(*Element)
		So(v, ShouldNotBeNil)
		So(v.Key(), ShouldEqual, 3)

		So(iter.HasNext(), ShouldBeTrue)
		iter.Next()
		So(iter.Current(), ShouldBeNil)
	})
}

func TestSkipListIterator_Iterator(t *testing.T) {
	s, _ := NewSkipList(DefaultMaxLevel, helpers.IntCompare)
	for i := 0; i < 100; i++ {
		s.Add(i, nil)
	}
	for i := 101; i < 150; i += 3 {
		s.Add(i, nil)
	}

	Convey("Next", t, func() {
		iter := s.Iterator()
		So(iter.HasNext(), ShouldBeTrue)
		v := iter.Current().(*Element)
		So(v, ShouldNotBeNil)
		So(v.Key(), ShouldEqual, 0)

		iter.Next()
		elem := iter.Current()
		So(elem, ShouldNotBeNil)
		So(elem.(*Element), ShouldNotBeNil)
		So(elem.(*Element).Key(), ShouldEqual, 1)

		v = iter.GetGE(5).(*Element)
		So(v, ShouldNotBeNil)
		So(v.Key(), ShouldEqual, 5)
		elem = iter.Current()
		So(elem, ShouldNotBeNil)
		So(elem.(*Element), ShouldNotBeNil)
		So(elem.(*Element).Key(), ShouldEqual, 5)

		iter.Next()
		elem = iter.Current()
		So(elem, ShouldNotBeNil)
		So(elem.(*Element), ShouldNotBeNil)
		So(elem.(*Element).Key(), ShouldEqual, 6)

		v = iter.GetGE(102).(*Element)
		So(v, ShouldNotBeNil)
		So(v.Key(), ShouldEqual, 104)
		elem = iter.Current()
		So(elem, ShouldNotBeNil)
		So(elem.(*Element), ShouldNotBeNil)
		So(elem.(*Element).Key(), ShouldEqual, 104)
		So(iter.HasNext(), ShouldBeTrue)

		v = iter.GetGE(147).(*Element)
		So(v, ShouldNotBeNil)
		So(v.Key(), ShouldEqual, 149)
		elem = iter.Current()
		So(elem, ShouldNotBeNil)
		So(elem.(*Element), ShouldNotBeNil)
		So(elem.(*Element).Key(), ShouldEqual, 149)
		So(iter.HasNext(), ShouldBeTrue)

		v = iter.GetGE(160).(*Element)
		So(v, ShouldBeNil)
		elem = iter.Current()
		So(elem, ShouldBeNil)
		So(iter.HasNext(), ShouldBeFalse)
	})
}

func TestSkipListIterator_GetGE(t *testing.T) {
	s, _ := NewSkipList(DefaultMaxLevel, helpers.IntCompare)
	for i := 0; i < 100; i++ {
		s.Add(i, nil)
	}
	a := s.Iterator()
	fmt.Println(a.GetGE(99))
	fmt.Println(a.GetGE(99))
	fmt.Println(a.GetGE(99))
}

//
//func getGE(s *SkipListIterator) {
//	for i := 0; i < 100000; i++ {
//		s.GetGE(arr[i])
//	}
//}
//
//func add1(s *SkipListIterator) {
//	for i := 0; i < 200000; i++ {
//		s.Add(arr[i], [1]byte{})
//	}
//}
//
//func TestSkipListIterator_First(t *testing.T) {
//	a := NewSkipListIterator(DEFAULT_MAX_LEVEL, helpers.IntCompare)
//	for i := 0; i < 1000; i++ {
//		a.Add(i, nil)
//	}
//	fmt.Println(a.GetGE(10))
//	fmt.Println(a.GetGE(324))
//	a.Del(10)
//	a.Del(324)
//	fmt.Println(a.GetGE(10))
//	fmt.Println(a.GetGE(324))
//}
//
//func BenchmarkSkipListIterator_GetGE(b *testing.B) {
//	a := NewSkipListIterator(DEFAULT_MAX_LEVEL, helpers.IntCompare)
//	add1(a)
//	b.ResetTimer()
//	b.ReportAllocs()
//	for i := 0; i < b.N; i++ {
//		getGE(a)
//	}
//}
//
//func BenchmarkSkipListIterator_GetGE_RunParallel(b *testing.B) {
//	a := NewSkipListIterator(DEFAULT_MAX_LEVEL, helpers.IntCompare)
//	add1(a)
//	b.ResetTimer()
//	b.ReportAllocs()
//	b.RunParallel(func(pb *testing.PB) {
//		for pb.Next() {
//			getGE(a)
//		}
//	})
//}
//
//func BenchmarkNewSkipListIterator_Next(b *testing.B) {
//	a := NewSkipListIterator(DEFAULT_MAX_LEVEL, helpers.IntCompare)
//	add1(a)
//	b.ResetTimer()
//	b.ReportAllocs()
//	for i := 0; i < b.N; i++ {
//		for a.HasNext() {
//			a.Next()
//		}
//	}
//}
//
//func BenchmarkNewSkipListIterator_Next_RunParallel(b *testing.B) {
//	a := NewSkipListIterator(DEFAULT_MAX_LEVEL, helpers.IntCompare)
//	add1(a)
//	b.ResetTimer()
//	b.ReportAllocs()
//	b.RunParallel(func(pb *testing.PB) {
//		for pb.Next() {
//			for a.HasNext() {
//				a.Next()
//			}
//		}
//	})
//}
