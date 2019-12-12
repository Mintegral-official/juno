package helpers

import "github.com/pkg/errors"

var (
	ElementNotfound    = errors.New("Element Not Found")
	ParseError         = errors.New("Parse Error")
	DocumentError      = errors.New("doc is nil")
	DelFailed          = errors.New("del failed")
	DocIdNotFound      = errors.New("DocId Not Found")
	NoMoreData         = errors.New("no more data")
	ComparableError    = errors.New("Comparable not nil")
	MongoCfgError      = errors.New("mongo config should not nil")
	ConnectError       = errors.New("database connect failed")
	CollectionNotFound = errors.New("collection not found")
	CursorError        = errors.New("cursor wrong")
)
