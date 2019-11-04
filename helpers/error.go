package helpers

import "errors"

/**
 * @author: tangye
 * @Date: 2019/11/4 18:33
 * @Description:
 */

var (
	ERROR_ELEMENT_ERROR = errors.New("Element Not Found")
	ERROR_PARSE_ERROR = errors.New("Parse Error")
	ERROR_DOCUMENT_ERROR = errors.New("Parse Error")
)