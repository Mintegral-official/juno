package datastruct

type Iterator interface {
	HasNext() bool
	Next()
	Current() interface{}
	GetGE(id interface{}) interface{}
}
