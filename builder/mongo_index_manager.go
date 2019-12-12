package builder

import "github.com/Mintegral-official/juno/index"

type MongoIndexManager struct {
	ops        *MongoIndexManagerOps
	innerIndex index.Index
}

func NewMongoIndexManager(ops *MongoIndexManagerOps) *MongoIndexManager {
	return &MongoIndexManager{
		ops:        ops,
		innerIndex: nil,
	}
}

func (mim *MongoIndexManager) Update() error {
	if e := mim.base(); e != nil {
		return e
	}
	// TODO
	return nil
}

func (mim *MongoIndexManager) base() error { // TODO
	// TODO
	return nil
}

func (mim *MongoIndexManager) inc() error {
	// TODO
	return nil
}
