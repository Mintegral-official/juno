package main

import (
	"context"
	"fmt"
	"github.com/Mintegral-official/juno/builder"
	"github.com/Mintegral-official/juno/index"
	"github.com/Mintegral-official/juno/query"
	"github.com/Mintegral-official/juno/query/check"
	"github.com/Mintegral-official/juno/query/operation"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"strconv"
	"time"
)

type myOperations struct {
	value interface{}
}

func (o *myOperations) Equal(value interface{}) bool {
	// your logic
	return true
}

func (o *myOperations) Less(value interface{}) bool {
	// your logic
	return true
}

func (o *myOperations) In(value interface{}) bool {
	switch value.(type) {
	case map[string]bool:
		// industryId int64   condition.BlockIndustryIds map[string]bool
		if v, ok := o.value.(int64); ok {
			if _, ok := value.(map[string]bool)[strconv.FormatInt(v, 10)]; ok {
				return true
			}
			return false
		}
		// Domain string conditon.Badv map[string]bool
		// AppCategory string conditon.BAppCategory map[string]bool
		// AppSubCategory string  condition.BAppSubCategory map[string]bool
		if v, ok := o.value.(string); ok {
			if _, ok := value.(map[string]bool)[v]; ok {
				return true
			}
			return false
		}
	// SubCategoryName string
	case []string:
		v, ok := o.value.([]string)
		if !ok {
			return false
		}
		for i := 0; i < len(value.([]string)); i++ {
			for j := 0; j < len(v); j++ {
				if v[j] == value.([]string)[i] {
					continue
				} else if j == len(v)-1 {
					return false
				}
			}
			if i == len(value.([]string))-1 {
				return true
			}
		}
	}
	return true
}

func (o *myOperations) SetValue(value interface{}) {
	o.value = value
}

/**
campaignId : invert
campaignType: both
os: both
countryCode, CityCode: invert
Direct: both
industryId: both
domain:both

*/

func campaignIdQuery(idx *index.Indexer, cond *CampaignCondition) query.Query {
	invert := idx.GetInvertedIndex()
	var oq, naq []query.Query
	if len(cond.WhiteOfferList) > 0 {
		for wol := range cond.WhiteOfferList {
			oq = append(oq, query.NewTermQuery(invert.Iterator("CampaignId", wol)))
		}
	}
	naq = append(naq, query.NewOrQuery(oq, nil))
	if len(condition.BlackOfferList) > 0 {
		for wol := range cond.BlackOfferList {
			naq = append(naq, query.NewTermQuery(invert.Iterator("CampaignId", wol)))
		}
	}
	return query.NewNotAndQuery(naq, nil)
}

func campaignTypeQuery(idx *index.Indexer, cond *CampaignCondition) query.Query {
	storage := idx.GetStorageIndex()
	var q []query.Query
	if cond.ForbidApk == true {
		q = append(q, query.NewAndQuery([]query.Query{
			query.NewTermQuery(storage.Iterator("CampaignType")),
		}, []check.Checker{
			check.NewChecker(storage.Iterator("CampaignType"), mvutil.CAMPAIGN_TYPE_APK, operation.NE, nil, false),
		}))
	}
	if _, ok := condition.BlockAuditIndustryIds[mvutil.AppDownload]; ok {
		q = append(q, query.NewAndQuery([]query.Query{
			query.NewTermQuery(storage.Iterator("CampaignType")),
		}, []check.Checker{
			check.NewChecker(storage.Iterator("CampaignType"), mvutil.CAMPAIGN_TYPE_GOOGLEPLAY, operation.NE, nil, false),
			check.NewChecker(storage.Iterator("CampaignType"), mvutil.CAMPAIGN_TYPE_APK, operation.NE, nil, false),
		}))
	}
	return query.NewAndQuery(q, nil)
}

