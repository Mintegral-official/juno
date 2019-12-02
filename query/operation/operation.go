package operation

type Operation interface {
	Equal(value interface{}) bool
	Less(value interface{}) bool
	In(value []interface{}) bool
}

