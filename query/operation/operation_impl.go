package operation

import "github.com/Mintegral-official/juno/helpers"

type OperationImpl struct {
	FieldValue interface{}
}

func NewOperationImpl(fieldValue interface{}) *OperationImpl {
	if fieldValue == nil {
		return nil
	}
	return &OperationImpl{
		FieldValue: fieldValue,
	}
}

func (ee *OperationImpl) Equal(value interface{}) bool {
	if value == nil {
		return false
	}
	return helpers.Compare(ee.FieldValue, value) == 0
}

func (ee *OperationImpl) Less(value interface{}) bool {
	if value == nil {
		return false
	}
	return helpers.Compare(ee.FieldValue, value) == -1
}

func (ee *OperationImpl) In(value []interface{}) bool {
	if value == nil {
		return false
	}
	for _, v := range value {
		if helpers.Compare(ee.FieldValue, v) == 0 {
			return true
		}
	}
	return false
}
