package query

import (
	"github.com/Mintegral-official/juno/helpers"
)

type OperationImpl struct {
	Op         OP
	FieldValue interface{}
	cmp        helpers.Comparable
}

func NewOperationImpl(fieldValue interface{}, op OP, cmp helpers.Comparable) *OperationImpl {
	if fieldValue == nil || cmp == nil {
		return nil
	}
	return &OperationImpl{
		Op:         op,
		FieldValue: fieldValue,
		cmp:        cmp,
	}
}

func (ee *OperationImpl) Equal(value interface{}) bool {
	if value == nil {
		return false
	}
	return ee.cmp.Compare(ee.FieldValue, value) == 0
}

func (ee *OperationImpl) Less(value interface{}) bool {
	if value == nil {
		return false
	}
	return ee.cmp.Compare(ee.FieldValue, value) == -1
}

func (ee *OperationImpl) In(value []interface{}) bool {
	if value == nil {
		return false
	}
	for _, v := range value {
		if ee.cmp.Compare(ee.FieldValue, v) == 0 {
			return true
		}
	}
	return false
}
