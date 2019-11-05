package document

type DocId uint64
type FieldType int64
type IndexType int64

const (
	INVERTED_INDEX_TYPE = iota
	STORAGE_INDEX_TYPE
)

type Field struct {
	Name      string
	IndexType IndexType
	Value     interface{}
}

type DocInfo struct {
	Id     DocId
	Fields []*Field
}
