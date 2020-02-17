package main

import (
	"context"
	"fmt"
	"github.com/Mintegral-official/juno/builder"
	"github.com/Mintegral-official/juno/document"
	"github.com/Mintegral-official/juno/query"
	"github.com/Mintegral-official/juno/query/check"
	"github.com/Mintegral-official/juno/query/operation"
	"github.com/Mintegral-official/juno/search"
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

	var p = []float64{2.3, 1.4, 3.65, 2.46, 2.5}
	var a0 = []int64{647, 658, 670}
	var dev = []int64{4, 5}

	for i := 0; i < 10; i++ {
		q := query.NewOrQuery([]query.Query{
			// ==
			query.NewOrQuery([]query.Query{
				query.NewTermQuery(invertIdx.Iterator("Platform", "1")),
			}, nil),
			// ==
			query.NewOrQuery([]query.Query{
				query.NewTermQuery(invertIdx.Iterator("AdvertiserId", "457")),
			}, nil),
			/* special example */
			// [campaign] <-> [condition]
			// or in , one in
			query.NewOrQuery([]query.Query{
				query.NewTermQuery(storageIdx.Iterator("DeviceTypeV2")),
			}, []check.Checker{
				check.NewInChecker(storageIdx.Iterator("DeviceTypeV2"), dev, &myOperation{}, false),
			}),
			// or not
			query.NewOrQuery([]query.Query{
				query.NewTermQuery(storageIdx.Iterator("DeviceTypeV2")),
			}, []check.Checker{
				check.NewNotChecker(storageIdx.Iterator("DeviceTypeV2"), dev, &myOperation{}, false),
			}),
			// and
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
			// !=
			query.NewNotAndQuery([]query.Query{
				query.NewTermQuery(storageIdx.Iterator("AdvertiserId")),
				query.NewTermQuery(invertIdx.Iterator("AdvertiserId", "457")),
			}, nil),
			// !=
			query.NewAndQuery([]query.Query{
				query.NewTermQuery(storageIdx.Iterator("AdvertiserId")),
			}, []check.Checker{
				check.NewChecker(storageIdx.Iterator("AdvertiserId"), 457, operation.NE, nil, false),
			}),
		},
			nil,
		)

		tquery := time.Now()
		r1 := search.NewSearcher()
		r1.Search(tIndex, q)
		fmt.Println("query: ", time.Since(tquery))
		fmt.Println("+****************************+")
		fmt.Println("res: ", len(r1.Docs), r1.Time)
		//fmt.Println("+****************************+")
		//fmt.Println(r1.QueryDebug)
		//fmt.Println("+****************************+")
		//fmt.Println(r1.IndexDebug)
		//fmt.Println("+****************************+")

		tIndex.UnsetDebug()

		a := "AdvertiserId=457 or Platform=1 or (Price in [2.3, 1.4, 3.65, 2.46, 2.5] and AdvertiserId !in [647, 658, 670])"

		tsql := time.Now()
		sq := query.NewSqlQuery(a, nil, false)
		m := sq.LRD(tIndex)
		fmt.Println("sql parse: ", time.Since(tsql))
		r2 := search.NewSearcher()
		r2.Search(tIndex, m)
		fmt.Println("sql: ", time.Since(tsql))

		//fmt.Println(r2.QueryDebug)
		//fmt.Println(r2.IndexDebug)
		fmt.Println("+****************************+")
		fmt.Println("res: ", len(r2.Docs), r2.Time)

		fmt.Println(SliceEqual(r1.Docs, r2.Docs))
	}

	c := make(chan os.Signal)
	signal.Notify(c)
	s := <-c
	fmt.Println("退出信号", s)

	//// 1. os: campaign.os == condition.Os
	//query.NewTermQuery(invertIdx.Iterator(campaign.Os, condition.Os))
	////	2. country, city:  campaign.country == condition.Country  and campaign.CityCode == condition.CityCode) or country == "ALL"
	//query.NewOrQuery([]query.Query{
	//	query.NewAndQuery([]query.Query{
	//		query.NewTermQuery(invertIdx.Iterator(campaign.CountryCode, condition.Country)),
	//		query.NewTermQuery(invertIdx.Iterator(campaign.CityCode, condition.CityCode)),
	//	}, nil),
	//	query.NewTermQuery(invertIdx.Iterator(campaign.countryCode, "ALL")),
	//}, nil)
	//// 3. osv: condition.osv >= campaign.osvMin and condition.osv <= osvMax
	//query.NewAndQuery([]query.Query{
	//	query.NewTermQuery(storageIdx.Iterator(campaign.OsVersionMin)),
	//	query.NewTermQuery(storageIdx.Iterator(campaign.OsVersionMax)),
	//}, []check.Checker{
	//	check.NewChecker(storageIdx.Iterator(campaign.OsVersionMin), condition.Osv, operation.LE, nil, false),
	//	check.NewChecker(storageIdx.Iterator(campaign.OsVersionMax), condition.Osv, operation.GE, nil, false),
	//})
	////4. adx: condition.Adx in  campaign.AdxWhiteBlack[“1”]  and conditon.Adx not in campaign.AdxWhiteBlack["2"]
	//query.NewAndQuery([]query.Query{
	//	query.NewTermQuery(storageIdx.Iterator(campaign.AdxWhiteBlack + "_1")),
	//	query.NewTermQuery(storageIdx.Iterator(campaign.AdxWhiteBlack + "_2")),
	//}, []check.Checker{
	//	check.NewInChecker(storageIdx.Iterator(campaign.AdxWhiteBlack+"_1"), condition.Adx, nil, true),
	//	check.NewNotChecker(storageIdx.Iterator(campaign.AdxWhiteBlack+"_2"), condition.Adx, nil, true),
	//})
	////5. devicelanguage: conditon.DeviceLanguage not in campaign.DeviceLanguage["2"] and condition.DeviceLanguage in  campaign.DeviceLanguage[“1”]
	//query.NewAndQuery([]query.Query{
	//	query.NewTermQuery(storageIdx.Iterator(campaign.DeviceLanguage + "_1")),
	//	query.NewTermQuery(storageIdx.Iterator(campaign.DeviceLanguage + "_2")),
	//}, []check.Checker{
	//	check.NewInChecker(storageIdx.Iterator(campaign.DeviceLanguage+"_1"), condition.DeviceLanguage, nil, true),
	//	check.NewNotChecker(storageIdx.Iterator(campaign.DeviceLanguage+"_2"), condition.DeviceLanguage, nil, true),
	//})
	////6. adSchedule: “-2“ in campaign.AdSchedule or ("-1-curHour" in campaign.AdSchedule ) or (curDay-curHour in campaign.AdSchedule)
	//query.NewOrQuery([]query.Query{
	//	query.NewTermQuery(storageIdx.Iterator(campaign.AdSchedule)),
	//}, []check.Checker{
	//	check.NewInChecker(storageIdx.Iterator(campaign.AdSchedule), "-2", &myOperation{}, true),
	//	check.NewInChecker(storageIdx.Iterator(campaign.AdSchedule), "-1-curHour", &myOperation{}, true),
	//	check.NewInChecker(storageIdx.Iterator(campaign.AdSchedule), curDay-curHour, &myOperation{}, true),
	//})
	////7. advertiserId: len(condition.AdvertiserBlocklist) > 0   campaign.advertiserId not in condition.AdvertiserBlocklist
	////					len(condition.AdvertiserWhitelist) > 0  campaign.advertiser in condition.AdvertiserWhitelist
	//query.NewAndQuery([]query.Query{
	//	query.NewTermQuery(storageIdx.Iterator(campaign.advertiserId)),
	//}, []check.Checker{
	//	check.NewNotChecker(storageIdx.Iterator(campaign.advertiserId), condition.AdvertiserBlocklist, &myOperation{}, false),
	//	check.NewInChecker(storageIdx.Iterator(campaign.advertiserId), condition.AdvertiserWhitelist, &myOperation{}, false),
	//})
	////8. industryId:  len(condition.BlockIndustryIds) > 0  campaign.IndustryId not in condition.BlockIndustryIds
	//query.NewAndQuery([]query.Query{
	//	query.NewTermQuery(storageIdx.Iterator(campaign.IndustryId)),
	//}, []check.Checker{
	//	check.NewNotChecker(storageIdx.Iterator(campaign.IndustryId), condition.BlockIndustryIds, &myOperation{}, false),
	//})
	////9. advertiserAudit: condition.NeedAdvAudit == true condition.Adx in campaign.AuditAdvertiserMap
	////   creaticeAudit:  condition.NeedCreativeAudit == true condition.Adx in campaign.AuditCreativeMap
	//query.NewAndQuery([]query.Query{
	//	query.NewTermQuery(storageIdx.Iterator(campaign.AuditAdvertiserMap)),
	//}, []check.Checker{
	//	check.NewInChecker(storageIdx.Iterator(campaign.AuditAdvertiserMap), condition.Adx, &myOperation{}, true),
	//})
	////10. ctype:  condition.DevImpLimit  == true      campaign.ctype != 1(cpa) and campaign.ctype != 5(cpe)}
	////            condition.NeedSearchCpc  == false  campaign.ctype != 2(cpc) and campaign.ctype != 3(cpm)
	////			  condtion.TrafficType == "site"     campaign.ctype != 1(cpa) and campaign.ctype != 5(cpe)
	//query.NewAndQuery([]query.Query{
	//	query.NewTermQuery(storageIdx.Iterator(campaign.ctype)),
	//}, []check.Checker{
	//	check.NewChecker(storageIdx.Iterator(campaign.ctype), 1, operation.NE, nil, false),
	//	check.NewChecker(storageIdx.Iterator(campaign.ctype), 5, operation.NE, nil, false),
	//})
	//
	//query.NewAndQuery([]query.Query{
	//	query.NewTermQuery(storageIdx.Iterator(campaign.ctype)),
	//}, []check.Checker{
	//	check.NewChecker(storageIdx.Iterator(campaign.ctype), 2, operation.NE, nil, false),
	//	check.NewChecker(storageIdx.Iterator(campaign.ctype), 3, operation.NE, nil, false),
	//})
	////11. supportHttps:  condition.Https == true	campaign.supportHttps == true
	//query.NewTermQuery(invertIdx.Iterator(campaign.supportHttps, "true"))
	////12. domain: len(condition.Badv) != 0   campaign.domain not in conditon.Badv
	//query.NewAndQuery([]query.Query{
	//	query.NewTermQuery(storageIdx.Iterator(campaign.domain)),
	//}, []check.Checker{
	//	check.NewNotChecker(storageIdx.Iterator(campaign.domain), conditon.Badv, &myOperation{}, false),
	//})
	////*****13. deviceType: condtion.DeviceType != 0 len(campaign.DeviceTypeV2) == 0 or
	////(4 in campaign.DeviceTypeV2 and 5 in campaign.DeviceTypeV2) or conditon.DeviceType in campaign.DeviceTypeV2
	//query.NewOrQuery([]query.Query{
	//	query.NewTermQuery(invertIdx.Iterator(campaign.DeviceTypeV2), conditon.DeviceType),
	//	query.NewAndQuery([]query.Query{
	//		query.NewTermQuery(invertIdx.Iterator(campaign.DeviceTypeV2, 4)),
	//		query.NewTermQuery(invertIdx.Iterator(campaign.DeviceTypeV2, 5)),
	//	}, nil),
	//}, nil)
	////14. trafficType: condition.TrafficType == mvutil.Site     mvutil.SiteDirect in campaign.trafficType and
	////if len(campaign.AppSite){  condition.AppSite in campaign.AppSite }  TODO
	////15. AppSite TODO
	////16. iabCategory: TODO
	////17. condition.Adx == doubleclick campaign.Direct != 2
	//query.NewAndQuery([]query.Query{
	//	query.NewTermQuery(storageIdx.Iterator(campaign.Direct)),
	//}, []check.Checker{
	//	check.NewChecker(storageIdx.Iterator(campaign.Direct), 2, operation.NE, nil, false),
	//})
	////18. pkgName: len(condition.BApp)>0   campaign.pkgName not in condition.BApp
	//query.NewAndQuery([]query.Query{
	//	query.NewTermQuery(storageIdx.Iterator(campaign.pkgName)),
	//}, []check.Checker{
	//	check.NewNotChecker(storageIdx.Iterator(campaign.pkgName), condition.BApp, &myOperation{}, false),
	//})
	////19. AppCategory: len(condition.BAppCategory)>0  campaign.AppCategory not in condition.BAppCatetory
	//query.NewAndQuery([]query.Query{
	//	query.NewTermQuery(storageIdx.Iterator(campaign.AppCategory)),
	//}, []check.Checker{
	//	check.NewNotChecker(storageIdx.Iterator(campaign.AppCategory), condition.BAppCatetory, &myOperation{}, false),
	//})
	////20. AppSubCategory： len(condition.BAppSubCategory)>0  campaign.AppSubCategory not in condition.BAppSubCatetory
	//query.NewAndQuery([]query.Query{
	//	query.NewTermQuery(storageIdx.Iterator(campaign.AppSubCategory)),
	//}, []check.Checker{
	//	check.NewNotChecker(storageIdx.Iterator(campaign.AppSubCategory), condition.BAppCatetory, &myOperation{}, false),
	//})
	////21. SubCategoryName: condition.BSubCategoryName != nil condition.BSubCategoryName not in campaign.SubCategoryName 切片对切片
	//query.NewAndQuery([]query.Query{
	//	query.NewTermQuery(storageIdx.Iterator(campaign.SubCategoryName)),
	//}, []check.Checker{
	//	check.NewNotChecker(storageIdx.Iterator(campaign.SubCategoryName), condition.BSubCategoryName, &myOperation{}, true),
	//})
	////22. ContentRating: condition.Coppa==1 campaign.ContentRateing <=12
	//query.NewAndQuery([]query.Query{
	//	query.NewTermQuery(storageIdx.Iterator(campaign.ContentRateing)),
	//}, []check.Checker{
	//	check.NewChecker(storageIdx.Iterator(campaign.ContentRateing), 12, operation.LE, nil, false),
	//})
	////23. campaignId len(condition.WhiteOfferList) > 0   campaign.campaignId in condition.WhiteOfferList
	////				 len(condition.BlackOfferList) > 0    campaign.campaignId not in condition.BlackOfferList
	//query.NewAndQuery([]query.Query{
	//	query.NewTermQuery(storageIdx.Iterator(campaign.campaignId)),
	//}, []check.Checker{
	//	check.NewInChecker(storageIdx.Iterator(campaign.campaignId), condition.WhiteOfferList, &myOperation{}, false),
	//})
	//query.NewAndQuery([]query.Query{
	//	query.NewTermQuery(storageIdx.Iterator(campaign.campaignId)),
	//}, []check.Checker{
	//	check.NewNotChecker(storageIdx.Iterator(campaign.campaignId), condition.BlackOfferList, &myOperation{}, false),
	//})
	////24.subCategoryV2: len(condition.BSubCategorySDK) != 0  campaign.SubCategoryV2 not  in condition.BSUbCategorySDK TODO
	////25. campaignType  vvcondition.ForbidApk == true   campaign.CampaignType != 3  (mvutil.CAMPAIGN_TYPE_APK)
	//// 	mvutil.AppDownload in condition.BlockAuditIndustryIds campaign.campaignType != 2(mvutil.CAMPAIGN_TYPE_GOOGLEPLAY)and campaign.campaignType != 3(mvutil.CAMPAIGN_TYPE_APK)
	//query.NewAndQuery([]query.Query{
	//	query.NewTermQuery(storageIdx.Iterator(campaign.campaignType)),
	//}, []check.Checker{
	//	check.NewChecker(storageIdx.Iterator(campaign.campaignType), 3, operation.NE, nil, false),
	//})
	//query.NewAndQuery([]query.Query{
	//	query.NewTermQuery(storageIdx.Iterator(campaign.campaignType)),
	//}, []check.Checker{
	//	check.NewChecker(storageIdx.Iterator(campaign.campaignType), 2, operation.NE, nil, false),
	//})
	////26. effectiveCountryCode: len(campaign.effectiveCountryCode) > 0
	////	campaign.effectiveCountryCode[condition.Country] == 1 or campaign.effectiveCountryCode["ALL"] == 1
	//query.NewOrQuery([]query.Query{
	//	query.NewTermQuery(storageIdx.Iterator(campaign.effectiveCountryCode)),
	//}, []check.Checker{
	//	check.NewChecker(storageIdx.Iterator(campaign.effectiveCountryCode), condition.Country, operation.EQ, &myOperation{}, false),
	//	check.NewChecker(storageIdx.Iterator(campaign.effectiveCountryCode), "ALL", operation.EQ, &myOperation{}, false),
	//})
	////27. adx: condition.adx == tencent	condition.adx in campaign.AdxInclude  TODO
	//// 28. DeviceAndIpuaRetarget:  if (cond.Os == android && country==china ){
	////								  if IsImei == false && isAndroidId == false{
	////									  campaign.DeviceAndIpuaRetarget != 1
	////								  }else{
	////									  campaign.DeviceAndIpuaRetarget != 2
	////								  }
	////							  }else{
	////								if cond.IsGoolleAdid == false{
	////									campaign.DeviceAndIpuaRetarget != 1
	////								}else{
	////									campaign.DeviceAndIpuaRetarget != 2
	////								}
	////							} TODO
	////29. campaignType  AppDownload in condition.BlockAuditIndustryIds	campaign.CampaignType =!= GooglePlay && campaign.CampaignType != apk
	//query.NewAndQuery([]query.Query{
	//	query.NewTermQuery(storageIdx.Iterator(campaign.CampaignType)),
	//}, []check.Checker{
	//	check.NewChecker(storageIdx.Iterator(campaign.CampaignType), GooglePlay, operation.NE, nil, false),
	//	check.NewChecker(storageIdx.Iterator(campaign.CampaignType), apk, operation.NE, nil, false),
	//})
	////30.mvAppId: campaign.WhiteList != nil condition.MvAppid in campaign.WhiteList
	////			  campaign.BlackList != nil condition.MvAppid not in campaign.BlackList
	//query.NewAndQuery([]query.Query{
	//	query.NewTermQuery(storageIdx.Iterator(campaign.WhiteList)),
	//}, []check.Checker{
	//	check.NewInChecker(storageIdx.Iterator(campaign.WhiteList), condition.MvAppid, nil, true),
	//})
	//query.NewAndQuery([]query.Query{
	//	query.NewTermQuery(storageIdx.Iterator(campaign.BlackList)),
	//}, []check.Checker{
	//	check.NewNotChecker(storageIdx.Iterator(campaign.BlackList), condition.MvAppid, nil, true),
	//})
	////31. adtype: len(campaign.InventoryV2.AdnAdtype)>0 condition.AdType in campaign.InventoryV2.AdnAdtype
	//query.NewAndQuery([]query.Query{
	//	query.NewTermQuery(storageIdx.Iterator(campaign.InventoryV2.AdnAdtype)),
	//}, []check.Checker{
	//	check.NewInChecker(storageIdx.Iterator(campaign.InventoryV2.AdnAdtype), condition.AdType, nil, true),
	//})
	////32. networkType: len(campaign.networkType)>0  condition.NetworkType in campaign.NetworkType
	//query.NewAndQuery([]query.Query{
	//	query.NewTermQuery(storageIdx.Iterator(campaign.NetworkType)),
	//}, []check.Checker{
	//	check.NewInChecker(storageIdx.Iterator(campaign.NetworkType), condition.NetworkType, nil, true),
	//})
	////33.deviceModel: len(campaign.deviceModelV3)>0 condition.Make in campaign.deviceModelV3 or generateKey in deviceModelV3
	//// len(campaign.UserInterestV2) == 0 or campaign.UserInterestV2的每个二级数组中至少有一个元素在condition.DmpInterests里 TODO
	////34.gender: campaign.NeedGender == true  condition.RequestGender in campaign.Gender
	//query.NewTermQuery(invertIdx.Iterator(campaign.Gender, condition.RequestGender))
	////35. pkgName: len(campaign.inventoryBlackList) > 0 condition.PkgName not in campaign.inventoryBlackList
	//query.NewAndQuery([]query.Query{
	//	query.NewTermQuery(storageIdx.Iterator(campaign.inventoryBlackList)),
	//}, []check.Checker{
	//	check.NewNotChecker(storageIdx.Iterator(campaign.inventoryBlackList), condition.PkgName, nil, true),
	//})
	////36.endTime: campaign.endTime > 0 campaign.endTime > time.Now().Unix()
	//query.NewAndQuery([]query.Query{
	//	query.NewTermQuery(storageIdx.Iterator(campaign.endTime)),
	//}, []check.Checker{
	//	check.NewChecker(storageIdx.Iterator(campaign.endTime), time.Now().Unix(), operation.GT, nil, false),
	//})
	////37.电商广告	campaign.isEcAdv == 1 (mvutil.EcAdv）  campaign.retargetVisitionType == 1 (include)
	//query.NewTermQuery(invertIdx.Iterator(campaign.retargetVisitionType, "1"))
	////38.InventoryV2.iabcategory	campaign.NeedIabCategory == true && len(condition.IabCategory)!=0	 TODO
	////condition.IabCategory 某个元素使用”-“split的第一个元素在campaign.IabCategoryTag1 or 某个元素在campaign.IabCategryTag2
	////39. mobileCode	if campaign.mobileCode not have ”all“	condition.Carrier in TODO
	////40. UserAge	if NeedUserAge	condition.UserAge in campaign.UserAge
	//query.NewAndQuery([]query.Query{
	//	query.NewTermQuery(storageIdx.Iterator(campaign.UserAge)),
	//}, []check.Checker{
	//	check.NewInChecker(storageIdx.Iterator(campaign.UserAge), condition.UserAge, nil, true),
	//})
	////41.InstallApps	include	if len(campaign.InstallApps) >0 {campaign.InstallApps in condtion.InstallApp}
	//query.NewAndQuery([]query.Query{
	//	query.NewTermQuery(storageIdx.Iterator(campaign.InstallApps)),
	//}, []check.Checker{
	//	check.NewInChecker(storageIdx.Iterator(campaign.InstallApps), condition.InstallApp, nil, false),
	//})
	////42.ExcludeInstalledApps	exclude	if len(campaign.InstallApps) >0 {campaign.ExcludeInstalledApps not in condition.InstallApps}
	//query.NewAndQuery([]query.Query{
	//	query.NewTermQuery(storageIdx.Iterator(campaign.ExcludeInstalledApps)),
	//}, []check.Checker{
	//	check.NewNotChecker(storageIdx.Iterator(campaign.ExcludeInstalledApps), condition.InstallApps, nil, false),
	//})
	//43.campaignType	if campaign.CampaignType in mvutil.MapCampaignClickThroughType value not in condition.BlockCLickType(单子点击跳转类型） TODO
	//adx的第三方黑名单 TODO 词表相关
	//adx in FilterAdxThirdPartBlacklistTable
	//campaign.ThirdParty not in value
	//
	//
	//
	//
	//adx的第三方白名单
	//adx in FilterAdxThirdPartWhitelistTable
	//campaign.ThirdParty in value
	//流量包名白名单过滤
	//condition.Bundle in BundleCampaignWhitelistTable
	//
	//campaign.campaignId in value
	//
	//adx的包名黑名单过滤
	//condition.Adx in FilterAdxPackageTable
	//
	//campaign.pkgName not in value
	//adx的adomain黑名单过滤	condition.Adx in FilterAdxPackageTable	campaign.domain not in value
	//DevIdTable	campaign.retargetingDevice/DeviceId/Package/CampaignId
	//BlackCampaignTable	adx维度单子黑名单	campaignId not in value

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

func SliceEqual(a, b []document.DocId) bool {
	if len(a) != len(b) {
		return false
	}

	if (a == nil) != (b == nil) {
		return false
	}

	for i, v := range a {
		if v != b[i] {
			return false
		}
	}

	return true
}
