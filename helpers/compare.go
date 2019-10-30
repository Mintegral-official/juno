package helpers

import (
	"github.com/Mintegral-official/juno/document"
	"strings"
)

/**
 * @author: tangye
 * @Date: 2019/10/30 14:19
 * @Description:
 */

type Comparable interface {
	Compare(a, b interface{}) int
}

type Func func(a, b interface{}) int

func (f Func) Compare(a, b interface{}) int {
	return f(a, b)
}

var IntCompare Func = func(a, b interface{}) int {
	if a.(int) < b.(int) {
		return -1
	} else if a.(int) > b.(int) {
		return 1
	}
	return 0
}

var Float32Compare Func = func(a, b interface{}) int {
	if a.(float32) < b.(float32) {
		return -1
	} else if a.(float32) > b.(float32) {
		return 1
	}
	return 0
}

var Float64Compare Func = func(a, b interface{}) int {
	if a.(float64) < b.(float64) {
		return -1
	} else if a.(float64) > b.(float64) {
		return 1
	}
	return 0
}

var StringCompare Func = func(a, b interface{}) int {
	return strings.Compare(a.(string), b.(string))
}

var DocIdFunc Func = func(a, b interface{}) int {
	if a.(document.DocId) > b.(document.DocId) {
		return 1
	} else if a.(document.DocId) < b.(document.DocId) {
		return -1
	}
	return 0
}


