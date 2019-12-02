package juno

import (
	"fmt"
	"testing"
)

func TestIn(t *testing.T) {
	f(100, 99)
	f("abc", "abc")
}

func f(a, b interface{}) {
	fmt.Println(a == b)
}