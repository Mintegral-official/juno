package helpers

import "github.com/pkg/errors"

// Error codes returned by failures to parse an expression.
var (
	ElementNotfound    = errors.New("Element Not Found")
	ParseError         = errors.New("Parse Error")
	DocumentError      = errors.New("doc is nil")
	DocIdNotFound      = errors.New("DocId Not Found")
	NoMoreData         = errors.New("no more data")
	ComparableError    = errors.New("Comparable not nil")
	MongoCfgError      = errors.New("mongo config should not nil")
	ConnectError       = errors.New("database connect failed")
	PingError          = errors.New("connect ping is not through")
	CollectionNotFound = errors.New("collection not found")
	//CursorError        = errors.New("cursor wrong")
	//DelFailed          = errors.New("del failed")
)
