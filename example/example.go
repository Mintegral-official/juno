package main

import (
	"context"
	"fmt"
	"github.com/Mintegral-official/juno/builder"
	"github.com/Mintegral-official/juno/document"
	"github.com/Mintegral-official/juno/helpers"
	"github.com/Mintegral-official/juno/query"
	"github.com/Mintegral-official/juno/query/check"
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
	DeviceTypeV2   []int64  `bson:"deviceTypeV2,omitempty" json:"deviceTypeV2,omitempty"`
	Uptime         int64    `bson:"updated,omitempty"`
}

type CampaignParser struct {
}

type UserData struct {
	upTime int64
}

func MakeInfo(info *CampaignInfo) *document.DocInfo {
	if info == nil {
		return nil
	}
	docInfo := &document.DocInfo{
		Fields: []*document.Field{},
	}
	docInfo.Id = document.DocId(info.CampaignId)
	if info.AdvertiserId != nil {
		docInfo.Fields = append(docInfo.Fields, &document.Field{
			Name:      "AdvertiserId",
			IndexType: 2,
			Value:     int64(*info.AdvertiserId),
			ValueType: document.IntFieldType,
		})
	}
	if info.Platform != nil {
		docInfo.Fields = append(docInfo.Fields, &document.Field{
			Name:      "Platform",
			IndexType: 2,
			Value:     int64(*info.Platform),
			ValueType: document.IntFieldType,
		})
	}
	if info.Price != nil {
		docInfo.Fields = append(docInfo.Fields, &document.Field{
			Name:      "Price",
			IndexType: 2,
			Value:     *info.Price,
			ValueType: document.FloatFieldType,
		})
	}
	if info.StartTime != nil {
		docInfo.Fields = append(docInfo.Fields, &document.Field{
			Name:      "StartTime",
			IndexType: 1,
			Value:     int64(*info.StartTime),
			ValueType: document.IntFieldType,
		})
	}

	if info.EndTime != nil {
		docInfo.Fields = append(docInfo.Fields, &document.Field{
			Name:      "EndTime",
			IndexType: 1,
			Value:     int64(*info.EndTime),
			ValueType: document.IntFieldType,
		})
	}

	if info.CampaignType != nil {
		docInfo.Fields = append(docInfo.Fields, &document.Field{
			Name:      "CampaignType",
			IndexType: 1,
			Value:     int64(*info.CampaignType),
			ValueType: document.IntFieldType,
		})
	}

	if info.OsVersionMaxV2 != nil {
		docInfo.Fields = append(docInfo.Fields, &document.Field{
			Name:      "OsVersionMaxV2",
			IndexType: 1,
			Value:     int64(*info.OsVersionMaxV2),
			ValueType: document.IntFieldType,
		})
	}

	if info.OsVersionMinV2 != nil {
		docInfo.Fields = append(docInfo.Fields, &document.Field{
			Name:      "OsVersionMinV2",
			IndexType: 1,
			Value:     int64(*info.OsVersionMinV2),
			ValueType: document.IntFieldType,
		})
	}

	docInfo.Fields = append(docInfo.Fields, &document.Field{
		Name:      "PackageName",
		IndexType: 1,
		Value:     info.PackageName,
		ValueType: document.StringFieldType,
	})

	docInfo.Fields = append(docInfo.Fields, &document.Field{
		Name:      "DeviceTypeV2",
		IndexType: 0,
		Value:     info.DeviceTypeV2,
		ValueType: document.SliceFieldType,
	})

	return docInfo
}

