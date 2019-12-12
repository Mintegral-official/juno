package main

import (
	"fmt"
	"github.com/Mintegral-official/juno/builder"
	"github.com/Mintegral-official/juno/datastruct"
	"github.com/Mintegral-official/juno/index"
	"github.com/Mintegral-official/juno/query"
	"github.com/Mintegral-official/juno/query/check"
	"github.com/Mintegral-official/juno/query/operation"
	"go.mongodb.org/mongo-driver/bson"
	"sync"
	"testing"
)

func BenchmarkIndexBuilderImpl_CampaignFilter(b *testing.B) {
	cfg := &builder.MongoCfg{
		URI:            "mongodb://192.168.1.198:27017",
		DB:             "new_adn",
		Collection:     "campaign",
		ConnectTimeout: 10000,
		ReadTimeout:    20000,
	}

	ib := NewIndexBuilder(cfg, bson.M{"status": 1})

	ii := ib.build()
	if1 := ii.GetStorageIndex().Iterator("AdvertiserId").(*datastruct.SkipListIterator)
	if2 := ii.GetStorageIndex().Iterator("Platform").(*datastruct.SkipListIterator)
	if3 := ii.GetStorageIndex().Iterator("Price").(*datastruct.SkipListIterator)

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
	var res *index.SearchResult
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		res = ii.Search(q)
	}
	b.StopTimer()
	if res != nil {
		fmt.Println(len(res.Docs))
		fmt.Println(res.Time)
	}
}

func TestNewIndexBuilder(t *testing.T) {
	var m sync.Map
	m.Store(1, 1)
	m.Store(2, 2)
	m.Store(3, 3)
	count := 1
	m.Range(func(key, value interface{}) bool {
		if key != nil {
			count++
			return true
		}
		return false
	})
	fmt.Println(count)
}
