package operation

import "github.com/MintegralTech/juno/helpers"

type Operations struct {
	FieldValue interface{}
}

func NewOperations(fieldValue interface{}) *Operations {
	return &Operations{
		FieldValue: fieldValue,
	}
}

func (ee *Operations) Equal(value interface{}) bool {
	return helpers.Compare(ee.FieldValue, value) == 0
}

func (ee *Operations) Less(value interface{}) bool {
	return helpers.Compare(ee.FieldValue, value) < 0
}

func (ee *Operations) In(value interface{}) bool {
	switch value.(type) {
	case []int:
		for _, v := range value.([]int) {
			if helpers.Compare(ee.FieldValue, v) == 0 {
				return true
			}
		}
	case []int32:
		for _, v := range value.([]int32) {
			if helpers.Compare(ee.FieldValue, v) == 0 {
				return true
			}
		}
	case []int64:
		for _, v := range value.([]int64) {
			if helpers.Compare(ee.FieldValue, v) == 0 {
				return true
			}
		}
	case []float32:
		for _, v := range value.([]float32) {
			if helpers.Compare(ee.FieldValue, v) == 0 {
				return true
			}
		}
	case []float64:
		for _, v := range value.([]float64) {
			if helpers.Compare(ee.FieldValue, v) == 0 {
				return true
			}
		}
	case []string:
		for _, v := range value.([]string) {
			if helpers.Compare(ee.FieldValue, v) == 0 {
				return true
			}
		}
	}
	return false
}

func (ee *Operations) SetValue(value interface{}) {
	ee.FieldValue = value
}