func queryDsp() {
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

	// invert list
	invertIdx := tIndex.GetInvertedIndex()

	// storage
	storageIdx := tIndex.GetStorageIndex()

	// os
	osQuery := query.NewTermQuery(invertIdx.Iterator("Os", condition.Os))
	// country city
	ccQuery := NewOrQuery([]Query{
		NewAndQuery([]Query{
			NewTermQuery(invertIdx.Iterator(campaign.country, condition.Country)),
			NewTermQuery(invertIdx.Iterator(campaign.citycode, conditon.CityCode)),
		}, nil),
		NewTermQuery(invertIdx.Iterator(country, "ALL")),
	}, nil)
	// Direct
	var directQuery query.Query
	if condition.Adx == doubleclick {
		directQuery = query.NewAndQuery([]query.Query{
			query.NewTermQuery(storageIdx.Iterator("Direct")),
		}, []check.Checker{
			check.NewChecker(storageIdx.Iterator("Direct"), 2, operation.NE, nil, false),
		})
	}
	// industryId
	var industryIdQuery []query.Query
	if len(condition.BlockIndustryIds) > 0 {
		for bi := range condition.BlockIndustryIds {
			industryIdQuery = append(industryIdQuery, query.NewTermQuery(invertIdx.Iterator("IndustryId", bi)))
		}
	}
	// domain
	var domainQuery query.Query
	if len(condition.Badv) != 0 {
		domainQuery = query.NewAndQuery([]query.Query{
			query.NewTermQuery(storageIdx.Iterator("Domain")),
		}, []check.Checker{
			check.NewNotChecker(storageIdx.Iterator("Domain"), conditon.Badv, &myOperations{}, false),
		})
	}
	// DeviceAndIpuaRetarget TODO
	// AppCategory
	var appQuery query.Query
	if len(condition.BAppCategory) > 0 {
		appQuery = query.NewAndQuery([]query.Query{
			query.NewTermQuery(storageIdx.Iterator("AppCategory")),
		}, []check.Checker{
			check.NewNotChecker(storageIdx.Iterator("AppCategory"), conditon.BAppCategory, &myOperations{}, false),
		})
	}
	// AppSubCategory
	var appSubCategoryQuery query.Query
	if len(condition.BAppSubCategory) > 0 {
		appSubCategoryQuery = query.NewAndQuery([]query.Query{
			query.NewTermQuery(storageIdx.Iterator("appSubCategoryQuery")),
		}, []check.Checker{
			check.NewNotChecker(storageIdx.Iterator("appSubCategoryQuery"), conditon.BAppSubCategory, &myOperations{}, false),
		})
	}
	// SubCategoryName
	var subCategoryNameQuery query.Query
	if condition.BSubCategoryName != nil {
		appSubCategoryQuery = query.NewAndQuery([]query.Query{
			query.NewTermQuery(storageIdx.Iterator("SubCategoryName")),
		}, []check.Checker{
			check.NewNotChecker(storageIdx.Iterator("SubCategoryName"), conditon.BSubCategoryName, &myOperations{}, true),
		})
	}
	// ContentRating
	var contentRatingQuery query.Query
	if condition.Coppa == 1 {
		contentRatingQuery = query.NewAndQuery([]query.Query{
			query.NewTermQuery(storageIdx.Iterator("ContentRating")),
		}, []check.Checker{
			check.NewChecker(storageIdx.Iterator("ContentRating"), 12, operation.LE, nil, false),
		})
	}
	// subCategoryV2 TODO
	// osv
	var osvQuery query.Query
	osvQuery = query.NewAndQuery([]query.Query{
		query.NewTermQuery(storageIdx.Iterator("OsVersionMin")),
		query.NewTermQuery(storageIdx.Iterator("OsVersionMax")),
	}, []check.Checker{
		check.NewChecker(storageIdx.Iterator("OsVersionMin"), condition.Osv, operation.LE, nil, false),
		check.NewChecker(storageIdx.Iterator("OsVersionMax"), condition.Osv, operation.GE, nil, false),
	})
	//

}
