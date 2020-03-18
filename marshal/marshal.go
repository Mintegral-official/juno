package marshal

import "github.com/MintegralTech/juno/operation"

type MarshalInfo struct {
	Name          string         `json:"name,omitempty"`
	QueryValue    interface{}    `json:"query_value,omitempty"`
	IndexValue    interface{}    `json:"index_value,omitempty"`
	Operation     string         `json:"operation,omitempty"`
	Op            operation.OP   `json:"op,omitempty"`
	Transfer      bool           `json:"transfer,omitempty"`
	SelfOperation bool           `json:"self_operation,omitempty"`
	Label         string         `json:"label,omitempty"`
	Result        bool           `json:"result,omitempty"`
	Nodes         []*MarshalInfo `json:"nodes,omitempty"`
}
