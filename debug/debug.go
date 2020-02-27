package debug

import (
	"encoding/json"
	"github.com/Mintegral-official/juno/document"
	"unsafe"
)

type Debug struct {
	//Name document.DocId                `json:"name"`
	Node map[document.DocId][][]string `json:"node"`
}

func NewDebug(name string) *Debug {
	return &Debug{
		Node: map[document.DocId][][]string{},
	}
}

func (d *Debug) AddDebug(debug ...*Debug) {
	//for _, v := range debug {
	//	d.Node = append(d.Node, v)
	//}
}

func (d *Debug) AddDebugMsg(msg ...string) {
	//for _, v := range msg {
	//	d.Msg = append(d.Msg, v)
	//}
}

func (d *Debug) String() string {
	res, err := json.Marshal(d)
	if err == nil {
		return *(*string)(unsafe.Pointer(&res))
	} else {
		return err.Error()
	}
}
