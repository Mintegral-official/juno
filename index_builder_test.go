package juno

import (
	"fmt"
	"github.com/Mintegral-official/juno/conf"
	"github.com/Mintegral-official/juno/datastruct"
	"github.com/Mintegral-official/juno/query"
	"testing"
	"time"
)

func TestNewIndexBuilder(t *testing.T) {
	cfg := &conf.MongoCfg{
		URI:            "mongodb://localhost:27017",
		DB:             "new_adn",
		Collection:     "campaign",
		ConnectTimeout: 10000,
		ReadTimeout:    20000,
	}

	ib := NewIndexBuilder(cfg)
	if ib == nil {
		fmt.Println("*********")
	} else {
		ii := ib.build()
		if1 := ii.GetInvertedIndex().Iterator("AdvertiserId").(*datastruct.SkipListIterator)
		if2 := ii.GetInvertedIndex().Iterator("Platform").(*datastruct.SkipListIterator)
		if3 := ii.GetInvertedIndex().Iterator("Price").(*datastruct.SkipListIterator)
		//for if1.HasNext() {
		//	fmt.Println(if1.Next())
		//}
		t := time.Now()
		q := query.NewOrQuery([]query.Query{
			query.NewTermQuery(if3),
			query.NewAndQuery([]query.Query{
				query.NewTermQuery(if1),
				query.NewTermQuery(if2),
				query.NewTermQuery(if3),
			}, nil),
		},
			nil,
		)
		fmt.Println(time.Since(t))

        t = time.Now()
		res := ii.Search(q)
		fmt.Println(time.Since(t))
        fmt.Println(len(res.Docs))
	}
}
