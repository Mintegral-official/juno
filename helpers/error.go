package helpers

import "errors"

// Error codes returned by failures to parse an expression.
var (
	ElementNotfound    = errors.New("element not found")
	ParseError         = errors.New("parse error")
	DocumentError      = errors.New("doc is nil")
	NoMoreData         = errors.New("no more data")
	MongoCfgError      = errors.New("mongo config should not nil")
	ConnectError       = errors.New("database connect failed")
	PingError          = errors.New("connect ping is not through")
	CollectionNotFound = errors.New("collection not found")
	//CursorError        = errors.New("cursor wrong")
	//DelFailed          = errors.New("del failed")
)
