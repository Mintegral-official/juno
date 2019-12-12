package builder

import (
	"github.com/Mintegral-official/juno/document"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoParser interface {
	Parse([]byte, interface{}) (*document.DocInfo, error)
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
