package datastruct

type Iterator interface {
	HasNext() bool
	Next() interface{}
	Current() interface{}
	GetGE(id interface{}) interface{}
}
