package query

import (
	"encoding/json"
	"fmt"
	"github.com/Mintegral-official/juno/datastruct"
	"github.com/Mintegral-official/juno/document"
	"reflect"
	"runtime"
	"strings"
	"testing"
)

func TestMarshalQueryInfo_Marshal(t *testing.T) {
	var a interface{} = 4
	var b int = 100
	r, e := json.Marshal(b)
	fmt.Println(string(r), e)
	_ = json.Unmarshal(r, &a)
	fmt.Println(a)
	q := NewAndQuery([]Query{}, nil)

	fmt.Println(GetFunctionName(NewAndQuery, '/'))
	fmt.Println(GetFunctionName(q))
}

func TestNewAndQuery(t *testing.T) {
	sl := datastruct.NewSkipList(datastruct.DefaultMaxLevel)

	sl.Add(document.DocId(1), [1]byte{})
	sl.Add(document.DocId(3), [1]byte{})
	sl.Add(document.DocId(6), [1]byte{})
	sl.Add(document.DocId(10), [1]byte{})

	sl1 := datastruct.NewSkipList(datastruct.DefaultMaxLevel)

	sl1.Add(document.DocId(1), [1]byte{})
	sl1.Add(document.DocId(4), [1]byte{})
	sl1.Add(document.DocId(6), [1]byte{})
	sl1.Add(document.DocId(9), [1]byte{})

	a := NewAndQuery([]Query{NewTermQuery(sl.Iterator())}, nil)
	rs := reflect.ValueOf(a)
	fmt.Println("result:", rs.String())
	fmt.Println(reflect.ValueOf(NewAndQuery).String())

}

func GetFunctionName(i interface{}, seps ...rune) string {
	// 获取函数名称
	fn := runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()

	// 用 seps 进行分割
	fields := strings.FieldsFunc(fn, func(sep rune) bool {
		for _, s := range seps {
			if sep == s {
				return true
			}
		}
		return false
	})

	// fmt.Println(fields)

	if size := len(fields); size > 0 {
		return fields[size-1]
	}
	return ""
}
