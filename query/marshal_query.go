package query

import (
	"encoding/json"
	"unsafe"
)

type JSONFormatter struct {
}

func NewJSONFormatter() *JSONFormatter {
	return &JSONFormatter{}
}

func (j *JSONFormatter) Marshal(cond interface{}) (string, error) {
	if r, e := json.Marshal(cond); e != nil {
		return "", e
	} else {
		return *(*string)(unsafe.Pointer(&r)), nil
	}

}

func (j *JSONFormatter) Unmarshal(str string, cond interface{}, queryFunc func(cond interface{}) Query) (Query, error) {
	if e := json.Unmarshal([]byte(str), &cond); e != nil {
		return nil, e
	}
	return queryFunc(cond), nil
}
