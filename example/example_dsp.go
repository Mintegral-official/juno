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
	"regexp"
	"strconv"
	"strings"
	"time"
)

type operations struct {
	value interface{}
}

func (o *operations) Equal(value interface{}) bool {
	return true
}

func (o *operations) Less(value interface{}) bool {
	return true
}

func (o *operations) In(value interface{}) bool {
	switch value.(type) {
	// mobileCode
	case string:
		for _, v := range o.value.([]string) {
			if strings.Contains(value.(string), v) {
				return true
			}
		}
		return false
		// UserInterest
	case map[int]bool:
		ov, v := o.value.([][]int), value.(map[int]bool)
		nInterestNum := len(ov)
		if 0 == nInterestNum {
			return true
		}
		if 0 == len(v) {
			// mvutil.InterestOthers
			v[999999] = true
		}

		for _, Interests := range ov {
			flag := false
			for _, interest := range Interests {
				if _, ok := v[interest]; ok {
					flag = true
					break
				}
			}
			if false == flag {
				return false
			}
		}
		return true
	}
	return false
}

func (o *operations) SetValue(value interface{}) {
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

//  *************query code*****************
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
	return query.NewTermQuery(invertIdx.Iterator("Os", strconv.Itoa(int(cond.Os))))
}

// country city
// ondition.Country in campaign.CountryCode   and conditon.CityCode in campaign.citycode or "ALL" in campaign.CountryCode
func countryCodeQuery(idx *index.Indexer, cond *CampaignCondition) query.Query {
	invertIdx := idx.GetInvertedIndex()
	return query.NewOrQuery([]query.Query{
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

//// osv
//// condition.osv >= campaign.osvMin and condition.osv <= osvMax
//func osvQuery(idx *index.Indexer, cond *CampaignCondition) query.Query {
//	storageIdx := idx.GetStorageIndex()
//	return query.NewAndQuery([]query.Query{
//		query.NewTermQuery(storageIdx.Iterator("OsVersionMax")),
//		query.NewTermQuery(storageIdx.Iterator("OsVersionMin")),
//	}, []check.Checker{
//		check.NewChecker(storageIdx.Iterator("OsVersionMin"), cond.osv, operation.LE, nil, false),
//		check.NewChecker(storageIdx.Iterator("OsVersionMax"), cond.osv, operation.GE, nil, false),
//	})
//}

//devicelanguage
//conditon.DeviceLanguage not in campaign.DeviceLanguage["2"] and
//condition.DeviceLanguage in  campaign.DeviceLanguage[“1”]
func devicelanguageQuery(idx *index.Indexer, cond *CampaignCondition) query.Query {
	invertIdx := idx.GetInvertedIndex()
	return query.NewNotAndQuery([]query.Query{
		query.NewTermQuery(invertIdx.Iterator("DeviceLanguage_1", cond.DeviceLanguage)),
		query.NewTermQuery(invertIdx.Iterator("DeviceLanguage_2", cond.DeviceLanguage)),
	}, nil)
}

// advertiserId
//len(condition.AdvertiserBlocklist) > 0 campaign.advertiserId not in condition.AdvertiserBlocklist
//advertiserId	len(condition.AdvertiserWhitelist) > 0 campaign.advertiser in condition.AdvertiserWhitelist
func advertiserIdQuery(idx *index.Indexer, cond *CampaignCondition) query.Query {
	invertIdx := idx.GetInvertedIndex()
	var blocklist, whitelist query.Query
	if len(cond.AdvertiserBlocklist) > 0 {
		var q []query.Query
		for k, v := range cond.AdvertiserBlocklist {
			if !v {
				continue
			}
			q = append(q, query.NewTermQuery(invertIdx.Iterator("AdvertiserId", k)))
		}
		blocklist = query.NewOrQuery(q, nil)
	}
	if len(cond.AdvertiserWhitelist) > 0 {
		var q []query.Query
		for k, v := range cond.AdvertiserWhitelist {
			if !v {
				continue
			}
			q = append(q, query.NewTermQuery(invertIdx.Iterator("AdvertiserId", k)))
		}
		whitelist = query.NewOrQuery(q, nil)
	}
	return query.NewNotAndQuery([]query.Query{whitelist, blocklist}, nil)
}

// industryId
// len(condition.BlockIndustryIds) > 0 campaign.IndustryId not in condition.BlockIndustryIds
func industryIdQuery(idx *index.Indexer, cond *CampaignCondition) query.Query {
	invertIdx := idx.GetInvertedIndex()
	var q []query.Query
	if len(cond.BlockIndustryIds) > 0 {
		for k, v := range cond.BlockIndustryIds {
			if !v {
				continue
			}
			q = append(q, query.NewTermQuery(invertIdx.Iterator("IndustryId", k)))
		}
	}
	return query.NewOrQuery(q, nil)
}

//supportHttps
func supportHttpsQuery(idx *index.Indexer, cond *CampaignCondition) query.Query {
	invertIdx := idx.GetInvertedIndex()
	var supportHttps query.Query
	if cond.Https == true {
		supportHttps = query.NewTermQuery(invertIdx.Iterator("SupportHttps", "1"))
	}
	return supportHttps
}

// domain
//len(condition.Badv) != 0  campaign.domain not in condition.Badv
func domainQuery(idx *index.Indexer, cond *CampaignCondition) query.Query {
	invertIdx := idx.GetInvertedIndex()
	var q []query.Query
	if len(cond.Badv) != 0 {
		for k := range cond.Badv {
			q = append(q, query.NewTermQuery(invertIdx.Iterator("Domain", k)))
		}
	}
	return query.NewOrQuery(q, nil)
}

//trafficType
// mvutil.SiteDirect in campaign.trafficType and
//（(campaign.NeedAppSite ==1 and  condition.AppSite in campaign.AppSite ) or campaign.NeedAppSite==0})

// AppSite
// ((campaign.NeedTrafficType ==1 and mvutil.AppDirect in campaign.TrafficType )or campaign.NeedTrafficType ==0)
//and ((campaign.NeedAppSite == 1 and condition.AppSite in campaign.AppSite) or campaign.NeetTrafficType == 0)
func trafficTypeQuery(idx *index.Indexer, cond *CampaignCondition) query.Query {
	invertIdx := idx.GetInvertedIndex()
	var trafficType query.Query
	if cond.TrafficType == mvutil.Site {
		trafficType = query.NewAndQuery([]query.Query{
			query.NewTermQuery(invertIdx.Iterator("TrafficType", strconv.Itoa(int(mvutil.SiteDirect)))),
			query.NewOrQuery([]query.Query{
				query.NewAndQuery([]query.Query{
					query.NewTermQuery(invertIdx.Iterator("NeedAppSite", "1")),
					query.NewTermQuery(invertIdx.Iterator("AppSite", cond.AppSite)),
				}, nil),
				query.NewTermQuery(invertIdx.Iterator("NeedAppSite", "0")),
			}, nil),
		}, nil)
	}
	trafficType = query.NewAndQuery([]query.Query{
		query.NewOrQuery([]query.Query{
			query.NewAndQuery([]query.Query{
				query.NewTermQuery(invertIdx.Iterator("NeedAppSite", "1")),
				query.NewTermQuery(invertIdx.Iterator("TrafficType", strconv.Itoa(int(mvutil.AppDirect)))),
			}, nil),
			query.NewTermQuery(invertIdx.Iterator("NeedTrafficType", "0")),
		}, nil),
		query.NewOrQuery([]query.Query{
			query.NewAndQuery([]query.Query{
				query.NewTermQuery(invertIdx.Iterator("NeedAppSite", "1")),
				query.NewTermQuery(invertIdx.Iterator("AppSite", cond.AppSite)),
			}, nil),
			query.NewTermQuery(invertIdx.Iterator("NeedTrafficType", "0")),
		}, nil),
	}, nil)
	return trafficType
}

//iabCategory
//len(condition.Bcat) > 0
//(campaign.NeedIabCategory == 1 and campaign.IabCategory not in condition.Bcat ) or campaign.NeedIabCategory == 0
func iabCategoryQuery(idx *index.Indexer, cond *CampaignCondition) query.Query {
	invertIdx := idx.GetInvertedIndex()
	var q []query.Query
	if len(cond.Bcat) > 0 {
		for k, v := range cond.Bcat {
			if v {
				q = append(q, query.NewTermQuery(invertIdx.Iterator("IabCategory", k)))
			}
		}
	}
	return query.NewNotAndQuery([]query.Query{
		query.NewTermQuery(invertIdx.Iterator("NeedIabCategory", "1")),
		query.NewOrQuery(q, nil),
	}, nil)

}

// Direct
// condition.Adx == doubleclick campaign.Direct != 2
func directQuery(idx *index.Indexer, cond *CampaignCondition) query.Query {
	invertIdx := idx.GetInvertedIndex()
	var directQuery query.Query
	if cond.Adx == doubleclick {
		directQuery = query.NewTermQuery(invertIdx.Iterator("Direct", "2"))
	}
	return directQuery
}

// AppCategory
// len(condition.BAppCategory)>0 campaign.AppCategory not in condition.BAppCategory
func appCategoryQuery(idx *index.Indexer, cond *CampaignCondition) query.Query {
	invertIdx := idx.GetInvertedIndex()
	var q []query.Query
	if len(cond.BAppCategory) > 0 {
		for k, v := range cond.BAppCategory {
			if !v {
				continue
			}
			q = append(q, query.NewTermQuery(invertIdx.Iterator("AppCategory", k)))
		}
	}
	return query.NewOrQuery(q, nil)
}

//AppSubCategory
//len(condition.BAppSubCategory)>0 campaign.AppSubCategory not in condition.BAppSubCatetory
func appSubCategoryQuery(idx *index.Indexer, cond *CampaignCondition) query.Query {
	invertIdx := idx.GetInvertedIndex()
	var q []query.Query
	if len(cond.BAppSubCategory) > 0 {
		for k, v := range cond.BAppSubCategory {
			if !v {
				continue
			}
			q = append(q, query.NewTermQuery(invertIdx.Iterator("AppSubCategory", k)))
		}
	}
	return query.NewOrQuery(q, nil)
}

//SubCategoryName
//condition.BSubCategoryName != nil condition.BSubCategoryName not in campaign.SubCategoryName 切片对切片
func subCategoryNameQuery(idx *index.Indexer, cond *CampaignCondition) query.Query {
	invertIdx := idx.GetInvertedIndex()
	var subCategoryName query.Query
	if cond.BSubCategoryName != nil {
		var q []query.Query
		for _, v := range cond.BSubCategoryName {
			q = append(q, query.NewTermQuery(invertIdx.Iterator("SubCategoryName", v)))
		}
		subCategoryName = query.NewAndQuery(q, nil)
	}
	return subCategoryName
}

////ContentRating
////condition.Coppa==1 campaign.ContentRating <=12
//func contentRatingQuery(idx *index.Indexer, cond *CampaignCondition) query.Query {
//	storageIdx := idx.GetStorageIndex()
//	var contentRating query.Query
//	if cond.Coppa == 1 {
//		contentRating = query.NewAndQuery([]query.Query{
//			query.NewTermQuery(storageIdx.Iterator("ContentRating")),
//		}, []check.Checker{
//			check.NewChecker(storageIdx.Iterator("ContentRating"), 12, operation.LE, nil, false),
//		})
//	}
//	return contentRating
//}

//subCategoryV2
//len(condition.BSubCategorySDK) != 0 campaign.SubCategoryV2 not  in condition.BSUbCategorySDK
func subCategoryV2Query(idx *index.Indexer, cond *CampaignCondition) query.Query {
	invertIdx := idx.GetInvertedIndex()
	var q []query.Query
	if len(cond.BSubCategorySDK) != 0 {
		for k := range cond.BSubCategorySDK {
			q = append(q, query.NewTermQuery(invertIdx.Iterator("SubCategoryV2", strconv.Itoa(k))))
		}
	}
	return query.NewOrQuery(q, nil)
}

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
func deviceAndIpuaRetargetQuery(idx *index.Indexer, cond *CampaignCondition) query.Query {
	invertIdx := idx.GetInvertedIndex()
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
	return deviceAndIpuaRetargetQuery
}

// campaign.WhiteList != nil
// (campaign.NeedMvAppidWhiteList ==1 and conditon.MvAppid in campaign.WhiteList) or campaign.NeedMvAppidWhiteList == 0
// (campaign.NeedMvAppIdBlackList == 1 and conditon.MvAppid not in campaign.BlackList) or campaign.NeedMvAppIdBlackList == 0
func mvAppIdQuery(idx *index.Indexer, cond *CampaignCondition) query.Query {
	invertIdx := idx.GetInvertedIndex()
	return query.NewNotAndQuery([]query.Query{
		query.NewOrQuery([]query.Query{
			query.NewAndQuery([]query.Query{
				query.NewTermQuery(invertIdx.Iterator("NeedMvAppidWhiteList", "1")),
				query.NewTermQuery(invertIdx.Iterator("WhiteList", cond.MvAppid)),
			}, nil),
			query.NewTermQuery(invertIdx.Iterator("NeedMvAppidWhiteList", "0")),
		}, nil),
		query.NewOrQuery([]query.Query{
			query.NewAndQuery([]query.Query{
				query.NewTermQuery(invertIdx.Iterator("NeedMvAppIdBlackList", "1")),
				query.NewTermQuery(invertIdx.Iterator("BlackList", cond.MvAppid)),
			}, nil),
			query.NewTermQuery(invertIdx.Iterator("NeedMvAppIdBlackList", "0")),
		}, nil),
	}, nil)
}

// adtype
//len(campaign.InventoryV2.AdnAdtype)>0 condition.AdType in campaign.InventoryV2.AdnAdtype
// (campaign.NeedAdType == 1 and condition.AdType in campaign.InventoryV2.AdnAdtype) or campaign.NeedAdType == 0
func adtypeQuery(idx *index.Indexer, cond *CampaignCondition) query.Query {
	invertIdx := idx.GetInvertedIndex()
	return query.NewOrQuery([]query.Query{
		query.NewAndQuery([]query.Query{
			query.NewTermQuery(invertIdx.Iterator("NeedAdType", "1")),
			query.NewTermQuery(invertIdx.Iterator("AdType", cond.AdType)),
		}, nil),
		query.NewTermQuery(invertIdx.Iterator("NeedAdType", "0")),
	}, nil)
}

//networkType len(campaign.networkType)>0 condition.NetworkType in campaign.NetworkType
//(campaign.NeedNetWorkType ==1 and condition.NetworkType in campaign.NetworkType) or campaign.NeedNetWorkType ==0
func networkTypeQuery(idx *index.Indexer, cond *CampaignCondition) query.Query {
	invertIdx := idx.GetInvertedIndex()
	return query.NewOrQuery([]query.Query{
		query.NewAndQuery([]query.Query{
			query.NewTermQuery(invertIdx.Iterator("NeedNetWorkType", "1")),
			query.NewTermQuery(invertIdx.Iterator("NetworkType", strconv.Itoa(int(cond.NetworkType)))),
		}, nil),
		query.NewTermQuery(invertIdx.Iterator("NeedNetWorkType", "0")),
	}, nil)
}

//deviceModel
//len(campaign.deviceModelV3)>0 condition.Make in campaign.deviceModelV3 or generateKey in deviceModelV3
//(campaign.NeedDeviceModel == 1 and (condition.Make in campaign.deviceModelV3 or generateKey in deviceModelV3)) \
//or campaign.NeedDeviceModel == 0
func deviceModelQuery(idx *index.Indexer, cond *CampaignCondition) query.Query {
	invertIdx := idx.GetInvertedIndex()
	return query.NewOrQuery([]query.Query{
		query.NewAndQuery([]query.Query{
			query.NewTermQuery(invertIdx.Iterator("NeedDeviceModel", "1")),
			query.NewOrQuery([]query.Query{
				query.NewTermQuery(invertIdx.Iterator("DeviceModel", cond.Make)),
				query.NewTermQuery(invertIdx.Iterator("DeviceModel", generateKey)),
			}, nil),
		}, nil),
		query.NewTermQuery(invertIdx.Iterator("NeedDeviceModel", "0")),
	}, nil)
}

//effectiveCountryCode
//len(campaign.effectiveCountryCode) > 0 campaign.effectiveCountryCode[condition.Country] == 1 or campaign.effectiveCountryCode["ALL"] == 1
//
//(campaign.NeedEffectiveCountry ==1 and (condition.Country_1 in campaign.effectiveCountryCode or
//"ALL_1" in campaign.effectiveCountryCode)) or campaign.NeedEffectiveCountry
func effectiveCountryCodeQuery(idx *index.Indexer, cond *CampaignCondition) query.Query {
	invertIdx := idx.GetInvertedIndex()
	return query.NewOrQuery([]query.Query{
		query.NewAndQuery([]query.Query{
			query.NewTermQuery(invertIdx.Iterator("NeedEffectiveCountry", "1")),
			query.NewOrQuery([]query.Query{
				query.NewTermQuery(invertIdx.Iterator("EffectiveCountryCode", cond.Country_1)),
				query.NewTermQuery(invertIdx.Iterator("EffectiveCountryCode", "ALL_1")),
			}, nil),
		}, nil),
		query.NewTermQuery(invertIdx.Iterator("NeedEffectiveCountry", "0")),
	}, nil)
}

// Gender NeedGender
// (campaign.NeedGender ==1 and condition.RequestGender in campaign.Gender) or campaign.NeedGender == 0
func genderQuery(idx *index.Indexer, cond *CampaignCondition) query.Query {
	invertIdx := idx.GetInvertedIndex()
	return query.NewOrQuery([]query.Query{
		query.NewAndQuery([]query.Query{
			query.NewTermQuery(invertIdx.Iterator("NeedGender", "1")),
			query.NewTermQuery(invertIdx.Iterator("Gender", strconv.Itoa(condition.RequestGender))),
		}, nil),
		query.NewTermQuery(invertIdx.Iterator("NeedGender", "0")),
	}, nil)

}

//pkgName
// ( campaign.NeedInventoryBlackList ==1 and condition.PkgName not in campaign.inventoryBlackList ) or
//campaign.NeedInventoryBlackList == 0
func pkgNameQuery(idx *index.Indexer, cond *CampaignCondition) query.Query {
	invertIdx := idx.GetInvertedIndex()
	return query.NewOrQuery([]query.Query{
		query.NewAndQuery([]query.Query{
			query.NewTermQuery(invertIdx.Iterator("NeedInventoryBlackList", "1")),
			query.NewTermQuery(invertIdx.Iterator("InventoryBlackList", cond.PkgName)),
		}, nil),
		query.NewTermQuery(invertIdx.Iterator("NeedInventoryBlackList", "0")),
	}, nil)
}

// IsEcAdv  RetargetVisitorType
// campaign.isEcAdv == 1 (mvutil.EcAdv） campaign.retargetVisitionType == 1 (include)
func isecadvQuery(idx *index.Indexer, cond *CampaignCondition) query.Query {
	invertIdx := idx.GetInvertedIndex()
	return query.NewAndQuery([]query.Query{
		query.NewTermQuery(invertIdx.Iterator("IsEcAdv", "1")),
		query.NewTermQuery(invertIdx.Iterator("RetargetVisitorType", "1")),
	}, nil)
}

// iab
// (campaign.NeedIabCategoryTag == 1 and condition.IabCategory 某个元素使用”-“split的第一个元素在campaign.IabCategoryTag1
//or 某个元素在campaign.IabCategryTag2) or campaign.NeedIabCategoryTag == 0
func iabQuery(idx *index.Indexer, cond *CampaignCondition) query.Query {
	invertIdx := idx.GetInvertedIndex()
	//IabCategoryTag1 IabCategoryTag2
	var iabCategoryTag1, iabCategoryTag2 []query.Query
	if len(cond.IabCategory) != 0 {
		for _, v := range cond.IabCategory {
			tmp := strings.Split(v, "-")
			iabCategoryTag1 = append(iabCategoryTag1, query.NewTermQuery(invertIdx.Iterator("IabCategoryTag1", tmp[0])))
			iabCategoryTag2 = append(iabCategoryTag2, query.NewTermQuery(invertIdx.Iterator("IabCategoryTag1", v)))
		}
	}
	return query.NewOrQuery([]query.Query{
		query.NewAndQuery([]query.Query{
			query.NewTermQuery(invertIdx.Iterator("NeedIabCategoryTag", "1")),
			query.NewOrQuery(iabCategoryTag1, nil),
			query.NewOrQuery(iabCategoryTag2, nil),
		}, nil),
		query.NewTermQuery(invertIdx.Iterator("NeedIabCategoryTag", "0")),
	}, nil)
}

// UserAge
// (campaign.NeedUserAge == 1 and condition.UserAge in campaign.UserAge) and campaign.NeedUserAge == 0
func userAgeQuery(idx *index.Indexer, cond *CampaignCondition) query.Query {
	invertIdx := idx.GetInvertedIndex()
	return query.NewOrQuery([]query.Query{
		query.NewAndQuery([]query.Query{
			query.NewTermQuery(invertIdx.Iterator("UserAge", strconv.Itoa(cond.UserAge))),
		}, nil),
		query.NewTermQuery(invertIdx.Iterator("NeedUserAge", "0")),
	}, nil)
}

// InstallApps
// include	if len(campaign.InstallApps) >0 {campaign.InstallApps in condition.InstallApps}
func installAppsQuery(idx *index.Indexer, cond *CampaignCondition) query.Query {
	invertIdx := idx.GetInvertedIndex()
	var q []query.Query
	for k, v := range cond.InstallApps {
		if !v {
			continue
		}
		q = append(q, query.NewTermQuery(invertIdx.Iterator("InstallApps", strconv.Itoa(k))))
	}
	return query.NewOrQuery(q, nil)
}

// ExcludeInstalledApps
//exclude	if len(campaign.InstallApps) >0 {campaign.ExcludeInstalledApps not in condition.InstallApps}
func excludeInstalledAppsQuery(idx *index.Indexer, cond *CampaignCondition) query.Query {
	invertIdx := idx.GetInvertedIndex()
	var q []query.Query
	for k := range cond.InstallApps {
		q = append(q, query.NewTermQuery(invertIdx.Iterator("ExcludeInstalledApps", strconv.Itoa(k))))
	}
	return query.NewOrQuery(q, nil)
}

// adx
// condition.Adx in  campaign.AdxWhiteBlack_1 and conditon.Adx not in campaign.AdxWhiteBlack_2
func adxQuery(idx *index.Indexer, cond *CampaignCondition) query.Query {
	invertIdx := idx.GetInvertedIndex()
	if cond.Adx == tencent {
		return query.NewTermQuery(invertIdx.Iterator("AdxInclude", cond.Adx))
	}
	return query.NewNotAndQuery([]query.Query{
		query.NewTermQuery(invertIdx.Iterator("AdxWhiteBlack_1", cond.Adx)),
		query.NewTermQuery(invertIdx.Iterator("AdxWhiteBlack_2", cond.Adx)),
	}, nil)
}

// advertiserAudit
// condition.Adx in campaign.AuditAdvertiserMap
func advertiserAuditQuery(idx *index.Indexer, cond *CampaignCondition) query.Query {
	invertIdx := idx.GetInvertedIndex()
	var advertiserAudit query.Query
	if cond.NeedAdvAudit == true {
		advertiserAudit = query.NewTermQuery(invertIdx.Iterator("AuditAdvertiserMap", cond.Adx))
	}
	return advertiserAudit
}

// creativeAudit
// condition.Adx in campaign.AuditCreativeMap
func creativeAuditQuery(idx *index.Indexer, cond *CampaignCondition) query.Query {
	invertIdx := idx.GetInvertedIndex()
	var advertiserAudit query.Query
	if cond.NeedCreativeAudit == true {
		advertiserAudit = query.NewTermQuery(invertIdx.Iterator("AuditCreativeMap", cond.Adx))
	}
	return advertiserAudit
}

// deviceType
// len(campaign.DeviceTypeV2) == 0 or (4 in campaign.DeviceTypeV2 and 5 in campaign.DeviceTypeV2) or
//conditon.DeviceType in campaign.DeviceTypeV2
func deviceTypeQuery(idx *index.Indexer, cond *CampaignCondition) query.Query {
	invertIdx, storage := idx.GetInvertedIndex(), idx.GetStorageIndex()
	if condtion.DeviceType != 0 {
		return query.NewOrQuery([]query.Query{
			query.NewAndQuery([]query.Query{
				query.NewTermQuery(invertIdx.Iterator("NeedDeviceType", "1")),
				query.NewTermQuery(invertIdx.Iterator("DeviceType", "4")),
				query.NewTermQuery(invertIdx.Iterator("DeviceType", "5")),
			}, nil),
			query.NewTermQuery(invertIdx.Iterator("NeedDeviceType", "0")),
			query.NewTermQuery(invertIdx.Iterator("DeviceType", strconv.Itoa(int(cond.DeviceType)))),
		}, nil)
	}
	return query.NewTermQuery(storage.Iterator("DeviceType"))
}

// mobileCode
func mobileCodeQuery(idx *index.Indexer, cond *CampaignCondition) query.Query {
	invertIdx, storageIdx := idx.GetInvertedIndex(), idx.GetStorageIndex()
	var q []query.Query
	q = append(q, query.NewTermQuery(invertIdx.Iterator("MobileCode", "All")))
	if isNum, _ := regexp.MatchString("(^[0-9]+$)", cond.Carrier); !isNum {
		q = append(q, query.NewAndQuery([]query.Query{
			query.NewTermQuery(storageIdx.Iterator("MobileCode")),
		}, []check.Checker{
			check.NewInChecker(storageIdx.Iterator("MobileCode"), cond.Carrier, &operations{}, false),
		}))
	}
	return query.NewOrQuery(q, nil)
}

// PackageName
// campaign.PackageName not in condition.BApp
func packageNameQuery(idx *index.Indexer, cond *CampaignCondition) query.Query {
	invertIdx := idx.GetInvertedIndex()
	var q []query.Query
	for k := range cond.BApp {
		q = append(q, query.NewTermQuery(invertIdx.Iterator("PackageName", k)))
	}
	return query.NewOrQuery(q, nil)
}

// UserInterest
// len(campaign.UserInterestV2) == 0 or campaign.UserInterestV2的每个二级数组中至少有一个元素在condition.DmpInterests里
func userInterestQuery(idx *index.Indexer, cond *CampaignCondition) query.Query {
	storageIdx := idx.GetStorageIndex()
	return query.NewAndQuery([]query.Query{
		query.NewTermQuery(storageIdx.Iterator("UserInterest")),
	}, []check.Checker{
		check.NewInChecker(storageIdx.Iterator("UserInterest"), cond.DmpInterests, &operations{}, false),
	})
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

	// adSchedule
	//“-2“ in campaign.AdSchedule or ("-1-curHour" in campaign.AdSchedule ) or (curDay-curHour in campaign.AdSchedule)
	var adSchedule = query.NewOrQuery([]query.Query{
		query.NewTermQuery(invertIdx.Iterator("AdSchedule", "-2")),
		query.NewTermQuery(invertIdx.Iterator("AdSchedule", "-1-curHour")),
		query.NewTermQuery(invertIdx.Iterator("AdSchedule", curDay-curHour)),
	}, nil)

	// endTime startTime TODO 这个在下面合并的时候拆开了写的
	var endTime = query.NewOrQuery([]query.Query{
		query.NewAndQuery([]query.Query{
			query.NewTermQuery(storageIdx.Iterator("EndTime")),
		}, []check.Checker{
			check.NewChecker(storageIdx.Iterator("EndTime"), time.Now().Unix(), operation.GT, nil, false),
			check.NewChecker(storageIdx.Iterator("EndTime"), 0, operation.GT, nil, false),
		}),
		query.NewTermQuery(storageIdx.Iterator("EndTime")),
	}, []check.Checker{
		check.NewChecker(storageIdx.Iterator("EndTime"), 0, operation.LE, nil, false),
	})

	resQuery := query.NewAndQuery([]query.Query{
		query.NewNotAndQuery([]query.Query{
			campaignIdQuery(tIndex, cond),
			campaignTypeQuery(tIndex, cond),
			ctypeQuery(tIndex, cond),
			directQuery(tIndex, cond),
			isecadvQuery(tIndex, cond),
			industryIdQuery(tIndex, cond),
			domainQuery(tIndex, cond),
			deviceAndIpuaRetargetQuery(tIndex, cond),
			appCategoryQuery(tIndex, cond),
			appSubCategoryQuery(tIndex, cond),
			subCategoryV2Query(tIndex, cond),
			excludeInstalledAppsQuery(tIndex, cond),
			packageNameQuery(tIndex, cond),
			iabCategoryQuery(tIndex, cond),
			subCategoryNameQuery(tIndex, cond),
		}, nil),
		osQuery(tIndex, cond),
		countryCodeQuery(tIndex, cond),
		iabQuery(tIndex, cond),
		installAppsQuery(tIndex, cond),
		trafficTypeQuery(tIndex, cond),
		adtypeQuery(tIndex, cond),
		effectiveCountryCodeQuery(tIndex, cond),
		supportHttpsQuery(tIndex, cond),
		genderQuery(tIndex, cond),
		devicelanguageQuery(tIndex, cond),
		adSchedule,
		networkTypeQuery(tIndex, cond),
		deviceModelQuery(tIndex, cond),
		advertiserIdQuery(tIndex, cond),
		userInterestQuery(tIndex, cond),
		adxQuery(tIndex, cond),
		userInterestQuery(tIndex, cond),
		advertiserAuditQuery(tIndex, cond),
		creativeAuditQuery(tIndex, cond),
		deviceTypeQuery(tIndex, cond),
		mobileCodeQuery(tIndex, cond),
		userAgeQuery(tIndex, cond),
		pkgNameQuery(tIndex, cond),
		mvAppIdQuery(tIndex, cond),
		query.NewTermQuery(storageIdx.Iterator("EndTime")),
		query.NewTermQuery(storageIdx.Iterator("OsVersionMax")),
		query.NewTermQuery(storageIdx.Iterator("OsVersionMin")),
		query.NewTermQuery(storageIdx.Iterator("ContentRating")),
		query.NewTermQuery(storageIdx.Iterator("UserInterest")),
	}, []check.Checker{
		check.NewChecker(storageIdx.Iterator("OsVersionMin"), condition.osv, operation.LE, nil, false),
		check.NewChecker(storageIdx.Iterator("OsVersionMax"), condition.osv, operation.GE, nil, false),
		check.NewChecker(storageIdx.Iterator("EndTime"), time.Now().Unix(), operation.GT, nil, false),
		check.NewChecker(storageIdx.Iterator("ContentRating"), 12, operation.LE, nil, false),
		check.NewInChecker(storageIdx.Iterator("UserInterest"), cond.DmpInterests, &operations{}, false),
	})

	searcher := search.NewSearcher()
	searcher.Search(tIndex, resQuery)
	fmt.Println(searcher.Docs)
	fmt.Println(searcher.Time)
}
