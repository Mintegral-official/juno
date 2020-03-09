package register

import "github.com/MintegralTech/juno/operation"

var FieldMap = make(map[string]operation.Operation, 16)

type Register struct {
}

func NewRegister() *Register {
	return &Register{}
}

func (r *Register) Register(fieldName string, e operation.Operation) {
	FieldMap[fieldName] = e
}
