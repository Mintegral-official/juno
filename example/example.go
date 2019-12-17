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
	"github.com/Mintegral-official/juno/search"
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
	docInfo.Fields = append(docInfo.Fields,
		&document.Field{
			Name:
			"AdvertiserId",
			IndexType: 2,
			Value:     info.AdvertiserId,
		})
	docInfo.Fields = append(docInfo.Fields,
		&document.Field{
			Name:      "Platform",
			IndexType: 2,
			Value:     info.Platform,
		})
	if info.Price != nil {
		docInfo.Fields = append(docInfo.Fields,
			&document.Field{
				Name:      "Price",
				IndexType: 1,
				Value:     *info.Price,
			})
	}

	docInfo.Fields = append(docInfo.Fields,
		&document.Field{
			Name:      "StartTime",
			IndexType: 1,
			Value:     info.StartTime,
		})

	docInfo.Fields = append(docInfo.Fields,
		&document.Field{
			Name:      "EndTime",
			IndexType: 1,
			Value:     info.EndTime,
		})

	docInfo.Fields = append(docInfo.Fields,
		&document.Field{
			Name:      "PackageName",
			IndexType: 2,
			Value:     info.PackageName,
		})

	docInfo.Fields = append(docInfo.Fields,
		&document.Field{
			Name:      "CampaignType",
			IndexType: 1,
			Value:     info.CampaignType,
		})

	docInfo.Fields = append(docInfo.Fields,
		&document.Field{
			Name:      "OsVersionMaxV2",
			IndexType: 1,
			Value:     info.OsVersionMaxV2,
		})

	docInfo.Fields = append(docInfo.Fields,
		&document.Field{
			Name:      "OsVersionMinV2",
			IndexType: 1,
			Value:     info.OsVersionMinV2,
		})
	return docInfo
}

func (c *CampaignParser) Parse(bytes []byte) (*builder.ParserResult, error) {
	campaign := &model.CampaignInfo{}
	if err := bson.Unmarshal(bytes, &campaign); err != nil {
		return nil, err
	}
	var info = MakeInfo(campaign)
	var mode builder.DataMod
	if campaign.Status == 1 {
		mode = 1
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
	})
	if e != nil {
		fmt.Println(e)
		return
	}
	if e := b.Build(ctx); e != nil {
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
