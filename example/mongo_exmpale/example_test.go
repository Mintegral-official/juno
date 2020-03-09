package mongo_exmpale

import (
	"context"
	"fmt"
	"github.com/MintegralTech/juno/builder"
	"github.com/MintegralTech/juno/check"
	"github.com/MintegralTech/juno/query"
	"github.com/MintegralTech/juno/search"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"testing"
	"time"
)

func BenchmarkSliceEqual(b *testing.B) {

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

	// search: advertiserId=457 or platform=android or (price in [20.0, 1.4, 3.6, 5.7, 2.5] And price >= 1.4)
	// invert list
	invertIdx := tIndex.GetInvertedIndex()

	// storage
	storageIdx := tIndex.GetStorageIndex()

	var p = []float64{2.3, 1.4, 3.65, 2.46, 2.5}
	var pi = make([]interface{}, len(p))
	for _, v := range p {
		pi = append(pi, v)
	}
	var a0 = []int64{647, 658, 670}
	var ai = make([]interface{}, len(a0))
	for _, v := range a0 {
		ai = append(ai, v)
	}

	var dev = []int64{4, 5}
	var devi = make([]interface{}, len(dev))
	devi = append(devi, dev)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		q := query.NewOrQuery([]query.Query{
			query.NewOrQuery([]query.Query{
				query.NewTermQuery(invertIdx.Iterator("Platform", "1")),
			}, nil),
			query.NewOrQuery([]query.Query{
				query.NewTermQuery(invertIdx.Iterator("AdvertiserId", "457")),
			}, nil),
			/* special example */
			query.NewOrQuery([]query.Query{
				query.NewTermQuery(storageIdx.Iterator("DeviceTypeV2")),
			}, []check.Checker{
				check.NewInChecker(storageIdx.Iterator("DeviceTypeV2"), devi, nil, false),
			}),
			query.NewAndQuery([]query.Query{
				query.NewAndQuery([]query.Query{
					query.NewTermQuery(storageIdx.Iterator("Price")),
				}, []check.Checker{
					check.NewInChecker(storageIdx.Iterator("Price"), pi, nil, false),
				}),
				query.NewAndQuery([]query.Query{
					query.NewTermQuery(storageIdx.Iterator("AdvertiserId")),
				}, []check.Checker{
					check.NewNotChecker(storageIdx.Iterator("AdvertiserId"), ai, nil, false),
				})}, nil)},
			nil,
		)

		tquery := time.Now()
		r1 := search.NewSearcher()
		r1.Search(tIndex, q)
		fmt.Println("query: ", time.Since(tquery))
		fmt.Println("+****************************+")
		fmt.Println("query res: ", len(r1.Docs), r1.Time)
		//fmt.Println("+****************************+")
		//fmt.Println(r1.QueryDebug)
		//fmt.Println("+****************************+")
		//fmt.Println(r1.IndexDebug)
		//fmt.Println("+****************************+")

		a := "AdvertiserId=457 or Platform=1 or (Price in [2.3, 1.4, 3.65, 2.46, 2.5] and AdvertiserId !in [647, 658, 670])"
		sq := query.NewSqlQuery(a, nil, false)

		tsql := time.Now()
		m := sq.LRD(tIndex)
		r2 := search.NewSearcher()
		r2.Search(tIndex, m)
		fmt.Println("sql: ", time.Since(tsql))
		//fmt.Println(r2.QueryDebug)
		//fmt.Println(r2.IndexDebug)
		fmt.Println("+****************************+")
		fmt.Println("sql res: ", len(r2.Docs), r2.Time)
	}
	b.StopTimer()
	tIndex = nil
}
