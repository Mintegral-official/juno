package index

import (
	"fmt"
	"testing"
	"unsafe"
)

type A struct {
	int
}

func TestSimpleInvertedIndex_Add(t *testing.T) {
	fmt.Println(unsafe.Sizeof(A{}))
}
