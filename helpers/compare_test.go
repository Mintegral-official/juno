package helpers

import (
	"fmt"
	"testing"
)

func TestFunc_Compare(t *testing.T) {
	var a, b *rune
	var c, d rune = 2, 1
	a = &c
	b = &d
	fmt.Println(intCompare(*a, *b))
	fmt.Println(intCompare(c, d))
	fmt.Println(Compare(a, b))
}
