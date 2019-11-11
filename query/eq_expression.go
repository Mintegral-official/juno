package query

type EqExpression struct {
	Op         OP
	FieldValue interface{}
}

func NewEqExpression(fieldValue interface{}) *EqExpression {
	return &EqExpression{
		Op:         0,
		FieldValue: fieldValue,
	}
}
