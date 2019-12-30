package helpers

import (
	"fmt"
	"github.com/Mintegral-official/juno/document"
)

func Compare(i, j interface{}) int {
	switch i.(type) {
	case int8, int16, int32, int64, int, *int8, *int16, *int32, *int64, *int:
		return intCompare(i, j)
	case string, *string:
		return stringCompare(i, j)
	case float32, *float32, float64, *float64:
		return floatCompare(i, j)
	default:
		panic(fmt.Sprintf("parameters[%T - %T] type wrong.", i, j))
	}
}

func In(target int, arr []int) bool {
	for _, v := range arr {
		if target == v {
			return true
		}
	}
	return false
}

func Merge(arr1, arr2 []document.DocId) []document.DocId {
	if len(arr1) == 0 {
		return arr2
	}
	if len(arr2) == 0 {
		return arr1
	}
	i, j := 0, 0
	var res []document.DocId
	for i < len(arr1) && j < len(arr2) {
		if arr1[i] < arr2[j] {
			res = append(res, arr1[i])
			i++
		} else if arr1[i] == arr2[j] {
			res = append(res, arr1[i])
			i++
			j++
		} else {
			res = append(res, arr2[j])
			j++
		}
	}
	for ; i < len(arr1); i++ {
		res = append(res, arr1[i])
	}
	for ; j < len(arr2); j++ {
		res = append(res, arr1[j])
	}
	return res
}
