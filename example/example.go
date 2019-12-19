package main

import (
	"context"
	"fmt"
	"github.com/Mintegral-official/juno/builder"
	"github.com/Mintegral-official/juno/datastruct"
	"github.com/Mintegral-official/juno/document"
	"github.com/Mintegral-official/juno/query"
	"github.com/Mintegral-official/juno/query/check"
	"github.com/Mintegral-official/juno/query/operation"
	"github.com/Mintegral-official/juno/search"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"os"
	"os/signal"
	"time"
)

type CampaignInfo struct {
	CampaignId     int64    `bson:"campaignId,omitempty" json:"campaignId,omitempty"`
	AdvertiserId   *int32   `bson:"advertiserId,omitempty" json:"advertiserId,omitempty"`
	Price          *float64 `bson:"price,omitempty" json:"price,omitempty"`
	Status         int32    `bson:"status,omitempty" json:"status,omitempty"`
	PackageName    string   `bson:"packageName,omitempty" json:"packageName,omitempty"`
	CampaignType   *int32   `bson:"campaignType,omitempty" json:"campaignType,omitempty"`
	Platform       *int32   `bson:"platform,omitempty" json:"platform,omitempty"`
	OsVersionMinV2 *int     `bson:"oVersionMinV2,omitempty" json:"osVersionMinV2,omitempty"`
	OsVersionMaxV2 *int     `bson:"osVersionMaxV2,omitempty" json:"osVersionMaxV2,omitempty"`
	StartTime      *int     `bson:"startTime,omitempty" json:"startTime,omitempty"`
	EndTime        *int     `bson:"endTime,omitempty" json:"endTime,omitempty"`
}

type CampaignParser struct {
}

func MakeInfo(info *CampaignInfo) *document.DocInfo {
	if info == nil {
		return nil
	}
	docInfo := &document.DocInfo{
		Fields: []*document.Field{},
	}
	docInfo.Id = document.DocId(info.CampaignId)
	if info.Price == nil {
		return docInfo
	}
	docInfo.Fields = []*document.Field{
		{
			Name:      "AdvertiserId",
			IndexType: 2,
			Value:     info.AdvertiserId,
		},
		{
			Name:      "Platform",
			IndexType: 2,
			Value:     info.Platform,
		},
		{
			Name:      "Price",
			IndexType: 2,
			Value:     *info.Price,
		},
		{
			Name:      "Platform",
			IndexType: 2,
			Value:     info.Platform,
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
			IndexType: 2,
			Value:     info.PackageName,
		},
		{
			Name:      "CampaignType",
			IndexType: 1,
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

func (c *CampaignParser) Parse(bytes []byte) (*builder.ParserResult, error) {
	campaign := &CampaignInfo{}
	if err := bson.Unmarshal(bytes, &campaign); err != nil {
		return nil, err
	}
	var info = MakeInfo(campaign)
	var mode builder.DataMod = builder.DataDel
	if campaign.Status == 1 {
		mode = builder.DataAddOrUpdate
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
	b, e := builder.NewMongoIndexBuilder(&builder.MongoIndexManagerOps{
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
		Logger:         logrus.New(),
	})
	if e != nil {
		fmt.Println(e)
		return
	}
	if e := b.Build(ctx, "indexName"); e != nil {
		fmt.Println("build error", e.Error())
	}
	tIndex := b.GetIndex()

	// search: advertiserId=457 or platform=android or (price < 20.0 And price >= 16.4) or advertiserId=646
	// invert list
	invertIdx := tIndex.GetInvertedIndex()

	// storage
	storageIdx := tIndex.GetStorageIndex()

	q := query.NewOrQuery([]query.Query{
		query.NewTermQuery(invertIdx.Iterator("AdvertiserId_457").(*datastruct.SkipListIterator)),
		query.NewTermQuery(invertIdx.Iterator("Platform_1").(*datastruct.SkipListIterator)),
		query.NewAndQuery([]query.Query{
			query.NewTermQuery(tIndex.GetStorageIndex().Iterator("Price").(*datastruct.SkipListIterator)),
			query.NewTermQuery(tIndex.GetStorageIndex().Iterator("Price").(*datastruct.SkipListIterator)),
		},
			[]check.Checker{
				check.NewCheckerImpl(storageIdx.Iterator("Price").(*datastruct.SkipListIterator), 20.0, operation.LT),
				check.NewCheckerImpl(storageIdx.Iterator("Price").(*datastruct.SkipListIterator), 16.4, operation.GE),
			}),
		query.NewTermQuery(invertIdx.Iterator("AdvertiserId_646").(*datastruct.SkipListIterator)),
	}, nil)

	res := search.Search(tIndex, q)
	fmt.Println("res: ", len(res.Docs), res.Time)

	c := make(chan os.Signal)
	signal.Notify(c)
	s := <-c
	fmt.Println("退出信号", s)
}
