package query

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestMarshalQueryInfo_Marshal(t *testing.T) {
	var a interface{} = 4
	var b int = 100
	r, e := json.Marshal(b)
	fmt.Println(string(r), e)
	_ = json.Unmarshal(r, &a)
	fmt.Println(a)

}
