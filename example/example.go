package main

import (
	"fmt"
	"github.com/Mintegral-official/juno/builder"
	"github.com/Mintegral-official/juno/datastruct"
	"github.com/Mintegral-official/juno/document"
	"github.com/Mintegral-official/juno/example/model"
	"github.com/Mintegral-official/juno/index"
	"github.com/Mintegral-official/juno/query"
	"github.com/Mintegral-official/juno/query/check"
	"github.com/Mintegral-official/juno/query/operation"
	"go.mongodb.org/mongo-driver/bson"
	"time"
)

type IndexBuilderImpl struct {
	Campaign []*model.CampaignInfo
}

var (
	ib, ibInc *IndexBuilderImpl
	tIndex    *index.IndexImpl
	cfg       = &builder.MongoCfg{
		URI:            "mongodb://13.250.108.190:27017",
		DB:             "new_adn",
		Collection:     "campaign",
		ConnectTimeout: 10000,
		ReadTimeout:    20000,
	}
)

func NewIndexBuilder(cfg *builder.MongoCfg, m bson.M) *IndexBuilderImpl {
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
	t := time.Now()
	for i := 0; i < len(c); i++ {
		if info = index.MakeInfo(c[i]); info != nil {
			_ = idx.Add(info)
		}
	}
	fmt.Println("index build:", time.Since(t))
	//fmt.Println(idx.GetBitMap().Count())
	return idx
}

func buildIndex() {
	ib = NewIndexBuilder(cfg, bson.M{"status": 1})
	tIndex = ib.build()

	go func() {

	}()
}

func main() {

	// build index
	buildIndex()

	// search
	if1 := tIndex.GetStorageIndex().Iterator("AdvertiserId").(*datastruct.SkipListIterator)
	if2 := tIndex.GetStorageIndex().Iterator("Platform").(*datastruct.SkipListIterator)
	if3 := tIndex.GetStorageIndex().Iterator("Price").(*datastruct.SkipListIterator)
	if331 := tIndex.GetStorageIndex().Iterator("Price").(*datastruct.SkipListIterator)
	if332 := tIndex.GetStorageIndex().Iterator("Price").(*datastruct.SkipListIterator)
	if333 := tIndex.GetStorageIndex().Iterator("Price").(*datastruct.SkipListIterator)
	if334 := tIndex.GetStorageIndex().Iterator("Price").(*datastruct.SkipListIterator)

	q := query.NewOrQuery([]query.Query{
		query.NewTermQuery(if3),
		query.NewAndQuery([]query.Query{
			query.NewTermQuery(if1),
			query.NewTermQuery(if2),
			query.NewTermQuery(if3),
		}, nil),
	},
		[]check.Checker{
			check.NewCheckerImpl(if331, 20.0, operation.LT),
			check.NewCheckerImpl(if332, 16.4, operation.LE),
			check.NewCheckerImpl(if333, 0.5, operation.EQ),
			check.NewCheckerImpl(if334, 1.24, operation.EQ),
		},
	)
	res := tIndex.Search(q)

	fmt.Println("res: ", len(res.Docs), res.Time)
	for {
		inc := time.After(5 * time.Second)
		base := time.After(120 * time.Second)
		now := time.Now().Unix()
		select {
		case <-inc:
			fmt.Println("now & now: ", now, time.Now().Unix())
			ibInc = NewIndexBuilder(cfg, bson.M{"updated": bson.M{"$gt": now}})
			if ibInc == nil {
				fmt.Println("ibInc is nil")
				return
			}
		case <-base:
			ib = NewIndexBuilder(cfg, bson.M{"status": 1})
			tIndex = ib.build()
		}
		fmt.Println(len(ibInc.Campaign))
		if ibInc != nil && ibInc.Campaign != nil && len(ibInc.Campaign) != 0 {
			a := time.Now()
			tIndex.IncBuild(ibInc.Campaign)
			fmt.Println("index inc time: ", time.Since(a))
			var docs = make([]document.DocId, len(res.Docs))
			t := time.Now()
			for i := 0; i < len(res.Docs); i++ {
				if !tIndex.GetBitMap().IsExist(int(tIndex.GetCampaignMap()[res.Docs[i]])) {
					continue
				}
				docs[i] = res.Docs[i]
			}
			fmt.Println(time.Since(t))
			ibInc = nil
		}
		fmt.Println("res & resInc: ", len(res.Docs), res.Time)
	}

}
