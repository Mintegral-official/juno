package builder

import (
	"github.com/Mintegral-official/juno/index"
	"go.mongodb.org/mongo-driver/bson"
)

type IndexBuilder struct {
	*MongoIndexManager
}

func NewIndexBuilder(ops *MongoIndexManagerOps) *IndexBuilder {
	if ops == nil {
		return nil
	}
	mongoIndexManager := NewMongoIndexManager(ops)
	return &IndexBuilder{
		mongoIndexManager,
	}
}

func (ib *IndexBuilder) filter() []*ParserResult {
	if ib == nil {
		return nil
	}
	c := ib.result // TODO
	return c
}

func (ib *IndexBuilder) Build() *index.IndexImpl {
	_ = ib.find(bson.M{"status": 1})
	if ib == nil || ib.result == nil || len(ib.result) == 0 {
		return index.NewIndex("empty")
	}
	ib.innerIndex = index.NewIndex("index")
	c := ib.result
	for i := 0; i < len(c); i++ {
		_ = ib.innerIndex.Add(c[i].Value)
	}
	return ib.innerIndex
}
