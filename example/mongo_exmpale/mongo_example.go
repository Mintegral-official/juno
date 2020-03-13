package main

import (
	"context"
	"fmt"
	"github.com/MintegralTech/juno/builder"
	"github.com/MintegralTech/juno/check"
	"github.com/MintegralTech/juno/document"
	"github.com/MintegralTech/juno/query"
	"github.com/MintegralTech/juno/search"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"os"
	"os/signal"
	"strconv"
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
	UpTime int64
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
			IndexType: 1,
			Value:     *info.Price,
			ValueType: document.FloatFieldType,
		})
	}
	if info.StartTime != nil {
		docInfo.Fields = append(docInfo.Fields, &document.Field{
			Name:      "StartTime",
			IndexType: 2,
			Value:     int64(*info.StartTime),
			ValueType: document.IntFieldType,
		})
	}

	if info.EndTime != nil {
		docInfo.Fields = append(docInfo.Fields, &document.Field{
			Name:      "EndTime",
			IndexType: 2,
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
			IndexType: 2,
			Value:     int64(*info.OsVersionMaxV2),
			ValueType: document.IntFieldType,
		})
	}

	if info.OsVersionMinV2 != nil {
		docInfo.Fields = append(docInfo.Fields, &document.Field{
			Name:      "OsVersionMinV2",
			IndexType: 2,
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
		IndexType: 2,
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
	if ud.UpTime < campaign.Uptime {
		ud.UpTime = campaign.Uptime
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
	logger := logrus.New()
	b, e := builder.NewMongoIndexBuilder(&builder.MongoIndexManagerOps{
		URI:            "mongodb://13.250.108.190:27017",
		IncInterval:    5,
		BaseInterval:   120,
		IncParser:      &CampaignParser{},
		BaseParser:     &CampaignParser{},
		BaseQuery:      bson.M{"advertiserId": 903, "publisherId": 0, "status": 1, "system": 5},
		DB:             "new_adn",
		Collection:     "campaign",
		ConnectTimeout: 10000,
		ReadTimeout:    20000,
		UserData:       &UserData{},
		Logger:         logger,
		OnBeforeInc: func(userData interface{}) interface{} {
			ud, ok := userData.(*UserData)
			if !ok {
				return nil
			}
			incQuery := bson.M{"advertiserId": 903, "publisherId": 0, "status": 1, "system": 5, "updated": bson.M{"$gte": ud.UpTime - 5, "$lte": time.Now().Unix()}}
			return incQuery
		},
		OnFinishInc: func(userData interface{}) {
			ud, ok := userData.(*UserData)
			if !ok {
				return
			}
			fmt.Printf("OnFinishInc[%d], [%d]", time.Now().UnixNano(), ud.UpTime)
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

	var p = []float64{2.3, 1.4, 3.65, 2.46, 2.5}
	var a0 = []int64{647, 658, 670}
	//var dev = []int64{4, 5}

	for i := 0; i < 0; i++ {
		q := query.NewOrQuery([]query.Query{
			// ==
			query.NewOrQuery([]query.Query{
				query.NewTermQuery(invertIdx.Iterator("Platform", "1")),
			}, nil),
			// ==
			query.NewOrQuery([]query.Query{
				query.NewTermQuery(invertIdx.Iterator("AdvertiserId", "457")),
			}, nil),
			query.NewAndQuery([]query.Query{
				// in
				query.NewAndQuery([]query.Query{
					query.NewTermQuery(storageIdx.Iterator("Price")),
				}, []check.Checker{
					check.NewInChecker(storageIdx.Iterator("Price"), p, nil, false),
				}),
				// not in
				query.NewAndQuery([]query.Query{
					query.NewTermQuery(storageIdx.Iterator("AdvertiserId")),
				}, []check.Checker{
					check.NewNotChecker(storageIdx.Iterator("AdvertiserId"), a0, nil, false),
				})}, nil),
		}, nil)

		q.SetDebug(1)

		tquery := time.Now()
		r1 := search.NewSearcher()
		r1.Search(tIndex, q)
		fmt.Println("query: ", time.Since(tquery))
		fmt.Println("+****************************+")
		fmt.Println("res: ", len(r1.Docs), r1.Time)
		fmt.Println(r1.Docs[0])
		fmt.Println(invertIdx.GetValueById(document.DocId(1526540701)))
		fmt.Println(r1.QueryDebug)
	}

	c := make(chan os.Signal)
	signal.Notify(c)
	s := <-c
	fmt.Println("退出信号", s)

}

type myOperation struct {
	value interface{}
}

func (o *myOperation) Equal(value interface{}) bool {
	// your logic
	switch o.value.(type) {
	case map[string]int:
		return o.value.(map[string]int)[value.(string)] == 1
	}
	return true
}

func (o *myOperation) Less(value interface{}) bool {
	// your logic
	return true
}

func (o *myOperation) In(value interface{}) bool {
	// your logic
	switch value.(type) {
	// campaign.AdSchedule
	case map[string][]int:
		v, ok := o.value.(string)
		if !ok {
			return false
		}
		if _, ok = value.(map[string][]int)[v]; ok {
			return true
		}
		//condition.AdvertiserBlocklist  AdvertiserWhitelist  BlockIndustryIds
	case map[string]bool:
		//if len(value.(map[string]bool)) <= 0 {
		//	return false
		//}
		v, ok := o.value.(int64)
		if !ok {
			return false
		}
		if _, ok = value.(map[string]bool)[strconv.FormatInt(v, 10)]; ok {
			return true
		}
		// campaign.AuditAdvertiserMap
	case map[string]string:
		v, ok := o.value.(string)
		if !ok {
			return false
		}
		if _, ok := value.(map[string]string)[v]; ok {
			return true
		}
		// campaign.SubCategoryName  condition.BSubCategoryName
	case []string:
		scn, ok := value.([]string)
		if !ok {
			return false
		}
		bscn, ok := o.value.([]string)
		if !ok {
			return false
		}
		if len(bscn) > len(scn) {
			return false
		}
		for _, v := range bscn {
			for i := 0; i < len(scn); i++ {
				if v == scn[i] {
					continue
				} else if i == len(scn) {
					return false
				}
			}
		}
		return true
	}
	return false
}

func (o *myOperation) SetValue(value interface{}) {
	o.value = value
}
