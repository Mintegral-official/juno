package index

type DocId uint64
type FieldType int64
type IndexType int64

type Field struct {
	Name      string
	IndexType IndexType
	Value     interface{}
}

type DocInfo struct {
	Id     DocId
	Fields []*Field
}
