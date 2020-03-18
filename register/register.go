package register

import "github.com/MintegralTech/juno/operation"

var FieldMap map[string]operation.Operation

type Register struct {
}

func NewRegister() *Register {
	FieldMap = make(map[string]operation.Operation, 16)
	return &Register{}
}

func (r *Register) Register(fieldName string, e operation.Operation) {
	FieldMap[fieldName] = e
}
