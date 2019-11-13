package index

import (
	"fmt"
	"github.com/Mintegral-official/juno/helpers"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestNewSkipListIterator(t *testing.T) {
	a := NewSkipListIterator(DEFAULT_MAX_LEVEL, helpers.IntCompare)
	for i := 0; i < 100; i++ {
		a.Add(i, nil)
	}
	// slt := NewSkipListIterator(DEFAULT_MAX_LEVEL, helpers.IntCompare)
	Convey("NewSkipListIterator", t, func() {
		So(NewSkipListIterator(DEFAULT_MAX_LEVEL, helpers.IntCompare), ShouldNotBeNil)
	})
}

func TestSkipListIterator_HasNext(t *testing.T) {
	a := NewSkipListIterator(DEFAULT_MAX_LEVEL, helpers.IntCompare)
	for i := 0; i < 100; i++ {
		a.Add(i, nil)
	}
	Convey("HasNext", t, func() {
		So(a.HasNext(), ShouldBeTrue)
		So(a.Iterator(), ShouldNotBeNil)
	})
}

func TestSkipListIterator_Iterator(t *testing.T) {
	a := NewSkipListIterator(DEFAULT_MAX_LEVEL, helpers.IntCompare)
	for i := 0; i < 100; i++ {
		a.Add(i, nil)
	}
	// fmt.Println(a.HasNext())
	//m := a.Iterator()
	//fmt.Println(m.Next())
	//fmt.Println(m.Next())
	//fmt.Println(m.Next())
	//fmt.Println(m.Next())
	Convey("Next", t, func() {
		So(a.GetGE(5), ShouldNotBeNil)
		So(a.GetGE(10), ShouldNotBeNil)
		c := 0
	//	a = a.Iterator()
		for a.HasNext() {
			a.Next()
			c++
			if c == 10 {
				break
			}
		}
		So(a.Index, ShouldEqual, 10)
		So(a.GetGE(5), ShouldBeNil)
		fmt.Println(a.GetGE(10))
		So(a.GetGE(10), ShouldNotBeNil)
		So(a.GetGE(11), ShouldNotBeNil)

	})
}

func getGE(s *SkipListIterator) {
	for i := 0; i < 100000; i++ {
		s.GetGE(arr[i])
	}
}

func add1(s *SkipListIterator) {
	for i := 0; i < 200000; i++ {
		s.Add(arr[i], [1]byte{})
	}
}

func TestSkipListIterator_First(t *testing.T) {
	a := NewSkipListIterator(DEFAULT_MAX_LEVEL, helpers.IntCompare)
	for i := 0; i < 1000; i++ {
		a.Add(i, nil)
	}
	fmt.Println(a.GetGE(10))
	fmt.Println(a.GetGE(324))
	a.Del(10)
	a.Del(324)
	fmt.Println(a.GetGE(10))
	fmt.Println(a.GetGE(324))
}

func BenchmarkSkipListIterator_GetGE(b *testing.B) {
	a := NewSkipListIterator(DEFAULT_MAX_LEVEL, helpers.IntCompare)
	add1(a)
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		getGE(a)
	}
}

func BenchmarkSkipListIterator_GetGE_RunParallel(b *testing.B) {
	a := NewSkipListIterator(DEFAULT_MAX_LEVEL, helpers.IntCompare)
	add1(a)
	b.ResetTimer()
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			getGE(a)
		}
	})
}

func BenchmarkNewSkipListIterator_Next(b *testing.B) {
	a := NewSkipListIterator(DEFAULT_MAX_LEVEL, helpers.IntCompare)
	add1(a)
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		for a.HasNext() {
			a.Next()
		}
	}
}

func BenchmarkNewSkipListIterator_Next_RunParallel(b *testing.B) {
	a := NewSkipListIterator(DEFAULT_MAX_LEVEL, helpers.IntCompare)
	add1(a)
	b.ResetTimer()
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			for a.HasNext() {
				a.Next()
			}
		}
	})
}
