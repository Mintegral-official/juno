package datastruct

type Iterator interface {
	HasNext() bool
	Next() interface{}
	GetGE(id interface{}) interface{}
}
