package builder

import (
	"context"
	"github.com/MintegralTech/juno/index"
	"time"
)

type BuildInfo struct {
	TotalNumber int64            `json:"total_num"`
	ErrorNumber int64            `json:"error_num"`
	AddNum      int64            `json:"add_num,"`
	DeleteNum   int64            `json:"delete_num,"`
	MergeTime   time.Duration    `json:"merge_time,"`
	LastUpdated time.Time        `json:"last_updated,omitempty"`
	IndexInfo   *index.IndexInfo `json:"index_info,omitempty"`
}

type Builder interface {
	Build(ctx context.Context) error
	GetIndex() index.Index
	GetBuilderInfo() *BuildInfo
}
