package query

import (
	"bytes"
	"encoding/gob"
	"unsafe"
)

type MarshalQueryInfo struct {
	queries Query
}

func (q *MarshalQueryInfo) Marshal(queries Query) (string, error) {
	buf := new(bytes.Buffer)
	enc := gob.NewEncoder(buf)
	if err := enc.Encode(queries); err != nil {
		return "", err
	} else {
		return buf.String(), nil
	}
}

func (q *MarshalQueryInfo) UnMarshal(str string) (Query, error) {
	var res = &MarshalQueryInfo{}
	dec := gob.NewDecoder(bytes.NewBuffer(*(*[]byte)(unsafe.Pointer(&str))))
	if err := dec.Decode(&res); err != nil {
		return nil, err
	} else {
		return res.queries, nil
	}
}
