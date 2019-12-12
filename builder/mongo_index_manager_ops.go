package builder

import (
	"github.com/Mintegral-official/juno/document"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	DataModAdd = iota
	DataModUpdate
	DateModDel
)

type DataMod int

type ParserResult struct {
	DataMod DataMod
	Value   *document.DocInfo
}

type MongoParser interface {
	Parse([]byte, bool) (*ParserResult, error)
}

type MongoIndexManagerOps struct {
	Name           string
	IncInterval    int
	BaseInterval   int
	URI            string
	DB             string
	Collection     string
	ConnectTimeout int
	ReadTimeout    int
	BaseParser     MongoParser
	IncParser      MongoParser
	BaseQuery      interface{}
	IncQuery       interface{}
	UserData       interface{}
	FindOpt        *options.FindOptions
	OnBeforeBase   func(interface{}) interface{}
	OnBeforeInc    func(interface{}) interface{}
}
