package query

import (
	"encoding/json"
)

type JSONFormatter struct {
}

func NewJSONFormatter() *JSONFormatter {
	return &JSONFormatter{}
}

func (j *JSONFormatter) Marshal(queryInfo map[string]interface{}) (string, error) {
	if r, e := json.Marshal(queryInfo); e != nil {
		return "", e
	} else {
		return string(r), nil
	}

}

func (j *JSONFormatter) Unmarshal(str string) (map[string]interface{}, error) {
	var res map[string]interface{}
	if e := json.Unmarshal([]byte(str), &res); e != nil {
		return nil, e
	}
	return res, nil
}
