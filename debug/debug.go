package debug

import (
	"encoding/json"
	"unsafe"
)

type Debug struct {
	Level int      `json:"-"`
	Name  string   `json:"name"`
	Msg   []string `json:"msg,omitempty"`
	Node  []*Debug `json:"node,omitempty"`
}

func NewDebug(level int, name string) *Debug {
	return &Debug{
		Level: level,
		Name:  name,
		Msg:   []string{},
		//	Node: []*Debug{},
	}
}

func (d *Debug) AddDebug(debug ...*Debug) {
	for _, v := range debug {
		d.Node = append(d.Node, v)
	}
}

func (d *Debug) AddDebugMsg(msg ...string) {
	for _, v := range msg {
		d.Msg = append(d.Msg, v)
	}
}

func (d *Debug) String() string {
	res, err := json.Marshal(d)
	if err == nil {
		return *(*string)(unsafe.Pointer(&res))
	} else {
		return err.Error()
	}
}
