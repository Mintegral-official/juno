package helpers

import "github.com/Mintegral-official/juno/document"

func Compare(i, j interface{}) int {
	switch i.(type) {
	case document.DocId:
		return DocIdFunc(i, j)
	case int:
		return IntCompare(i, j)
	case string:
		return StringCompare(i, j)
	case float32:
		return Float32Compare(i, j)
	case float64:
		return Float64Compare(i, j)
	default:
		return IntCompare(i, j)
	}
}

func In(target interface{}, arrays ...interface{}) bool {
    for _, v := range arrays {
		if Compare(target, v) == 0 {
			return true
		}
	}
	return false
}

func Equal(i, j interface{}) bool {
	return Compare(i, j) == 0
}

