package juno

import (
	"github.com/Mintegral-official/juno/conf"
	"github.com/Mintegral-official/juno/document"
	"github.com/Mintegral-official/juno/helpers"
	"github.com/Mintegral-official/juno/index"
	"github.com/Mintegral-official/juno/model"
	"github.com/Mintegral-official/juno/query"
)

type IndexBuilder interface {
	build() *index.IndexImpl
}

type IndexBuilderImpl struct {
	Campaign []*model.CampaignInfo
}

func NewIndexBuilder(cfg *conf.MongoCfg) *IndexBuilderImpl {
	if cfg == nil {
		return nil
	}
	mon, err := model.NewMongo(cfg)
	if err != nil {
		return nil
	}
	c, err := mon.Find()
	if err != nil {
		return nil
	}
	return &IndexBuilderImpl{
		Campaign: c,
	}
}

func (ib *IndexBuilderImpl) CampaignFilter() []*model.CampaignInfo {
	if ib == nil {
		return nil
	}
	r := ib.Campaign
	c := make([]*model.CampaignInfo, len(r))
	for i := 0; i < len(r); i++ {
		if !r[i].IsSSPlatform() {
			continue
		}
		if int(*r[i].AdvertiserId) == 919 || int(*r[i].AdvertiserId) == 976 {
			continue
		}
		c = append(c, r[i])
	}
	return c
}

func (ib *IndexBuilderImpl) build() *index.IndexImpl {
	if ib == nil {
		return nil
	}
	c := ib.Campaign
	if c == nil || len(c) == 0 {
		return index.NewIndex("empty")
	}
	idx := index.NewIndex("")
	info := &document.DocInfo{
		Fields: []*document.Field{},
	}
	for i := 0; i < len(c); i++ {
		info.Id = document.DocId(c[i].CampaignId)
		info.Fields = []*document.Field{
			{
				Name:      "AdvertiserId",
				IndexType: 2,
				Value:     c[i].AdvertiserId,
			},
			{
				Name:      "Platform",
				IndexType: 2,
				Value:     c[i].Platform,
			},
			{
				Name:      "Price",
				IndexType: 2,
				Value:     c[i].Price,
			},
			{
				Name:      "StartTime",
				IndexType: 0,
				Value:     c[i].StartTime,
			},
			{
				Name:      "EndTime",
				IndexType: 0,
				Value:     c[i].EndTime,
			},
			{
				Name:      "PackageName",
				IndexType: 0,
				Value:     c[i].PackageName,
			},
			{
				Name:      "CampaignType",
				IndexType: 0,
				Value:     c[i].CampaignType,
			},
			{
				Name:      "OsVersionMaxV2",
				IndexType: 2,
				Value:     c[i].OsVersionMaxV2,
			},
			{
				Name:      "OsVersionMinV2",
				IndexType: 2,
				Value:     c[i].OsVersionMinV2,
			},
		}
		n := query.NewOperationImpl(*c[i].Price, 9, helpers.Float64Compare)
		v := []float64{1.05, 2.4, 0.57, 1.24, 1.05, 4.29}
		inf := make([]interface{}, len(v))
		for i, v := range v {
			inf[i] = v
		}
		if !n.In(inf) {
			continue
		}
		_ = idx.Add(info)
	}
	return idx
}