func (c *CampaignParser) Parse(bytes []byte, userData interface{}) *builder.ParserResult {
	ud, ok := userData.(*UserData)
	if !ok {
		return nil
	}
	campaign := &CampaignInfo{}
	if err := bson.Unmarshal(bytes, &campaign); err != nil {
		fmt.Println("bson.Unmarshal error:" + err.Error())
	}
	if ud.upTime < campaign.Uptime {
		ud.upTime = campaign.Uptime
	}
	var info = MakeInfo(campaign)
	var mode builder.DataMod = builder.DataDel
	if campaign.Status == 1 {
		mode = builder.DataAddOrUpdate
	}
	return &builder.ParserResult{
		DataMod: mode,
		Value:   info,
	}
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
	if e := b.Build(ctx, "indexName"); e != nil {
		fmt.Println("build error", e.Error())
	}

	tIndex := b.GetIndex()

	// search: advertiserId=457 or platform=android or (price in [20.0, 1.4, 3.6, 5.7, 2.5] And price >= 1.4)
	// invert list
	invertIdx := tIndex.GetInvertedIndex()

	// storage
	storageIdx := tIndex.GetStorageIndex()

	//q := query.NewOrQuery(
	//	[]query.Query{
	//		query.NewAndQuery(
	//			[]query.Query{
	//				query.NewTermQuery(storageIdx.Iterator("Price")),
	//				query.NewTermQuery(storageIdx.Iterator("AdvertiserId")),
	//			},
	//			[]check.Checker{
	//				check.NewInChecker(storageIdx.Iterator("Price"), 2.3, 1.4, 3.65, 2.46, 2.5),
	//				check.NewNotChecker(storageIdx.Iterator("AdvertiserId"), 647, 658, 670),
	//			},
	//		),
	//		query.NewTermQuery(invertIdx.Iterator("AdvertiserId", "457")),
	//		query.NewTermQuery(invertIdx.Iterator("Platform", "1")),
	//	}, nil,
	//)

	q := query.NewOrQuery([]query.Query{
		query.NewOrQuery([]query.Query{
			query.NewTermQuery(invertIdx.Iterator("Platform", "1")),
		}, nil),
		query.NewOrQuery([]query.Query{
			query.NewTermQuery(invertIdx.Iterator("AdvertiserId", "457")),
		}, nil),
		query.NewOrQuery([]query.Query{
			query.NewTermQuery(invertIdx.Iterator("DeviceTypeV2", "4")),
			query.NewTermQuery(invertIdx.Iterator("DeviceTypeV2", "5")),
		}, nil),
		query.NewAndQuery([]query.Query{
			query.NewAndQuery([]query.Query{
				query.NewTermQuery(storageIdx.Iterator("Price")),
			}, []check.Checker{
				check.NewInChecker(storageIdx.Iterator("Price"),
					2.3, 1.4, 3.65, 2.46, 2.5),
			}),
			query.NewAndQuery([]query.Query{
				query.NewTermQuery(storageIdx.Iterator("AdvertiserId")),
			}, []check.Checker{
				check.NewNotChecker(storageIdx.Iterator("AdvertiserId"), int64(647), int64(658), int64(670)),
			})}, nil)},
		nil,
	)

	r1 := search.NewSearcher()
	r1.Search(tIndex, q)
	fmt.Println("+****************************+")
	fmt.Println("res: ", len(r1.Docs), r1.Time)
	//fmt.Println("+****************************+")
	//fmt.Println(r1.QueryDebug)
	//fmt.Println("+****************************+")
	fmt.Println(r1.IndexDebug)
	//fmt.Println("+****************************+")

	tIndex.UnsetDebug()

	a := "AdvertiserId=457 or Platform=1 or (Price in [2.3, 1.4, 3.65, 2.46, 2.5] and AdvertiserId !in [647, 658, 670])"
	sq := query.NewSqlQuery(a)

	m := sq.LRD(tIndex)
	r2 := search.NewSearcher()
	r2.Search(tIndex, m)
	//fmt.Println(r2.QueryDebug)
	fmt.Println(r2.IndexDebug)
	fmt.Println("+****************************+")
	fmt.Println("res: ", len(r2.Docs), r2.Time)

	fmt.Println(helpers.SliceEqual(r1.Docs, r2.Docs))

	c := make(chan os.Signal)
	signal.Notify(c)
	s := <-c
	fmt.Println("退出信号", s)
}
