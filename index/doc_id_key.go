package index

import "github.com/Mintegral-official/juno/document"

type DocIdKey struct {
	value document.DocId
}

func (i *DocIdKey) PartitionKey() int64 {
	return int64(i.value)
}

func (i *DocIdKey) Value() interface{} {
	return i.value
}

func DocId(key document.DocId) *DocIdKey {
	return &DocIdKey{key}
}
