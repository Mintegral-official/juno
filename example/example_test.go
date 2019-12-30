package main

import (
	"context"
	"fmt"
	"github.com/Mintegral-official/juno/builder"
	"github.com/Mintegral-official/juno/query"
	"github.com/Mintegral-official/juno/query/check"
	"github.com/Mintegral-official/juno/search"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"testing"
	"time"
)

func BenchmarkCampaignParser_Parse(b *testing.B) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// build index
	bi, e := builder.NewMongoIndexBuilder(&builder.MongoIndexManagerOps{
		URI:            "mongodb://13.250.108.190:27017",
		IncInterval:    5,
		BaseInterval:   120,
		IncParser:      &CampaignParser{},
		BaseParser:     &CampaignParser{},
		BaseQuery:      bson.M{"status": 1},
		IncQuery:       bson.M{"updated": bson.M{"$gte": time.Now().Unix() - 5, "$lte": time.Now().Unix()}},
		DB:             "new_adn",
		Collection:     "campaign",
		ConnectTimeout: 10000,
		ReadTimeout:    20000,
		UserData:       &UserData{},
		Logger:         logrus.New(),
		OnBeforeInc: func(userData interface{}) interface{} {
			ud, ok := userData.(*UserData)
			if !ok {
				return nil
			}
			incQuery := bson.M{"updated": bson.M{"$gte": ud.upTime - 5, "$lte": time.Now().Unix()}}
			return incQuery
		},
	})
	if e != nil {
		fmt.Println(e)
		return
	}
	if e := bi.Build(ctx, "indexName"); e != nil {
		fmt.Println("build error", e.Error())
	}

	tIndex := bi.GetIndex()

	// search: advertiserId=457 or platform=android or (price in [20.0, 1.4, 3.6, 5.7, 2.5] And AdvertiserId not in [647, 658, 670])
	// invert list
	invertIdx := tIndex.GetInvertedIndex()

	// storage
	storageIdx := tIndex.GetStorageIndex()

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		q := query.NewOrQuery(
			[]query.Query{
				query.NewTermQuery(invertIdx.Iterator("AdvertiserId", 457)),
				query.NewTermQuery(invertIdx.Iterator("Platform", 1)),
				query.NewAndQuery(
					[]query.Query{
						query.NewTermQuery(storageIdx.Iterator("Price")),
						query.NewTermQuery(storageIdx.Iterator("Price")),
					},
					[]check.Checker{
						check.NewInChecker(storageIdx.Iterator("Price"), 20.0, 1.4, 3.6, 5.7, 2.5),
						check.NewNotChecker(storageIdx.Iterator("AdvertiserId"), 647, 658, 670),
					},
				),
			}, nil,
		)

		r := search.NewSearcher()
		r.Search(tIndex, q)
		//fmt.Println("+****************************+")
		fmt.Println("res: ", len(r.Docs), r.Time)
		//fmt.Println("+****************************+")
		//fmt.Println(r.QueryDebug)
		//fmt.Println("+****************************+")
		//fmt.Println(r.IndexDebug)
		//fmt.Println("+****************************+")

	}
	bi = nil
}
