package query

type Equal interface {
	Equal(value interface{}) bool
}

type Less interface {
	Less(value interface{}) bool
}

type In interface {
	In(value []interface{}) bool
}
