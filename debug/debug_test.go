package debug

import (
	"encoding/json"
	"fmt"
	"testing"
)

func Test_Debug(t *testing.T) {
	debug := &Debug{
		Name: "debug1",
		Msg:  []string{"msg1", "msg2"},
		Node: &Debug{
			Name: "msg2",
			Msg:  []string{"msg3", "msg4"},
			Node: nil,
		},
	}

	if j, err := json.Marshal(debug); err == nil {
		fmt.Println(string(j))
	}

}
