package datastruct

type Iterator interface {
	HasNext() bool
	Next() interface{}
	Current() (interface{}, error)
	GetGE(id interface{}) interface{}
}
