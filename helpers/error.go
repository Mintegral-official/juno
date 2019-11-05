package helpers

import "errors"

/**
 * @author: tangye
 * @Date: 2019/11/4 18:33
 * @Description:
 */

var (
	ELEMENT_NOTFOUND = errors.New("Element Not Found")
	PARSE_ERROR = errors.New("Parse Error")
	DOCUMENT_ERROR = errors.New("doc is nil")
	DEL_FAILED = errors.New("del failed")

)