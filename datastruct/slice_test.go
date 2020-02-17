package datastruct

import (
	"sort"
	"testing"
)

var slice []int
var m = make(map[int]interface{})

func binarySearch(sortedArray []int, lookingFor int) int {
	var low = 0
	var high = len(sortedArray) - 1
	for low <= high {
		var mid = low + (high-low)/2
		var midValue = sortedArray[mid]
		if midValue == lookingFor {
			return mid
		} else if midValue > lookingFor {
			high = mid - 1
		} else {
			low = mid + 1
		}
	}
	return -1
}

func add2() {
	for i := 0; i < 200000; i++ {
		slice = append(slice, arr[i])
	}
}

func add3() {
	for i := 0; i < 200000; i++ {
		m[arr[i]] = [1]byte{}
	}
}

func get2() {
	for i := 0; i < 100000; i++ {
		binarySearch(slice, arr[i])
	}
}

func get3() {
	for i := 0; i < 100000; i++ {
		_, _ = m[arr[i]]
	}
}

func BenchmarkSlice_Add(b *testing.B) {
	add2()
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		add2()
	}
}

func BenchmarkSlice_Get(b *testing.B) {
	add2()
	sort.Ints(slice)
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		get2()
	}
}

func BenchmarkSlice_Get_RunParallel(b *testing.B) {
	add2()
	sort.Ints(slice)
	b.ResetTimer()
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			get2()
		}
	})
}

func BenchmarkMap_Get(b *testing.B) {
	add3()
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		get3()
	}
}

func BenchmarkMap_Get_RunParallel(b *testing.B) {
	add3()
	b.ResetTimer()
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			get3()
		}
	})
}
