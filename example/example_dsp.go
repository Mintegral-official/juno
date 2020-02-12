package main

import (
	"context"
	"fmt"
	"github.com/Mintegral-official/juno/builder"
	"github.com/Mintegral-official/juno/index"
	"github.com/Mintegral-official/juno/query"
	"github.com/Mintegral-official/juno/query/check"
	"github.com/Mintegral-official/juno/query/operation"
	"github.com/Mintegral-official/juno/search"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"strconv"
	"strings"
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

	}
	return true
}

func (o *myOperations) SetValue(value interface{}) {
	o.value = value
}

/**
campaignId : both
campaignType: both
os: both
countryCode, CityCode: invert
Direct: invert
IsEcAdv  RetargetVisitorType: invert
industryId: both
CType: invert
domain:both
DeviceAndIpuaRetarget:invert
AppCategory: invert
AppSubCategory:invert
ContentRating:storage
subCategoryV2:invert
InstallApps, ExcludeInstalledApps:both
trafficType:invert
AppSite:invert
adtype:invert
InventoryBlackList:invert
effectiveCountryCode:both
supportHttps:invert
Gender NeedGender:invert
devicelanguage: invert
adSchedule: invert
iabCategory: both
networkType: invert
deviceModel: invert
OsVersionMax, OsVersionMin: storage
endTime startTime: storage
advertiserId: both
*/

//  campaignId
//  len(condition.WhiteOfferList) > 0 campaign.campaignId in condition.WhiteOfferList
//  len(condition.BlackOfferList) > 0 campaign.campaignId not in condition.BlackOfferList
func campaignIdQuery(idx *index.Indexer, cond *CampaignCondition) query.Query {
	invertIdx := idx.GetInvertedIndex()
	var oq, naq []query.Query
	if len(cond.WhiteOfferList) > 0 {
		for k, v := range cond.WhiteOfferList {
			if !v {
				continue
			}
			oq = append(oq, query.NewTermQuery(invertIdx.Iterator("CampaignId", k)))
		}
	}
	naq = append(naq, query.NewOrQuery(oq, nil))
	if len(condition.BlackOfferList) > 0 {
		for k, v := range cond.BlackOfferList {
			if !v {
				continue
			}
			naq = append(naq, query.NewTermQuery(invertIdx.Iterator("CampaignId", k)))
		}
	}
	return query.NewNotAndQuery(naq, nil)
}

// campaignType
// condition.ForbidApk == true campaign.CampaignType != 3  (mvutil.CAMPAIGN_TYPE_APK)
// mvutil.AppDownload in condition.BlockAuditIndustryIds campaign.campaignType != 2(mvutil.CAMPAIGN_TYPE_GOOGLEPLAY)and campaign.campaignType != 3(mvutil.CAMPAIGN_TYPE_APK)
func campaignTypeQuery(idx *index.Indexer, cond *CampaignCondition) query.Query {
	invertIdx := idx.GetInvertedIndex()
	var q []query.Query
	if cond.ForbidApk == true {
		q = append(q, query.NewTermQuery(invertIdx.Iterator("CampaignType", strconv.Itoa(mvutil.CAMPAIGN_TYPE_APK))))
	}
	if _, ok := condition.BlockAuditIndustryIds[mvutil.AppDownload]; ok {
		q = append(q, query.NewAndQuery([]query.Query{
			query.NewTermQuery(invertIdx.Iterator("CampaignType", strconv.Itoa(mvutil.CAMPAIGN_TYPE_GOOGLEPLAY))),
			query.NewTermQuery(invertIdx.Iterator("CampaignType", strconv.Itoa(mvutil.CAMPAIGN_TYPE_APK))),
		}, nil))
	}
	return query.NewAndQuery(q, nil)
}

// Os
// campaign.os == condition.Os
func osQuery(idx *index.Indexer, cond *CampaignCondition) query.Query {
	invertIdx := idx.GetInvertedIndex()
	return query.NewTermQuery(invertIdx.Iterator("Os", cond.Os))
}

// country city
// campaign.country == condition.Country  and campaign.citycode == conditon.CityCode) or country == "ALL"
func countryCodeQuery(idx *index.Indexer, cond *CampaignCondition) query.Query {
	invertIdx := idx.GetInvertedIndex()
	query.NewOrQuery([]query.Query{
		query.NewAndQuery([]query.Query{
			query.NewTermQuery(invertIdx.Iterator("CountryCode", cond.Country)),
			query.NewTermQuery(invertIdx.Iterator("CityCode", cond.CityCode)),
		}, nil),
		query.NewTermQuery(invertIdx.Iterator("CountryCode", "ALL")),
	}, nil)
}

