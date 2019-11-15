package helpers

import "github.com/pkg/errors"

/**
 * @author: tangye
 * @Date: 2019/11/4 18:33
 * @Description:
 */

var (
	ElementNotfound = errors.New("Element Not Found")
	ParseError      = errors.New("Parse Error")
	DocumentError   = errors.New("doc is nil")
	DelFailed       = errors.New("del failed")
	DocIdNotFound   = errors.New("DocId Not Found")
	NoMoreData      = errors.New("no more data")
)
