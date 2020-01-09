package operation

import "github.com/Mintegral-official/juno/helpers"

type Operations struct {
	FieldValue interface{}
}

func NewOperations(fieldValue interface{}) *Operations {
	if fieldValue == nil {
		return nil
	}
	return &Operations{
		FieldValue: fieldValue,
	}
}

func (ee *Operations) Equal(value interface{}) bool {
	if value == nil {
		return false
	}
	return helpers.Compare(ee.FieldValue, value) == 0
}

func (ee *Operations) Less(value interface{}) bool {
	if value == nil {
		return false
	}
	return helpers.Compare(ee.FieldValue, value) < 0
}

func (ee *Operations) In(value []interface{}) bool {
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

func (ee *Operations) SetValue(value interface{}) {
	ee.FieldValue = value
}