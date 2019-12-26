package helpers

import (
	"fmt"
	"strings"
)

const ACCURACY = 0.000001

type Comparable interface {
	Compare(a, b interface{}) int
}

type Func func(a, b interface{}) int

func (f Func) Compare(a, b interface{}) int {
	return f(a, b)
}

var intCompare Func = func(a, b interface{}) int {
	switch a.(type) {
	case int8:
		return int8Func(a.(int8), b.(int8))
	case *int8:
		return int8PtrFunc(a.(*int8), b.(*int8))
	case int16:
		return int16Func(a.(int16), b.(int16))
	case *int16:
		return int16PtrFunc(a.(*int16), b.(*int16))
	case int:
		return intFunc(a.(int), b.(int))
	case *int:
		return intPtrFunc(a.(*int), b.(*int))
	case int32:
		return int32Func(a.(int32), b.(int32))
	case *int32:
		return int32PtrFunc(a.(*int32), b.(*int32))
	case int64:
		return int64Func(a.(int64), b.(int64))
	case *int64:
		return int64PtrFunc(a.(*int64), b.(*int64))
	case byte:
		return byteFunc(a.(byte), b.(byte))
	case *byte:
		return bytePtrFunc(a.(*byte), b.(*byte))
	default:
		panic(fmt.Sprintf("parameters[%T - %T] type wrong.", a, b))
	}
}

var floatCompare Func = func(a, b interface{}) int {
	switch a.(type) {
	case float32:
		return float32Func(a.(float32), b.(float32))
	case *float32:
		return float32PtrFunc(a.(*float32), b.(*float32))
	case float64:
		return float64Func(a.(float64), b.(float64))
	case *float64:
		return float64PtrFunc(a.(*float64), b.(*float64))
	default:
		panic(fmt.Sprintf("parameters[%T - %T] type wrong.", a, b))
	}
}

var stringCompare Func = func(a, b interface{}) int {
	switch a.(type) {
	case string:
		return stringFunc(a.(string), b.(string))
	case *string:
		return stringPtrFunc(a.(*string), b.(*string))
	default:
		panic(fmt.Sprintf("parameters[%T - %T] type wrong.", a, b))
	}
}

func int8Func(i, j int8) int {
	return int(i - j)
}

func int16Func(i, j int16) int {
	return int(i - j)
}

func int32Func(i, j int32) int {
	return int(i - j)
}

func int64Func(i, j int64) int {
	return int(i - j)
}

func intFunc(i, j int) int {
	return i - j
}

func byteFunc(i, j byte) int {
	return int(i - j)
}

func float32Func(i, j float32) int {
	if j-i > ACCURACY {
		return -1
	} else if i-j > ACCURACY {
		return 1
	}
	return 0
}

func float64Func(i, j float64) int {
	if j-i > ACCURACY {
		return -1
	} else if i-j > ACCURACY {
		return 1
	}
	return 0
}

func int8PtrFunc(i, j *int8) int {
	return int8Func(*i, *j)
}

func int16PtrFunc(i, j *int16) int {
	return int16Func(*i, *j)
}

func int32PtrFunc(i, j *int32) int {
	return int32Func(*i, *j)
}

func int64PtrFunc(i, j *int64) int {
	return int64Func(*i, *j)
}

func intPtrFunc(i, j *int) int {
	return intFunc(*i, *j)
}

func bytePtrFunc(i, j *byte) int {
	return byteFunc(*i, *j)
}

func float32PtrFunc(i, j *float32) int {
	return float32Func(*i, *j)
}

func float64PtrFunc(i, j *float64) int {
	return float64Func(*i, *j)
}

func stringFunc(i, j string) int {
	return strings.Compare(i, j)
}

func stringPtrFunc(i, j *string) int {
	return strings.Compare(*i, *j)
}
