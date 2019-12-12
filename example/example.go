package main

import (
	"fmt"
	"github.com/Mintegral-official/juno/conf"
	"github.com/Mintegral-official/juno/datastruct"
	"github.com/Mintegral-official/juno/document"
	"github.com/Mintegral-official/juno/index"
	"github.com/Mintegral-official/juno/model"
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
	ii        *index.IndexImpl
	cfg       = &conf.MongoCfg{
		URI:            "mongodb://13.250.108.190:27017",
		DB:             "new_adn",
		Collection:     "campaign",
		ConnectTimeout: 10000,
		ReadTimeout:    20000,
	}
)

func init() {
	ib = NewIndexBuilder(cfg, bson.M{"status": 1})
	ii = ib.build()
}

func NewIndexBuilder(cfg *conf.MongoCfg, m bson.M) *IndexBuilderImpl {
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

func main() {

	if1 := ii.GetStorageIndex().Iterator("AdvertiserId").(*datastruct.SkipListIterator)
	if2 := ii.GetStorageIndex().Iterator("Platform").(*datastruct.SkipListIterator)
	if3 := ii.GetStorageIndex().Iterator("Price").(*datastruct.SkipListIterator)
	//if4 := ii.GetInvertedIndex().Iterator("Price_1.5").(*datastruct.SkipListIterator)
	if331 := ii.GetStorageIndex().Iterator("Price").(*datastruct.SkipListIterator)
	if332 := ii.GetStorageIndex().Iterator("Price").(*datastruct.SkipListIterator)
	if333 := ii.GetStorageIndex().Iterator("Price").(*datastruct.SkipListIterator)
	if334 := ii.GetStorageIndex().Iterator("Price").(*datastruct.SkipListIterator)

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
	//now := time.Now().Unix()
	//time.Sleep(5 * time.Second)
	//ibInc = NewIndexBuilder(cfg, bson.M{"updated": bson.M{"$gt": now}})
	//if ibInc == nil {
	//	fmt.Println("ibInc is nil")
	//	return
	//}
	//// iiInc := ibInc.build()
	//fmt.Println("inc change: ", len(ibInc.Campaign))
	res := ii.Search(q)
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
			ii = ib.build()
		}
		fmt.Println(len(ibInc.Campaign))
		if ibInc != nil && ibInc.Campaign != nil && len(ibInc.Campaign) != 0 {
			a := time.Now()
			ii.IncBuild(ibInc.Campaign)
			fmt.Println("index inc time: ", time.Since(a))
			var docs = make([]document.DocId, len(res.Docs))
			t := time.Now()
			for i := 0; i < len(res.Docs); i++ {
				if !ii.GetBitMap().IsExist(int(ii.GetCampaignMap()[res.Docs[i]])) {
					continue
				}
				docs[i] = res.Docs[i]
			}
			fmt.Println(time.Since(t))
			ibInc = nil
		}
		fmt.Println("res & resInc: ", len(res.Docs), res.Time)
	}

	//if11 := iiInc.GetStorageIndex().Iterator("AdvertiserId").(*datastruct.SkipListIterator)
	//if22 := iiInc.GetStorageIndex().Iterator("Platform").(*datastruct.SkipListIterator)
	//if33 := iiInc.GetStorageIndex().Iterator("Price").(*datastruct.SkipListIterator)
	////if44 := ii.GetInvertedIndex().Iterator("Price_1.5").(*datastruct.SkipListIterator)
	//
	//if3311 := iiInc.GetStorageIndex().Iterator("Price").(*datastruct.SkipListIterator)
	//if3322 := iiInc.GetStorageIndex().Iterator("Price").(*datastruct.SkipListIterator)
	//if3333 := iiInc.GetStorageIndex().Iterator("Price").(*datastruct.SkipListIterator)
	//if3344 := iiInc.GetStorageIndex().Iterator("Price").(*datastruct.SkipListIterator)
	//
	//q = query.NewOrQuery([]query.Query{
	//	query.NewTermQuery(if33),
	//	query.NewAndQuery([]query.Query{
	//		query.NewTermQuery(if11),
	//		query.NewTermQuery(if22),
	//		query.NewTermQuery(if33),
	//	}, nil),
	//},
	//	[]check.Checker{
	//		check.NewCheckerImpl(if3311, 20.0, operation.LT),
	//		check.NewCheckerImpl(if3322, 16.4, operation.LE),
	//		check.NewCheckerImpl(if3333, 0.5, operation.EQ),
	//		check.NewCheckerImpl(if3344, 1.24, operation.EQ),
	//	},
	//)
	//resInc := iiInc.Search(q)
	//if resInc == nil {
	//	res.Docs = docs
	//} else {
	//	fmt.Println(len(docs))
	//	docs = helpers.Merge(docs, resInc.Docs)
	//	res.Docs = docs
	//	res.Time = res.Time + resInc.Time
	//	fmt.Println("resInc: ", len(resInc.Docs), resInc.Time, len(docs))
	//}

}
