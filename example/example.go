package main

import (
	"context"
	"fmt"
	"github.com/Mintegral-official/juno/builder"
	"github.com/Mintegral-official/juno/datastruct"
	"github.com/Mintegral-official/juno/document"
	"github.com/Mintegral-official/juno/example/model"
	"github.com/Mintegral-official/juno/query"
	"github.com/Mintegral-official/juno/query/check"
	"github.com/Mintegral-official/juno/query/operation"
	"go.mongodb.org/mongo-driver/bson"
	"os"
	"os/signal"
	"time"
)

type CampaignParser struct {
}

func MakeInfo(info *model.CampaignInfo) *document.DocInfo {
	if info == nil {
		return nil
	}
	docInfo := &document.DocInfo{
		Fields: []*document.Field{},
	}
	docInfo.Id = document.DocId(info.CampaignId)
	docInfo.Fields = []*document.Field{
		{
			Name:      "AdvertiserId",
			IndexType: 1,
			Value:     info.AdvertiserId,
		},
		{
			Name:      "Platform",
			IndexType: 2,
			Value:     info.Platform,
		},
		{
			Name:      "Price",
			IndexType: 1,
			Value:     *info.Price,
		},
		{
			Name:      "StartTime",
			IndexType: 1,
			Value:     info.StartTime,
		},
		{
			Name:      "EndTime",
			IndexType: 1,
			Value:     info.EndTime,
		},
		{
			Name:      "PackageName",
			IndexType: 1,
			Value:     info.PackageName,
		},
		{
			Name:      "CampaignType",
			IndexType: 2,
			Value:     info.CampaignType,
		},
		{
			Name:      "OsVersionMaxV2",
			IndexType: 1,
			Value:     info.OsVersionMaxV2,
		},
		{
			Name:      "OsVersionMinV2",
			IndexType: 1,
			Value:     info.OsVersionMinV2,
		},
	}
	return docInfo
}

func (c *CampaignParser) Parse(bytes []byte, flag bool) (*builder.ParserResult, error) {
	campaign := &model.CampaignInfo{}
	if err := bson.Unmarshal(bytes, &campaign); err != nil {
		fmt.Println("bson.Unmarsnal error:" + err.Error())
	}
	var info = MakeInfo(campaign)
	var mode builder.DataMod
	if flag {
		return &builder.ParserResult{
			DataMod: 0,
			Value:   info,
		}, nil
	}
	if campaign.Status == 1 {
		mode = 1

	} else {
		mode = 2
	}
	return &builder.ParserResult{
		DataMod: mode,
		Value:   info,
	}, nil
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// build index
	b := builder.NewMongoIndexBuilder(&builder.MongoIndexManagerOps{
		URI:            "mongodb://13.250.108.190:27017",
		IncInterval:    5,
		BaseInterval:   120,
		IncParser:      &CampaignParser{},
		BaseParser:     &CampaignParser{},
		BaseQuery:      bson.M{"status": 1},
		IncQuery:       bson.M{"updated": bson.M{"$gt": time.Now().Unix() - int64(5*time.Second)}},
		DB:             "new_adn",
		Collection:     "campaign",
		ConnectTimeout: 10000,
		ReadTimeout:    20000,
	})
	if b == nil {
		fmt.Println("build index error")
		return
	}
	if e := b.Build(ctx); e != nil {
		fmt.Println("build error", e.Error())
	}
	tIndex := b.GetIndex()

	// search: advertiserId=123 or (advertiserId=456 and platform=ios and (price <= 20.0 or price <= 16.4 or price = 0.5 or price = 1.24))
	// invert list
	if1 := tIndex.GetStorageIndex().Iterator("AdvertiserId=123").(*datastruct.SkipListIterator)
	if2 := tIndex.GetStorageIndex().Iterator("Platform=ios").(*datastruct.SkipListIterator)

	// storage
	storageIdx := tIndex.GetStorageIndex()
	if3 := tIndex.GetStorageIndex().Iterator("Price").(*datastruct.SkipListIterator)

	q := query.NewOrQuery([]query.Query{
		query.NewTermQuery(if3),
		query.NewAndQuery([]query.Query{
			query.NewTermQuery(if1),
			query.NewTermQuery(if2),
			query.NewTermQuery(if3),
		}, nil),
	},
		[]check.Checker{
			check.NewCheckerImpl(storageIdx.Iterator("Price").(*datastruct.SkipListIterator), 20.0, operation.LT),
			check.NewCheckerImpl(storageIdx.Iterator("Price").(*datastruct.SkipListIterator), 16.4, operation.LE),
			check.NewCheckerImpl(storageIdx.Iterator("Price").(*datastruct.SkipListIterator), 0.5, operation.EQ),
			check.NewCheckerImpl(storageIdx.Iterator("Price").(*datastruct.SkipListIterator), 1.24, operation.EQ),
		},
	)
	res := tIndex.Search(q)
	fmt.Println("res: ", len(res.Docs), res.Time)

	c := make(chan os.Signal)
	signal.Notify(c)
	s := <-c
	fmt.Println("退出信号", s)
}