// CType
//condition.DevImpLimit  == true campaign.ctype != 1(cpa) and campaign.ctype != 5(cpe)}
//condition.NeedSearchCpc  == false campaign.ctype != 2(cpc) and campaign.ctype != 3(cpm)
//condition.TrafficType == "site" campaign.ctype != 1(cpa) and campaign.ctype != 5(cpe)
func ctypeQuery(idx *index.Indexer, cond *CampaignCondition) query.Query {
	invertIdx := idx.GetInvertedIndex()
	var q1, q2, q3 query.Query
	if condition.DevImpLimit == true {
		q1 = query.NewAndQuery([]query.Query{
			query.NewTermQuery(invertIdx.Iterator("CType", "1")),
			query.NewTermQuery(invertIdx.Iterator("CType", "5")),
		}, nil)
	}
	if condition.NeedSearchCpc == false {
		q2 = query.NewAndQuery([]query.Query{
			query.NewTermQuery(invertIdx.Iterator("CType", "2")),
			query.NewTermQuery(invertIdx.Iterator("CType", "3")),
		}, nil)
	}
	if condition.TrafficType == "site" {
		q3 = query.NewAndQuery([]query.Query{
			query.NewTermQuery(invertIdx.Iterator("CType", "1")),
			query.NewTermQuery(invertIdx.Iterator("CType", "5")),
		}, nil)
	}
	return query.NewAndQuery([]query.Query{q1, q2, q3}, nil)
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

	// storage list
	storageIdx := tIndex.GetStorageIndex()

	// Direct
	// condition.Adx == doubleclick campaign.Direct != 2
	var directQuery query.Query
	if cond.Adx == doubleclick {
		directQuery = query.NewTermQuery(invertIdx.Iterator("Direct", "2"))
	}

	// IsEcAdv  RetargetVisitorType
	// campaign.isEcAdv == 1 (mvutil.EcAdv） campaign.retargetVisitionType == 1 (include)
	var rtQuery query.Query
	rtQuery = query.NewAndQuery([]query.Query{
		query.NewTermQuery(invertIdx.Iterator("IsEcAdv", "1")),
		query.NewTermQuery(invertIdx.Iterator("RetargetVisitorType", "1")),
	}, nil)

	// industryId
	// len(condition.BlockIndustryIds) > 0 campaign.IndustryId not in condition.BlockIndustryIds
	var q []query.Query
	if len(condition.BlockIndustryIds) > 0 {
		for k, v := range condition.BlockIndustryIds {
			if !v {
				continue
			}
			q = append(q, query.NewTermQuery(invertIdx.Iterator("IndustryId", k)))
		}
	}
	var industryIdQuery = query.NewOrQuery(q, nil)
	// domain
	//len(condition.Badv) != 0  campaign.domain not in condition.Badv
	var q1 []query.Query
	if len(condition.Badv) != 0 {
		for k := range condition.Badv {
			q1 = append(q1, query.NewTermQuery(invertIdx.Iterator("Domain", k)))
		}
	}
	var domainQuery = query.NewOrQuery(q1, nil)

	// DeviceAndIpuaRetarget
	//if (cond.Os == android && country==china ){
	//if IsImei == false && isAndroidId == false{
	//campaign.DeviceAndIpuaRetarget != 1
	//
	//}else{
	//
	//campaign.DeviceAndIpuaRetarget != 2
	//
	//}
	//
	//}else{
	//
	//if cond.IsGoolleAdid == false{
	//
	//campaign.DeviceAndIpuaRetarget != 1
	//
	//}else{
	//
	//campaign.DeviceAndIpuaRetarget != 2
	//
	//}
	//}
	var deviceAndIpuaRetargetQuery query.Query
	if cond.Os == android && country == china {
		if IsImei == false && isAndroidId == false {
			deviceAndIpuaRetargetQuery = query.NewTermQuery(invertIdx.Iterator("DeviceAndIpuaRetarget", "1"))
		} else {
			deviceAndIpuaRetargetQuery = query.NewTermQuery(invertIdx.Iterator("DeviceAndIpuaRetarget", "2"))
		}
	} else {
		if cond.IsGoolleAdid == false {
			deviceAndIpuaRetargetQuery = query.NewTermQuery(invertIdx.Iterator("DeviceAndIpuaRetarget", "1"))
		} else {
			deviceAndIpuaRetargetQuery = query.NewTermQuery(invertIdx.Iterator("DeviceAndIpuaRetarget", "2"))
		}
	}

	// AppCategory
	// len(condition.BAppCategory)>0 campaign.AppCategory not in condition.BAppCategory
	var q2 []query.Query
	if len(condition.BAppCategory) > 0 {
		for k, v := range condition.BAppCategory {
			if !v {
				continue
			}
			q2 = append(q2, query.NewTermQuery(invertIdx.Iterator("AppCategory", k)))
		}
	}
	var appCategory = query.NewOrQuery(q2, nil)

	//AppSubCategory
	//len(condition.BAppSubCategory)>0 campaign.AppSubCategory not in condition.BAppSubCatetory
	var q3 []query.Query
	if len(condition.BAppSubCategory) > 0 {
		for k, v := range condition.BAppSubCategory {
			if !v {
				continue
			}
			q3 = append(q3, query.NewTermQuery(invertIdx.Iterator("AppSubCategory", k)))
		}
	}
	var appSubCategory = query.NewOrQuery(q3, nil)

	//SubCategoryName
	//condition.BSubCategoryName != nil condition.BSubCategoryName not in campaign.SubCategoryName 切片对切片
	var subCategoryName query.Query
	if condition.BSubCategoryName != nil {
		var q []query.Query
		for _, v := range condition.BSubCategoryName {
			q = append(q, query.NewTermQuery(invertIdx.Iterator("SubCategoryName", v)))
		}
		subCategoryName = query.NewAndQuery(q, nil)
	}
	//ContentRating
	//condition.Coppa==1 campaign.ContentRating <=12
	var contentRating query.Query
	if condition.Coppa == 1 {
		contentRating = query.NewAndQuery([]query.Query{
			query.NewTermQuery(storageIdx.Iterator("ContentRating")),
		}, []check.Checker{
			check.NewChecker(storageIdx.Iterator("ContentRating"), 12, operation.LE, nil, false)
		})
	}

	//subCategoryV2
	//len(condition.BSubCategorySDK) != 0 campaign.SubCategoryV2 not  in condition.BSUbCategorySDK
	var q4 []query.Query
	if len(condition.BSubCategorySDK) != 0 {
		for k := range condition.BSubCategorySDK {
			q4 = append(q4, query.NewTermQuery(invertIdx.Iterator("SubCategoryV2", strconv.Itoa(k))))
		}
	}
	var subCategoryV2 = query.NewOrQuery(q4, nil)

	// InstallApps
	// include	if len(campaign.InstallApps) >0 {campaign.InstallApps in condition.InstallApps}
	var q5 []query.Query
	for k := range condition.InstallApps {
		q5 = append(q5, query.NewTermQuery(invertIdx.Iterator("InstallApps", strconv.Itoa(k))))
	}
	var installApps = query.NewOrQuery(q5, nil)

	//ExcludeInstalledApps
	//exclude	if len(campaign.InstallApps) >0 {campaign.ExcludeInstalledApps not in condition.InstallApps}
	var q6 []query.Query
	for k := range condition.InstallApps {
		q6 = append(q6, query.NewTermQuery(invertIdx.Iterator("ExcludeInstalledApps", strconv.Itoa(k))))
	}
	var excludeInstalledApps = query.NewOrQuery(q6, nil)

	//trafficType
	//condition.TrafficType == mvutil.Site mvutil.SiteDirect in campaign.trafficType and if len(campaign.AppSite){  condition.AppSite in campaign.AppSite }
	var trafficType query.Query
	if condition.TrafficType == mvutil.Site {
		trafficType = query.NewAndQuery([]query.Query{
			query.NewTermQuery(invertIdx.Iterator("TrafficType", strconv.Itoa(int(mvutil.SiteDirect)))),
			query.NewTermQuery(invertIdx.Iterator("needAppSite", "1")),
			query.NewTermQuery(invertIdx.Iterator("AppSite", condition.AppSite)),
		}, nil)
	}

	// adtype
	//len(campaign.InventoryV2.AdnAdtype)>0 condition.AdType in campaign.InventoryV2.AdnAdtype
	var adType query.Query
	if len(campaign.InventoryV2.AdnAdtype) > 0 {
		adType = query.NewTermQuery(invertIdx.Iterator("AdType", condition.AdType))
	}
	//pkgName
	//len(campaign.inventoryBlackList) > 0 condition.PkgName not in campaign.inventoryBlackList
	var packageName query.Query
	if len(campaign.inventoryBlackList) > 0 {
		packageName = query.NewTermQuery(invertIdx.Iterator("InventoryBlackList", condition.PkgName))
	}

	//effectiveCountryCode
	//len(campaign.effectiveCountryCode) > 0 campaign.effectiveCountryCode[condition.Country] == 1 or campaign.effectiveCountryCode["ALL"] == 1
	var effective = query.NewOrQuery([]query.Query{
		query.NewTermQuery(invertIdx.Iterator("EffectiveCountryCode", condition.Country)),
		query.NewTermQuery(invertIdx.Iterator("EffectiveCountryCode", "ALL")),
	}, nil)

	//supportHttps
	var supportHttps query.Query
	if condition.Https == true {
		supportHttps = query.NewTermQuery(invertIdx.Iterator("SupportHttps", "1"))
	}

	// Gender NeedGender
	var gender = query.NewAndQuery([]query.Query{
		query.NewTermQuery(invertIdx.Iterator("NeedGender", "1")),
		query.NewTermQuery(invertIdx.Iterator("Gender", strconv.Itoa(condition.RequestGender))),
	}, nil)

	//devicelanguage
	//conditon.DeviceLanguage not in campaign.DeviceLanguage["2"] and
	//condition.DeviceLanguage in  campaign.DeviceLanguage[“1”]
	var devicelanguage = query.NewNotAndQuery([]query.Query{
		query.NewTermQuery(invertIdx.Iterator("DeviceLanguage_1", condition.DeviceLanguage)),
		query.NewTermQuery(invertIdx.Iterator("DeviceLanguage_2", condition.DeviceLanguage)),
	}, nil)

	// adSchedule
	//“-2“ in campaign.AdSchedule or ("-1-curHour" in campaign.AdSchedule ) or (curDay-curHour in campaign.AdSchedule)
	var adSchedule = query.NewOrQuery([]query.Query{
		query.NewTermQuery(invertIdx.Iterator("AdSchedule", "-2")),
		query.NewTermQuery(invertIdx.Iterator("AdSchedule", "-1-curHour")),
		query.NewTermQuery(invertIdx.Iterator("AdSchedule", curDay-curHour)),
	}, nil)

	//iabCategory
	//len(condition.Bcat) > 0
	//len(campaign.IabCategory)==0 or campaign.IabCategory key 和value（slice) not in condition.Bcat
	var q7 []query.Query
	if len(condition.Bcat) > 0 {
		for k := range condition.Bcat {
			q7 = append(q7, query.NewTermQuery(invertIdx.Iterator("IabCategory", k)))
		}
	}
	var iabCategory = query.NewOrQuery(q6, nil)

	//IabCategoryTag1 IabCategoryTag2
	var iabCategoryTag1, iabCategoryTag2 []query.Query
	if len(condition.IabCategory) != 0 {
		for _, v := range condition.IabCategory {
			tmp := strings.Split(v, "-")
			iabCategoryTag1 = append(iabCategoryTag1, query.NewTermQuery(invertIdx.Iterator("IabCategoryTag1", tmp[0])))
			iabCategoryTag2 = append(iabCategoryTag2, query.NewTermQuery(invertIdx.Iterator("IabCategoryTag1", v)))
		}
	}

	var iab = query.NewOrQuery([]query.Query{
		query.NewOrQuery(iabCategoryTag1, nil),
		query.NewAndQuery(iabCategoryTag2, nil),
	}, nil)

	//networkType len(campaign.networkType)>0 condition.NetworkType in campaign.NetworkType
	var networkType = query.NewTermQuery(invertIdx.Iterator("NetworkType", strconv.Itoa(int(condition.NetworkType))))

	//deviceModel
	//len(campaign.deviceModelV3)>0 condition.Make in campaign.deviceModelV3 or generateKey in deviceModelV3
	var deviceModel = query.NewOrQuery([]query.Query{
		query.NewTermQuery(invertIdx.Iterator("DeviceModel", condition.Make)),
		query.NewTermQuery(invertIdx.Iterator("DeviceModel", generateKey)),
	}, nil)

	// condition.osv >= campaign.osvMin and condition.osv <= osvMax
	var osv = query.NewAndQuery([]query.Query{
		query.NewTermQuery(storageIdx.Iterator("OsVersionMax")),
		query.NewTermQuery(storageIdx.Iterator("OsVersionMin")),
	}, []check.Checker{
		check.NewChecker(storageIdx.Iterator("OsVersionMin"), condition.osv, operation.LE, nil, false),
		check.NewChecker(storageIdx.Iterator("OsVersionMax"), condition.osv, operation.GE, nil, false),
	})

	// endTime startTime
	var endTime = query.NewAndQuery([]query.Query{
		query.NewTermQuery(storageIdx.Iterator("EndTime")),
	}, []check.Checker{
		check.NewChecker(storageIdx.Iterator("EndTime"), time.Now().Unix(), operation.GT, nil, false),
	})

	// advertiserId
	//len(condition.AdvertiserBlocklist) > 0 campaign.advertiserId not in condition.AdvertiserBlocklist
	//advertiserId	len(condition.AdvertiserWhitelist) > 0 campaign.advertiser in condition.AdvertiserWhitelist
	var blocklist, whitelist query.Query
	if len(condition.AdvertiserBlocklist) > 0 {
		var q []query.Query
		for k, v := range condition.AdvertiserBlocklist {
			if !v {
				continue
			}
			q = append(q, query.NewTermQuery(invertIdx.Iterator("AdvertiserId", k)))
		}
		blocklist = query.NewOrQuery(q, nil)
	}
	if len(condition.AdvertiserWhitelist) > 0 {
		var q []query.Query
		for k, v := range condition.AdvertiserWhitelist {
			if !v {
				continue
			}
			q = append(q, query.NewTermQuery(invertIdx.Iterator("AdvertiserId", k)))
		}
		whitelist = query.NewOrQuery(q, nil)
	}
	var advertiserId = query.NewNotAndQuery([]query.Query{whitelist, blocklist}, nil)

	//res := query.NewAndQuery([]query.Query{
	//	query.NewNotAndQuery([]query.Query{
	//		campaignIdQuery(tIndex, cond),
	//		campaignTypeQuery(tIndex, cond),
	//		ctypeQuery(tIndex, cond),
	//		directQuery,
	//		rtQuery,
	//		industryIdQuery,
	//		domainQuery,
	//		deviceAndIpuaRetargetQuery,
	//		appCategory,
	//		appSubCategory,
	//		subCategoryV2,
	//		excludeInstalledApps,
	//		packageName,
	//		iabCategory,
	//		subCategoryName,
	//	}, nil),
	//	osQuery(tIndex, cond),
	//	countryCodeQuery(tIndex, cond),
	//	iab,
	//	//contentRating,
	//	installApps,
	//	trafficType,
	//	adType,
	//	effective,
	//	supportHttps,
	//	gender,
	//	devicelanguage,
	//	adSchedule,
	//	networkType,
	//	deviceModel,
	//	//osv,
	//	advertiserId,
	//	//endTime,
	//}, nil)

	resQuery := query.NewAndQuery([]query.Query{
		query.NewNotAndQuery([]query.Query{
			campaignIdQuery(tIndex, cond),
			campaignTypeQuery(tIndex, cond),
			ctypeQuery(tIndex, cond),
			directQuery,
			rtQuery,
			industryIdQuery,
			domainQuery,
			deviceAndIpuaRetargetQuery,
			appCategory,
			appSubCategory,
			subCategoryV2,
			excludeInstalledApps,
			packageName,
			iabCategory,
			subCategoryName,
		}, nil),
		osQuery(tIndex, cond),
		countryCodeQuery(tIndex, cond),
		iab,
		installApps,
		trafficType,
		adType,
		effective,
		supportHttps,
		gender,
		devicelanguage,
		adSchedule,
		networkType,
		deviceModel,
		advertiserId,
		query.NewTermQuery(storageIdx.Iterator("EndTime")),
		query.NewTermQuery(storageIdx.Iterator("OsVersionMax")),
		query.NewTermQuery(storageIdx.Iterator("OsVersionMin")),
		query.NewTermQuery(storageIdx.Iterator("ContentRating")),
	}, []check.Checker{
		check.NewChecker(storageIdx.Iterator("OsVersionMin"), condition.osv, operation.LE, nil, false),
		check.NewChecker(storageIdx.Iterator("OsVersionMax"), condition.osv, operation.GE, nil, false),
		check.NewChecker(storageIdx.Iterator("EndTime"), time.Now().Unix(), operation.GT, nil, false),
		check.NewChecker(storageIdx.Iterator("ContentRating"), 12, operation.LE, nil, false),
	})

	searcher := search.NewSearcher()
	searcher.Search(tIndex, resQuery)
	fmt.Println(searcher.Docs)
	fmt.Println(searcher.Time)
}
