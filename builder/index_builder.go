package builder

import (
	"github.com/Mintegral-official/juno/conf"
	"github.com/Mintegral-official/juno/document"
	"github.com/Mintegral-official/juno/index"
	"github.com/Mintegral-official/juno/model"
	"go.mongodb.org/mongo-driver/bson"
)

type IndexBuilder struct {
	Builder
	Campaign []*model.CampaignInfo
}

func NewIndexBuilder(cfg *conf.MongoCfg, m bson.M) *IndexBuilder {
	if cfg == nil {
		return nil
	}
	mon, err := model.NewMongo(cfg)

	if err != nil {
		return nil
	}
	c, err := mon.Find(m)

	if err != nil {
		return nil
	}
	return &IndexBuilder{
		Campaign: c,
	}
}

func (ib *IndexBuilder) CampaignFilter() []*model.CampaignInfo {
	if ib == nil {
		return nil
	}
	r := ib.Campaign
	// TODO
	return r
}

func (ib *IndexBuilder) build() *index.IndexImpl {

	if ib == nil || ib.Campaign == nil || len(ib.Campaign) == 0 {
		return index.NewIndex("empty")
	}

	c := ib.CampaignFilter() // 过滤 TODO
	idx := index.NewIndex("IndexBuilder")
	info := &document.DocInfo{
		Fields: []*document.Field{},
	}
	for i := 0; i < len(c); i++ {
		info = index.MakeInfo(c[i])
		if info == nil {
			continue
		}
		_ = idx.Add(info)
	}
	return idx
}
