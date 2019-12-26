package debug

type Debug struct {
	Name string   `json:"name"`
	Msg  []string `json:"msg"`
	Node *Debug   `json:"node"`
}

func NewDebug(name string) *Debug {
	return &Debug{
		Name: name,
		Msg:  []string{},
	}
}

func (d *Debug) AddDebug(msg string) {
	d.Msg = append(d.Msg, msg)
}
