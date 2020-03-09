package datastruct

import (
	"github.com/MintegralTech/juno/document"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestNewSkipListIterator(t *testing.T) {
	sl := NewSkipList(DefaultMaxLevel)
	sl.Add(1, nil)
	sl.Add(3, nil)
	Convey("NewSkipListIterator", t, func() {
		iter := sl.Iterator()

		v := iter.Current()
		So(v, ShouldNotBeNil)
		So(v.Key(), ShouldEqual, 1)

		So(iter.HasNext(), ShouldBeTrue)
		iter.Next()
		v = iter.Current()
		So(v, ShouldNotBeNil)
		So(v.Key(), ShouldEqual, 3)

		So(iter.HasNext(), ShouldBeTrue)
		iter.Next()
		So(iter.Current(), ShouldBeNil)
	})
}

func TestSkipListIterator_Iterator(t *testing.T) {
	s := NewSkipList(DefaultMaxLevel)
	for i := 0; i < 100; i++ {
		s.Add(document.DocId(i), nil)
	}
	for i := 101; i < 150; i += 3 {
		s.Add(document.DocId(i), nil)
	}

	Convey("Next", t, func() {
		iter := s.Iterator()
		So(iter.HasNext(), ShouldBeTrue)
		v := iter.Current()
		So(v, ShouldNotBeNil)
		So(v.Key(), ShouldEqual, 0)

		iter.Next()
		elem := iter.Current()
		So(elem, ShouldNotBeNil)
		So(elem, ShouldNotBeNil)
		So(elem.Key(), ShouldEqual, 1)

		v = iter.GetGE(5)
		So(v, ShouldNotBeNil)
		So(v.Key(), ShouldEqual, 5)
		elem = iter.Current()
		So(elem, ShouldNotBeNil)
		So(elem, ShouldNotBeNil)
		So(elem.Key(), ShouldEqual, 5)

		iter.Next()
		elem = iter.Current()
		So(elem, ShouldNotBeNil)
		So(elem, ShouldNotBeNil)
		So(elem.Key(), ShouldEqual, 6)

		v = iter.GetGE(102)
		So(v, ShouldNotBeNil)
		So(v.Key(), ShouldEqual, 104)
		elem = iter.Current()
		So(elem, ShouldNotBeNil)
		So(elem, ShouldNotBeNil)
		So(elem.Key(), ShouldEqual, 104)
		So(iter.HasNext(), ShouldBeTrue)

		v = iter.GetGE(147)
		So(v, ShouldNotBeNil)
		So(v.Key(), ShouldEqual, 149)
		elem = iter.Current()
		So(elem, ShouldNotBeNil)
		So(elem, ShouldNotBeNil)
		So(elem.Key(), ShouldEqual, 149)
		So(iter.HasNext(), ShouldBeTrue)

		v = iter.GetGE(160)
		So(v, ShouldBeNil)
		elem = iter.Current()
		So(elem, ShouldBeNil)
		So(iter.HasNext(), ShouldBeFalse)
	})
}

func TestSkipListIterator_GetGE(t *testing.T) {
	s := NewSkipList(DefaultMaxLevel)
	for i := 0; i < 100; i++ {
		s.Add(document.DocId(i), nil)
	}
	a := s.Iterator()

	Convey("getGE", t, func() {
		v := a.GetGE(99)
		So(v.key, ShouldEqual, 99)

		v = a.GetGE(99)
		So(v.key, ShouldEqual, 99)

		v = a.GetGE(99)
		So(v.key, ShouldEqual, 99)
	})
}

func getGE(s *SkipListIterator) {
	for i := 0; i < 100000; i++ {
		s.GetGE(document.DocId(arr[i]))
	}
}

func add1(s *SkipList) {
	for i := 0; i < 200000; i++ {
		s.Add(document.DocId(arr[i]), [1]byte{})
	}
}

func TestSkipListIterator_First(t *testing.T) {
	a := NewSkipList(DefaultMaxLevel)
	for i := 0; i < 1000; i++ {
		a.Add(document.DocId(i), nil)
	}
	s := a.Iterator()
	Convey("del", t, func() {
		v := s.GetGE(10)
		So(v.key, ShouldEqual, 10)
		v = s.GetGE(324)
		So(v.key, ShouldEqual, 324)
		a.Del(10)
		a.Del(324)
		So(a.Len(), ShouldEqual, 998)
		v = s.GetGE(10)
		So(v.key, ShouldEqual, 324)
		v = s.GetGE(324)
		So(v.key, ShouldEqual, 324)
	})

}

func BenchmarkSkipListIterator_GetGE(b *testing.B) {
	a := NewSkipList(DefaultMaxLevel)
	add1(a)
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		getGE(a.Iterator())
	}
}

func BenchmarkSkipListIterator_GetGE_RunParallel(b *testing.B) {
	a := NewSkipList(DefaultMaxLevel)
	add1(a)
	b.ResetTimer()
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			getGE(a.Iterator())
		}
	})
}

func BenchmarkNewSkipListIterator_Next(b *testing.B) {
	a := NewSkipList(DefaultMaxLevel)
	add1(a)
	s := a.Iterator()
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		for s.HasNext() {
			s.Next()
		}
	}
}
