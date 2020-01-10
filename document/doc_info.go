package document

type DocId uint64
type FieldType int64
type IndexType int64

const (
	InvertedIndexType = iota
	StorageIndexType
	BothIndexType
)

const (
	BoolFieldType FieldType = iota
	IntFieldType
	FloatFieldType
	StringFieldType
	SliceFieldType
	MapFieldType
	SelfDefinedFieldType
	DefaultFieldType = StringFieldType
)

type Field struct {
	Name      string
	IndexType IndexType
	Value     interface{}
	ValueType FieldType
}

type DocInfo struct {
	Id     DocId
	Fields []*Field
}
