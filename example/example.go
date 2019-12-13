package main

import (
	"context"
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
	"os"
	"os/signal"
	"time"
)

type CampaignParser struct {
}

var (
	ib     *builder.IndexBuilder
	tIndex *index.IndexImpl
	cfg    = &builder.MongoIndexManagerOps{
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
	}
)

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
		if tIndex.GetBitMap().IsExist(int(tIndex.GetCampaignMap()[document.DocId(campaign.CampaignId)])) {
			mode = 1
		} else {
			mode = 0
		}
	} else {
		mode = 2
	}
	return &builder.ParserResult{
		DataMod: mode,
		Value:   info,
	}, nil
}

func buildIndex() {
	ib = builder.NewIndexBuilder(cfg)
	tIndex = ib.Build()
}

func main() {

	// build index
	buildIndex()

	// search
	//advertiserId or (advertiserId and platform and price and (price <= 20.0 or price <= 16.4 or price = 0.5 or price = 1.24))
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
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	c := make(chan os.Signal)
	signal.Notify(c)
	_ = ib.MongoIndexManager.Update(ctx)
	res1 := tIndex.Search(q)
	fmt.Println("res1: ", len(res1.Docs), res1.Time)
	s := <-c
	fmt.Println("退出信号", s)
}
