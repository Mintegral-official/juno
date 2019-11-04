package main

/**
 * @author: tangye
 * @Date: 2019/11/4 19:31
 * @Description:
 */

/*
#include <stdio.h>
int t() {
    return rand() % (1000000000 - 0 + 1) + 0;
}
 */
import "C"
import (
	"fmt"
	"github.com/Mintegral-official/juno/helpers"
	"github.com/Mintegral-official/juno/index"
	"sort"
	"time"
	"unsafe"
)
var s *index.SkipListIterator = index.NewSKipListIterator(index.DEFAULT_MAX_LEVEL, helpers.IntCompare)
var arr [500001]int
var s1 = make([]int, 500001)

func init() {
	b := time.Now()
	for i := 0; i < 15000000; i++ {
		_ = int(C.t())
	}
	fmt.Println(time.Since(b))

    a := time.Now()
	for i := 0; i < 15000000; i++ {
		//s.Add(int(C.t()), [1]byte{})
		s1 = append(s1, int(C.t()))
	}
	fmt.Println(time.Since(a))
	sort.Ints(s1)

	for i := 0; i < 500000; i++ {
		arr[i] = int(C.t())
	}

	var sl index.SkipList
	var el index.Element
	fmt.Printf("Structure sizes: SkipList is %v, Element is %v bytes\n", unsafe.Sizeof(sl), unsafe.Sizeof(el))
}

func main() {
	t := time.Now()

	for j := 0; j < 500001; j++ {
		//s.Get(arr[j])
		_ = s1[arr[j] % 15000000]
	}

	fmt.Println(time.Since(t))
}