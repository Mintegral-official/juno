package main

//
//import (
//	"github.com/Mintegral-official/juno/check"
//	"github.com/Mintegral-official/juno/index"
//	"github.com/Mintegral-official/juno/operation"
//	"github.com/Mintegral-official/juno/query"
//	"regexp"
//	"strconv"
//	"strings"
//	"time"
//)
//
//type operations struct {
//	value interface{}
//}
//
//func (o *operations) Equal(value interface{}) bool {
//	switch value.(type) {
//	// mobileCode
//	case string:
//		ov := o.value.(map[string]string)
//		return ov[value.(string)] == "ALL"
//	}
//	return false
//}
//
//func (o *operations) Less(value interface{}) bool {
//	return true
//}
//
//func (o *operations) In(value interface{}) bool {
//	switch value.(type) {
//	// mobileCode
//	case string:
//		ov, v := o.value.(map[string]string), value.(string)
//		if isNum, _ := regexp.MatchString("(^[0-9]+$)", v); isNum {
//			if len(ov[v]) == 0 {
//				return false
//			}
//		} else {
//			for _, value := range ov {
//				if strings.Contains(v, value) {
//					return true
//				}
//			}
//		}
//		return false
//		// UserInterest
//	case map[int]bool:
//		ov, v := o.value.([][]int), value.(map[int]bool)
//		nInterestNum := len(ov)
//		if 0 == nInterestNum {
//			return true
//		}
//		if 0 == len(v) {
//			// mvutil.InterestOthers
//			v[999999] = true
//		}
//
//		for _, Interests := range ov {
//			flag := false
//			for _, interest := range Interests {
//				if _, ok := v[interest]; ok {
//					flag = true
//					break
//				}
//			}
//			if false == flag {
//				return false
//			}
//		}
//		return false
//	}
//	return false
//}
//
//func (o *operations) SetValue(value interface{}) {
//	o.value = value
//}
//
///**
//campaignId : both
//campaignType: both
//os: both
//countryCode, CityCode: invert
//Direct: invert
//IsEcAdv  RetargetVisitorType: invert
//industryId: both
//CType: invert
//domain:both
//DeviceAndIpuaRetarget:invert
//AppCategory: invert
//AppSubCategory:invert
//ContentRating:storage
//subCategoryV2:invert
//InstallApps, ExcludeInstalledApps:both
//trafficType:invert
//AppSite:invert
//adtype:invert
//InventoryBlackList:invert
//effectiveCountryCode:both
//supportHttps:invert
//Gender NeedGender:invert
//devicelanguage: invert
//adSchedule: invert
//iabCategory: both
//networkType: invert
//deviceModel: invert
//OsVersionMax, OsVersionMin: storage
//endTime startTime: storage
//advertiserId: both
//*/
//
////  *************query code*****************
////  campaignId
////  len(condition.WhiteOfferList) > 0 campaign.campaignId in condition.WhiteOfferList
////  len(condition.BlackOfferList) > 0 campaign.campaignId not in condition.BlackOfferList
//func campaignIdQuery(idx *index.Indexer, cond *CampaignCondition) query.Query {
//	invertIdx := idx.GetInvertedIndex()
//	var whiteOq, blackOq []query.Query
//	if len(cond.WhiteOfferList) > 0 {
//		for k, v := range cond.WhiteOfferList {
//			if !v {
//				continue
//			}
//			whiteOq = append(whiteOq, query.NewTermQuery(invertIdx.Iterator("CampaignId", strconv.FormatInt(k, 10))))
//		}
//	}
//	if len(whiteOq) != 0 {
//		return nil
//
//	}
//	if len(cond.BlackOfferList) > 0 {
//		for k, v := range cond.BlackOfferList {
//			if !v {
//				continue
//			}
//			blackOq = append(blackOq, query.NewTermQuery(invertIdx.Iterator("CampaignId", strconv.FormatInt(k, 10))))
//		}
//	}
//	if len(blackOq) == 0 {
//		return query.NewOrQuery(whiteOq, nil)
//	}
//	return query.NewNotAndQuery([]query.Query{
//		query.NewOrQuery(whiteOq, nil),
//		query.NewOrQuery(blackOq, nil),
//	}, nil)
//}

//
//// campaignType
//// condition.ForbidApk == true campaign.CampaignType != 3  (mvutil.CAMPAIGN_TYPE_APK)
//// mvutil.AppDownload in condition.BlockAuditIndustryIds campaign.campaignType != 2(mvutil.CAMPAIGN_TYPE_GOOGLEPLAY)and campaign.campaignType != 3(mvutil.CAMPAIGN_TYPE_APK)
//func campaignTypeQuery(idx *index.Indexer, cond *CampaignCondition) query.Query {
//	invertIdx := idx.GetInvertedIndex()
//	var q []query.Query
//	if cond.ForbidApk == true {
//		q = append(q, query.NewTermQuery(invertIdx.Iterator("CampaignType", strconv.Itoa(mvutil.CAMPAIGN_TYPE_APK))))
//	}
//	if _, ok := condition.BlockAuditIndustryIds[mvutil.AppDownload]; ok {
//		q = append(q, query.NewAndQuery([]query.Query{
//			query.NewTermQuery(invertIdx.Iterator("CampaignType", strconv.Itoa(mvutil.CAMPAIGN_TYPE_GOOGLEPLAY))),
//			query.NewTermQuery(invertIdx.Iterator("CampaignType", strconv.Itoa(mvutil.CAMPAIGN_TYPE_APK))),
//		}, nil))
//	}
//	if len(q) == 0 {
//		return nil
//	}
//	return query.NewAndQuery(q, nil)
//}
//
//// Os
//// campaign.os == condition.Os
//func osQuery(idx *index.Indexer, cond *CampaignCondition) query.Query {
//	invertIdx := idx.GetInvertedIndex()
//	return query.NewTermQuery(invertIdx.Iterator("Os", strconv.Itoa(int(cond.Os))))
//}
//
//// country city
//// ondition.Country in campaign.CountryCode   and conditon.CityCode in campaign.citycode or "ALL" in campaign.CountryCode
//func countryCodeQuery(idx *index.Indexer, cond *CampaignCondition) query.Query {
//	invertIdx := idx.GetInvertedIndex()
//	return query.NewOrQuery([]query.Query{
//		query.NewAndQuery([]query.Query{
//			query.NewTermQuery(invertIdx.Iterator("CountryCode", cond.Country)),
//			query.NewTermQuery(invertIdx.Iterator("CityCode", cond.CityCode)),
//		}, nil),
//		query.NewTermQuery(invertIdx.Iterator("CountryCode", "ALL")),
//	}, nil)
//}
//
//// CType
////condition.DevImpLimit  == true campaign.ctype != 1(cpa) and campaign.ctype != 5(cpe)}
////condition.NeedSearchCpc  == false campaign.ctype != 2(cpc) and campaign.ctype != 3(cpm)
////condition.TrafficType == "site" campaign.ctype != 1(cpa) and campaign.ctype != 5(cpe)
//func ctypeQuery(idx *index.Indexer, cond *CampaignCondition) query.Query {
//	invertIdx := idx.GetInvertedIndex()
//	var q1, q2, q3 query.Query = nil, nil, nil
//	if cond.DevImpLimit == true {
//		q1 = query.NewAndQuery([]query.Query{
//			query.NewTermQuery(invertIdx.Iterator("CType", "1")),
//			query.NewTermQuery(invertIdx.Iterator("CType", "5")),
//		}, nil)
//	}
//	if cond.NeedSearchCpc == false {
//		q2 = query.NewAndQuery([]query.Query{
//			query.NewTermQuery(invertIdx.Iterator("CType", "2")),
//			query.NewTermQuery(invertIdx.Iterator("CType", "3")),
//		}, nil)
//	}
//	if cond.TrafficType == "site" {
//		q3 = query.NewAndQuery([]query.Query{
//			query.NewTermQuery(invertIdx.Iterator("CType", "1")),
//			query.NewTermQuery(invertIdx.Iterator("CType", "5")),
//		}, nil)
//	}
//	var q []query.Query
//	if q1 != nil {
//		q = append(q, q1)
//	}
//	if q2 != nil {
//		q = append(q, q2)
//	}
//	if q3 != nil {
//		q = append(q, q3)
//	}
//	if len(q) == 0 {
//		return nil
//	}
//	return query.NewAndQuery(q, nil)
//}
//
////// osv
////// condition.osv >= campaign.osvMin and condition.osv <= osvMax
////func osvQuery(idx *index.Indexer, cond *CampaignCondition) query.Query {
////	storageIdx := idx.GetStorageIndex()
////	return query.NewAndQuery([]query.Query{
////		query.NewTermQuery(storageIdx.Iterator("OsVersionMax")),
////		query.NewTermQuery(storageIdx.Iterator("OsVersionMin")),
////	}, []check.Checker{
////		check.NewChecker(storageIdx.Iterator("OsVersionMin"), cond.osv, operation.LE, nil, false),
////		check.NewChecker(storageIdx.Iterator("OsVersionMax"), cond.osv, operation.GE, nil, false),
////	})
////}
//
////devicelanguage
//// ((campaign.NeedDeviceLanguageWhiteList = 1 and conditon.DeviceLanguage not in campaign.DeviceLanguageWhiteList) or
////campaign.NeedDeviceLanguageWhiteList = 0) and
////((NeedDeviceLanguageBlackList = 1 and condition.DeviceLanguage in  campaign.DeviceLanguageBlackList) or campaign.NeedDeviceLanguageBlackList =0)
//func devicelanguageQuery(idx *index.Indexer, cond *CampaignCondition) query.Query {
//	invertIdx := idx.GetInvertedIndex()
//	return query.NewAndQuery([]query.Query{
//		query.NewOrQuery([]query.Query{
//			query.NewNotAndQuery([]query.Query{
//				query.NewTermQuery(invertIdx.Iterator("NeedDeviceLanguageWhiteList", "1")),
//				query.NewTermQuery(invertIdx.Iterator("DeviceLanguageWhiteList", cond.DeviceLanguage)),
//			}, nil),
//			query.NewTermQuery(invertIdx.Iterator("NeedDeviceLanguageWhiteList", "0")),
//		}, nil),
//		query.NewOrQuery([]query.Query{
//			query.NewAndQuery([]query.Query{
//				query.NewTermQuery(invertIdx.Iterator("NeedDeviceLanguageWhiteList", "1")),
//				query.NewTermQuery(invertIdx.Iterator("DeviceLanguageWhiteList", cond.DeviceLanguage)),
//			}, nil),
//			query.NewTermQuery(invertIdx.Iterator("NeedDeviceLanguageWhiteList", "0")),
//		}, nil),
//	}, nil)
//}
//
//// advertiserId
////len(condition.AdvertiserBlocklist) > 0 campaign.advertiserId not in condition.AdvertiserBlocklist
////advertiserId	len(condition.AdvertiserWhitelist) > 0 campaign.advertiser in condition.AdvertiserWhitelist
//func advertiserIdQuery(idx *index.Indexer, cond *CampaignCondition) query.Query {
//	invertIdx := idx.GetInvertedIndex()
//	var blocklist, whitelist query.Query
//	if len(cond.AdvertiserBlocklist) > 0 {
//		var q []query.Query
//		for k, v := range cond.AdvertiserBlocklist {
//			q = append(q, query.NewTermQuery(invertIdx.Iterator("AdvertiserId", k)))
//		}
//		blocklist = query.NewOrQuery(q, nil)
//	} else {
//		blocklist = nil
//	}
//	if len(cond.AdvertiserWhitelist) > 0 {
//		var q []query.Query
//		for k, v := range cond.AdvertiserWhitelist {
//			q = append(q, query.NewTermQuery(invertIdx.Iterator("AdvertiserId", k)))
//		}
//		whitelist = query.NewOrQuery(q, nil)
//	} else {
//		whitelist = nil
//	}
//	if blocklist == nil {
//		return nil
//	}
//	return query.NewNotAndQuery([]query.Query{whitelist, blocklist}, nil)
//}
//
//// industryId
//// len(condition.BlockIndustryIds) > 0 campaign.IndustryId not in condition.BlockIndustryIds
//func industryIdQuery(idx *index.Indexer, cond *CampaignCondition) query.Query {
//	invertIdx := idx.GetInvertedIndex()
//	var q []query.Query
//	if len(cond.BlockIndustryIds) > 0 {
//		for k, v := range cond.BlockIndustryIds {
//			q = append(q, query.NewTermQuery(invertIdx.Iterator("IndustryId", k)))
//		}
//		return query.NewOrQuery(q, nil)
//	}
//	return nil
//}
//
////supportHttps
//func supportHttpsQuery(idx *index.Indexer, cond *CampaignCondition) query.Query {
//	invertIdx := idx.GetInvertedIndex()
//	if cond.Https == true {
//		return query.NewTermQuery(invertIdx.Iterator("SupportHttps", "1"))
//	}
//	return nil
//}
//
//// domain
////len(condition.Badv) != 0  campaign.domain not in condition.Badv
//func domainQuery(idx *index.Indexer, cond *CampaignCondition) query.Query {
//	invertIdx := idx.GetInvertedIndex()
//	var q []query.Query
//	if len(cond.Badv) != 0 {
//		for k := range cond.Badv {
//			q = append(q, query.NewTermQuery(invertIdx.Iterator("Domain", k)))
//		}
//		return query.NewOrQuery(q, nil)
//	}
//	return nil
//}
//
////trafficType
//// mvutil.SiteDirect in campaign.trafficType and
////（(campaign.NeedAppSite ==1 and  condition.AppSite in campaign.AppSite ) or campaign.NeedAppSite==0})
//
//// AppSite
//// ((campaign.NeedTrafficType ==1 and mvutil.AppDirect in campaign.TrafficType )or campaign.NeedTrafficType ==0)
////and ((campaign.NeedAppSite == 1 and condition.AppSite in campaign.AppSite) or campaign.NeetTrafficType == 0)
//func trafficTypeQuery(idx *index.Indexer, cond *CampaignCondition) query.Query {
//	invertIdx := idx.GetInvertedIndex()
//	var trafficType query.Query
//	if cond.TrafficType == mvutil.Site {
//		trafficType = query.NewAndQuery([]query.Query{
//			query.NewTermQuery(invertIdx.Iterator("TrafficType", strconv.Itoa(int(mvutil.SiteDirect)))),
//			query.NewOrQuery([]query.Query{
//				query.NewAndQuery([]query.Query{
//					query.NewTermQuery(invertIdx.Iterator("NeedAppSite", "1")),
//					query.NewTermQuery(invertIdx.Iterator("AppSite", cond.AppSite)),
//				}, nil),
//				query.NewTermQuery(invertIdx.Iterator("NeedAppSite", "0")),
//			}, nil),
//		}, nil)
//	} else {
//		trafficType = query.NewAndQuery([]query.Query{
//			query.NewOrQuery([]query.Query{
//				query.NewAndQuery([]query.Query{
//					query.NewTermQuery(invertIdx.Iterator("NeedAppSite", "1")),
//					query.NewTermQuery(invertIdx.Iterator("TrafficType", strconv.Itoa(int(mvutil.AppDirect)))),
//				}, nil),
//				query.NewTermQuery(invertIdx.Iterator("NeedTrafficType", "0")),
//			}, nil),
//			query.NewOrQuery([]query.Query{
//				query.NewAndQuery([]query.Query{
//					query.NewTermQuery(invertIdx.Iterator("NeedAppSite", "1")),
//					query.NewTermQuery(invertIdx.Iterator("AppSite", cond.AppSite)),
//				}, nil),
//				query.NewTermQuery(invertIdx.Iterator("NeedTrafficType", "0")),
//			}, nil),
//		}, nil)
//	}
//	return trafficType
//}
//
////iabCategory
////len(condition.Bcat) > 0
////(campaign.NeedIabCategory == 1 and campaign.IabCategory not in condition.Bcat ) or campaign.NeedIabCategory == 0
//func iabCategoryQuery(idx *index.Indexer, cond *CampaignCondition) query.Query {
//	invertIdx := idx.GetInvertedIndex()
//	var q []query.Query
//	if len(cond.Bcat) > 0 {
//		for k, v := range cond.Bcat {
//			if v {
//				q = append(q, query.NewTermQuery(invertIdx.Iterator("IabCategory", k)))
//			}
//		}
//		return query.NewOrQuery([]query.Query{
//			query.NewTermQuery(invertIdx.Iterator("NeedIabCategory", "0")),
//			query.NewNotAndQuery([]query.Query{
//				query.NewTermQuery(invertIdx.Iterator("NeedIabCategory", "1")),
//				query.NewOrQuery(q, nil),
//			}, nil),
//		}, nil)
//	}
//	return nil
//}
//
//// Direct
//// condition.Adx == doubleclick campaign.Direct != 2
//func directQuery(idx *index.Indexer, cond *CampaignCondition) query.Query {
//	invertIdx := idx.GetInvertedIndex()
//	if cond.Adx == doubleclick {
//		return query.NewTermQuery(invertIdx.Iterator("Direct", "2"))
//	}
//	return nil
//}
//
//// AppCategory
//// len(condition.BAppCategory)>0 campaign.AppCategory not in condition.BAppCategory
//func appCategoryQuery(idx *index.Indexer, cond *CampaignCondition) query.Query {
//	invertIdx := idx.GetInvertedIndex()
//	var q []query.Query
//	if len(cond.BAppCategory) > 0 {
//		for k, v := range cond.BAppCategory {
//			q = append(q, query.NewTermQuery(invertIdx.Iterator("AppCategory", k)))
//		}
//		return query.NewOrQuery(q, nil)
//	}
//	return nil
//}
//
////AppSubCategory
////len(condition.BAppSubCategory)>0 campaign.AppSubCategory not in condition.BAppSubCatetory
//func appSubCategoryQuery(idx *index.Indexer, cond *CampaignCondition) query.Query {
//	invertIdx := idx.GetInvertedIndex()
//	var q []query.Query
//	if len(cond.BAppSubCategory) > 0 {
//		for k, v := range cond.BAppSubCategory {
//			q = append(q, query.NewTermQuery(invertIdx.Iterator("AppSubCategory", k)))
//		}
//	}
//	if len(q) == 0 {
//		return nil
//	}
//	return query.NewOrQuery(q, nil)
//}
//
////SubCategoryName
////condition.BSubCategoryName != nil condition.BSubCategoryName not in campaign.SubCategoryName 切片对切片
//func subCategoryNameQuery(idx *index.Indexer, cond *CampaignCondition) query.Query {
//	invertIdx := idx.GetInvertedIndex()
//	if cond.BSubCategoryName != nil {
//		var q []query.Query
//		for _, v := range cond.BSubCategoryName {
//			q = append(q, query.NewTermQuery(invertIdx.Iterator("SubCategoryName", v)))
//		}
//		return query.NewAndQuery(q, nil)
//	}
//	return nil
//}
//
//////ContentRating
//////condition.Coppa==1 campaign.ContentRating <=12
////func contentRatingQuery(idx *index.Indexer, cond *CampaignCondition) query.Query {
////	storageIdx := idx.GetStorageIndex()
////	var contentRating query.Query
////	if cond.Coppa == 1 {
////		contentRating = query.NewAndQuery([]query.Query{
////			query.NewTermQuery(storageIdx.Iterator("ContentRating")),
////		}, []check.Checker{
////			check.NewChecker(storageIdx.Iterator("ContentRating"), 12, operation.LE, nil, false),
////		})
////	}
////	return contentRating
////}
//
////subCategoryV2
////len(condition.BSubCategorySDK) != 0 campaign.SubCategoryV2 not  in condition.BSUbCategorySDK
//func subCategoryV2Query(idx *index.Indexer, cond *CampaignCondition) query.Query {
//	invertIdx := idx.GetInvertedIndex()
//	var q []query.Query
//	if len(cond.BSubCategorySDK) != 0 {
//		for k := range cond.BSubCategorySDK {
//			q = append(q, query.NewTermQuery(invertIdx.Iterator("SubCategoryV2", strconv.Itoa(k))))
//		}
//	}
//	return query.NewOrQuery(q, nil)
//}
//
//// DeviceAndIpuaRetarget
////if (cond.Os == android && country==china ){
////if IsImei == false && isAndroidId == false{
////campaign.DeviceAndIpuaRetarget != 1
////
////}else{
////
////campaign.DeviceAndIpuaRetarget != 2
////
////}
////
////}else{
////
////if cond.IsGoolleAdid == false{
////
////campaign.DeviceAndIpuaRetarget != 1
////
////}else{
////
////campaign.DeviceAndIpuaRetarget != 2
////
////}
////}
//func deviceAndIpuaRetargetQuery(idx *index.Indexer, cond *CampaignCondition) query.Query {
//	invertIdx := idx.GetInvertedIndex()
//	var deviceAndIpuaRetargetQuery query.Query
//	if cond.Os == android && country == china {
//		if IsImei == false && isAndroidId == false {
//			deviceAndIpuaRetargetQuery = query.NewTermQuery(invertIdx.Iterator("DeviceAndIpuaRetarget", "1"))
//		} else {
//			deviceAndIpuaRetargetQuery = query.NewTermQuery(invertIdx.Iterator("DeviceAndIpuaRetarget", "2"))
//		}
//	} else {
//		if cond.IsGoolleAdid == false {
//			deviceAndIpuaRetargetQuery = query.NewTermQuery(invertIdx.Iterator("DeviceAndIpuaRetarget", "1"))
//		} else {
//			deviceAndIpuaRetargetQuery = query.NewTermQuery(invertIdx.Iterator("DeviceAndIpuaRetarget", "2"))
//		}
//	}
//	return deviceAndIpuaRetargetQuery
//}
//
//// campaign.WhiteList != nil
//// (campaign.NeedMvAppidWhiteList ==1 and conditon.MvAppid in campaign.WhiteList) or campaign.NeedMvAppidWhiteList == 0
//// (campaign.NeedMvAppIdBlackList == 1 and conditon.MvAppid not in campaign.BlackList) or campaign.NeedMvAppIdBlackList == 0
//func mvAppIdQuery(idx *index.Indexer, cond *CampaignCondition) query.Query {
//	invertIdx := idx.GetInvertedIndex()
//	return query.NewNotAndQuery([]query.Query{
//		query.NewOrQuery([]query.Query{
//			query.NewAndQuery([]query.Query{
//				query.NewTermQuery(invertIdx.Iterator("NeedMvAppidWhiteList", "1")),
//				query.NewTermQuery(invertIdx.Iterator("WhiteList", cond.MvAppid)),
//			}, nil),
//			query.NewTermQuery(invertIdx.Iterator("NeedMvAppidWhiteList", "0")),
//		}, nil),
//		query.NewOrQuery([]query.Query{
//			query.NewAndQuery([]query.Query{
//				query.NewTermQuery(invertIdx.Iterator("NeedMvAppIdBlackList", "1")),
//				query.NewTermQuery(invertIdx.Iterator("BlackList", cond.MvAppid)),
//			}, nil),
//			query.NewTermQuery(invertIdx.Iterator("NeedMvAppIdBlackList", "0")),
//		}, nil),
//	}, nil)
//}
//
//// adtype
////len(campaign.InventoryV2.AdnAdtype)>0 condition.AdType in campaign.InventoryV2.AdnAdtype
//// (campaign.NeedAdType == 1 and condition.AdType in campaign.InventoryV2.AdnAdtype) or campaign.NeedAdType == 0
//func adtypeQuery(idx *index.Indexer, cond *CampaignCondition) query.Query {
//	invertIdx := idx.GetInvertedIndex()
//	return query.NewOrQuery([]query.Query{
//		query.NewAndQuery([]query.Query{
//			query.NewTermQuery(invertIdx.Iterator("NeedAdType", "1")),
//			query.NewTermQuery(invertIdx.Iterator("AdType", cond.AdType)),
//		}, nil),
//		query.NewTermQuery(invertIdx.Iterator("NeedAdType", "0")),
//	}, nil)
//}
//
////networkType len(campaign.networkType)>0 condition.NetworkType in campaign.NetworkType
////(campaign.NeedNetWorkType ==1 and condition.NetworkType in campaign.NetworkType) or campaign.NeedNetWorkType ==0
//func networkTypeQuery(idx *index.Indexer, cond *CampaignCondition) query.Query {
//	invertIdx := idx.GetInvertedIndex()
//	return query.NewOrQuery([]query.Query{
//		query.NewAndQuery([]query.Query{
//			query.NewTermQuery(invertIdx.Iterator("NeedNetWorkType", "1")),
//			query.NewTermQuery(invertIdx.Iterator("NetworkType", strconv.Itoa(int(cond.NetworkType)))),
//		}, nil),
//		query.NewTermQuery(invertIdx.Iterator("NeedNetWorkType", "0")),
//	}, nil)
//}
//
////deviceModel
////len(campaign.deviceModelV3)>0 condition.Make in campaign.deviceModelV3 or generateKey in deviceModelV3
////(campaign.NeedDeviceModel == 1 and (condition.Make in campaign.deviceModelV3 or generateKey in deviceModelV3)) \
////or campaign.NeedDeviceModel == 0
////
////(campaign.NeedDeviceModel == 1 and (condition.Make in campaign.deviceModelV3 or generateKey in deviceModelV3)) or
////campaign.NeedDeviceModel == 0
//func deviceModelQuery(idx *index.Indexer, cond *CampaignCondition) query.Query {
//	invertIdx := idx.GetInvertedIndex()
//	var q query.Query
//	if len(cond.Hwv) > 0 && cond.Model != cond.Hwv {
//		q = query.NewTermQuery(invertIdx.Iterator("DeviceModel", generateKey2))
//	}
//	if q != nil {
//		return query.NewOrQuery([]query.Query{
//			query.NewAndQuery([]query.Query{
//				query.NewTermQuery(invertIdx.Iterator("NeedDeviceModel", "1")),
//				query.NewOrQuery([]query.Query{
//					query.NewTermQuery(invertIdx.Iterator("DeviceModel", cond.Make)),
//					query.NewTermQuery(invertIdx.Iterator("DeviceModel", generateKey1)),
//					q,
//				}, nil),
//			}, nil),
//			query.NewTermQuery(invertIdx.Iterator("NeedDeviceModel", "0")),
//		}, nil)
//	} else {
//		return query.NewOrQuery([]query.Query{
//			query.NewAndQuery([]query.Query{
//				query.NewTermQuery(invertIdx.Iterator("NeedDeviceModel", "1")),
//				query.NewOrQuery([]query.Query{
//					query.NewTermQuery(invertIdx.Iterator("DeviceModel", cond.Make)),
//					query.NewTermQuery(invertIdx.Iterator("DeviceModel", generateKey1)),
//				}, nil),
//			}, nil),
//			query.NewTermQuery(invertIdx.Iterator("NeedDeviceModel", "0")),
//		}, nil)
//	}
//}
//
////effectiveCountryCode
////len(campaign.effectiveCountryCode) > 0 campaign.effectiveCountryCode[condition.Country] == 1 or campaign.effectiveCountryCode["ALL"] == 1
////
////(campaign.NeedEffectiveCountry ==1 and (condition.Country_1 in campaign.effectiveCountryCode or
////"ALL_1" in campaign.effectiveCountryCode)) or campaign.NeedEffectiveCountry
//func effectiveCountryCodeQuery(idx *index.Indexer, cond *CampaignCondition) query.Query {
//	invertIdx := idx.GetInvertedIndex()
//	return query.NewOrQuery([]query.Query{
//		query.NewAndQuery([]query.Query{
//			query.NewTermQuery(invertIdx.Iterator("NeedEffectiveCountry", "1")),
//			query.NewOrQuery([]query.Query{
//				query.NewTermQuery(invertIdx.Iterator("EffectiveCountryCode", cond.Country_1)),
//				query.NewTermQuery(invertIdx.Iterator("EffectiveCountryCode", "ALL_1")),
//			}, nil),
//		}, nil),
//		query.NewTermQuery(invertIdx.Iterator("NeedEffectiveCountry", "0")),
//	}, nil)
//}
//
//// Gender NeedGender
//// (campaign.NeedGender ==1 and condition.RequestGender in campaign.Gender) or campaign.NeedGender == 0
//func genderQuery(idx *index.Indexer, cond *CampaignCondition) query.Query {
//	invertIdx := idx.GetInvertedIndex()
//	return query.NewOrQuery([]query.Query{
//		query.NewAndQuery([]query.Query{
//			query.NewTermQuery(invertIdx.Iterator("NeedGender", "1")),
//			query.NewTermQuery(invertIdx.Iterator("Gender", strconv.Itoa(condition.RequestGender))),
//		}, nil),
//		query.NewTermQuery(invertIdx.Iterator("NeedGender", "0")),
//	}, nil)
//
//}
//
////pkgName
//// ( campaign.NeedInventoryBlackList ==1 and condition.PkgName not in campaign.inventoryBlackList ) or
////campaign.NeedInventoryBlackList == 0
//func pkgNameQuery(idx *index.Indexer, cond *CampaignCondition) query.Query {
//	invertIdx := idx.GetInvertedIndex()
//	return query.NewOrQuery([]query.Query{
//		query.NewAndQuery([]query.Query{
//			query.NewTermQuery(invertIdx.Iterator("NeedInventoryBlackList", "1")),
//			query.NewTermQuery(invertIdx.Iterator("InventoryBlackList", cond.PkgName)),
//		}, nil),
//		query.NewTermQuery(invertIdx.Iterator("NeedInventoryBlackList", "0")),
//	}, nil)
//}
//
//// IsEcAdv  RetargetVisitorType
//// campaign.isEcAdv == 1 (mvutil.EcAdv） campaign.retargetVisitionType == 1 (include)
//func isecadvQuery(idx *index.Indexer, cond *CampaignCondition) query.Query {
//	invertIdx := idx.GetInvertedIndex()
//	if len(cond.CampaignIds) > 0{
//		var q []query.Query
//		for k, _ := range cond.CampaignIds {
//			q = append(q, query.NewTermQuery(invertIdx.Iterator("CampaignId", strconv.FormatInt(k, 10))))
//		}
//		return query.NewOrQuery([]query.Query{
//			query.NewTermQuery(invertIdx.Iterator("IsEcAdv", "0")),
//			query.NewAndQuery([]query.Query{
//				query.NewTermQuery(invertIdx.Iterator("IsEcAdv", "1")),
//				query.NewOrQuery([]query.Query{
//					query.NewAndQuery([]query.Query{
//						query.NewTermQuery(invertIdx.Iterator("RetargetVisitorType", "1")),
//						query.NewOrQuery(q, nil),
//					}, nil),
//					query.NewNotAndQuery([]query.Query{
//						query.NewTermQuery(invertIdx.Iterator("RetargetVisitorType", "2")),
//						query.NewOrQuery(q, nil),
//					}, nil),
//				}, nil),
//			}, nil),
//		}, nil)
//	}
//	return query.NewOrQuery([]query.Query{
//		query.NewTermQuery(invertIdx.Iterator("IsEcAdv", "0")),
//		query.NewAndQuery([]query.Query{
//			query.NewTermQuery(invertIdx.Iterator("IsEcAdv", "1")),
//			query.NewTermQuery(invertIdx.Iterator("RetargetVisitorType", "2")),
//		}, nil),
//	}, nil)
//}

//
//// iab
//// (campaign.NeedIabCategoryTag == 1 and condition.IabCategory 某个元素使用”-“split的第一个元素在campaign.IabCategoryTag1
////or 某个元素在campaign.IabCategryTag2) or campaign.NeedIabCategoryTag == 0
//func iabQuery(idx *index.Indexer, cond *CampaignCondition) query.Query {
//	invertIdx := idx.GetInvertedIndex()
//	//IabCategoryTag1 IabCategoryTag2
//	var iabCategoryTag1, iabCategoryTag2 []query.Query
//	if len(cond.IabCategory) != 0 {
//		for _, v := range cond.IabCategory {
//			tmp := strings.Split(v, "-")
//			iabCategoryTag1 = append(iabCategoryTag1, query.NewTermQuery(invertIdx.Iterator("IabCategoryTag1", tmp[0])))
//			iabCategoryTag2 = append(iabCategoryTag2, query.NewTermQuery(invertIdx.Iterator("IabCategoryTag1", v)))
//		}
//		return query.NewOrQuery([]query.Query{
//			query.NewAndQuery([]query.Query{
//				query.NewTermQuery(invertIdx.Iterator("NeedIabCategoryTag", "1")),
//				query.NewOrQuery(iabCategoryTag1, nil),
//				query.NewOrQuery(iabCategoryTag2, nil),
//			}, nil),
//			query.NewTermQuery(invertIdx.Iterator("NeedIabCategoryTag", "0")),
//		}, nil)
//	}
//	return nil
//}
//
//// UserAge
//// (campaign.NeedUserAge == 1 and condition.UserAge in campaign.UserAge) and campaign.NeedUserAge == 0
//func userAgeQuery(idx *index.Indexer, cond *CampaignCondition) query.Query {
//	invertIdx := idx.GetInvertedIndex()
//	return query.NewOrQuery([]query.Query{
//		query.NewAndQuery([]query.Query{
//			query.NewTermQuery(invertIdx.Iterator("NeedUserAge", "1")),
//			query.NewTermQuery(invertIdx.Iterator("UserAge", strconv.Itoa(cond.UserAge))),
//		}, nil),
//		query.NewTermQuery(invertIdx.Iterator("NeedUserAge", "0")),
//	}, nil)
//}
//
//// InstallApps
//// campaign.NeedInstallApps ==1 && campaign.InstallApps in condtion.InstallApp
//func installAppsQuery(idx *index.Indexer, cond *CampaignCondition) query.Query {
//	if len(cond.InstallApps) == 0 {
//		return nil
//	}
//	invertIdx := idx.GetInvertedIndex()
//	var q []query.Query
//	for k, v := range cond.InstallApps {
//		q = append(q, query.NewTermQuery(invertIdx.Iterator("InstallApps", strconv.Itoa(k))))
//	}
//	return query.NewAndQuery([]query.Query{
//		query.NewOrQuery(q, nil),
//		query.NewTermQuery(invertIdx.Iterator("NeedInstallApps", "1")),
//	}, nil)
//}
//
//// ExcludeInstalledApps
////exclude	campaign.NeedInstallApps ==1 && campaign.ExcludeInstalledApps not in condition.InstallApps
//func excludeInstalledAppsQuery(idx *index.Indexer, cond *CampaignCondition) query.Query {
//	if len(cond.InstallApps) == 0 {
//		return nil
//	}
//	invertIdx := idx.GetInvertedIndex()
//	var q []query.Query
//	for k := range cond.InstallApps {
//		q = append(q, query.NewTermQuery(invertIdx.Iterator("ExcludeInstalledApps", strconv.Itoa(k))))
//	}
//	return query.NewNotAndQuery([]query.Query{
//		query.NewTermQuery(invertIdx.Iterator("NeedInstallApps", "1")),
//		query.NewOrQuery(q, nil),
//	}, nil)
//}
//
//// adx
//// （（NeedAdxWhiteList = 1 and condition.Adx in  campaign.AdxWhiteList）or NeedAdxWhiteList = 0） and
////（（NeedAdxBlackList = 1 and conditon.Adx not in campaign.AdxBlackList）or NeedAdxBlackList = 0）
//func adxQuery(idx *index.Indexer, cond *CampaignCondition) query.Query {
//	invertIdx := idx.GetInvertedIndex()
//	if cond.Adx == tencent {
//		return query.NewTermQuery(invertIdx.Iterator("AdxInclude", cond.Adx))
//	}
//	return query.NewAndQuery([]query.Query{
//		query.NewOrQuery([]query.Query{
//			query.NewAndQuery([]query.Query{
//				query.NewTermQuery(invertIdx.Iterator("NeedAdxWhiteList", "1")),
//				query.NewTermQuery(invertIdx.Iterator("AdxWhiteList", cond.Adx)),
//			}, nil),
//			query.NewTermQuery(invertIdx.Iterator("NeedAdxWhiteList", "0")),
//		}, nil),
//		query.NewOrQuery([]query.Query{
//			query.NewNotAndQuery([]query.Query{
//				query.NewTermQuery(invertIdx.Iterator("NeedAdxBlackList", "1")),
//				query.NewTermQuery(invertIdx.Iterator("AdxBlackList", cond.Adx)),
//			}, nil),
//			query.NewTermQuery(invertIdx.Iterator("NeedAdxBlackList", "0")),
//		}, nil),
//	}, nil)
//}
//
//// advertiserAudit TODO
//// condition.Adx in campaign.AuditAdvertiserMap
//func advertiserAuditQuery(idx *index.Indexer, cond *CampaignCondition) query.Query {
//	invertIdx := idx.GetInvertedIndex()
//	if cond.NeedAdvAudit == true {
//		return query.NewTermQuery(invertIdx.Iterator("AuditAdvertiserMap", cond.Adx))
//	}
//	return nil
//}
//
//// creativeAudit TODO
//// condition.Adx in campaign.AuditCreativeMap
//func creativeAuditQuery(idx *index.Indexer, cond *CampaignCondition) query.Query {
//	invertIdx := idx.GetInvertedIndex()
//	if cond.NeedCreativeAudit == true {
//		return query.NewTermQuery(invertIdx.Iterator("AuditCreativeMap", cond.Adx))
//	}
//	return nil
//}
//
//// deviceType
//// len(campaign.DeviceTypeV2) == 0 or (4 in campaign.DeviceTypeV2 and 5 in campaign.DeviceTypeV2) or
////conditon.DeviceType in campaign.DeviceTypeV2
//func deviceTypeQuery(idx *index.Indexer, cond *CampaignCondition) query.Query {
//	invertIdx := idx.GetInvertedIndex()
//	if cond.DeviceType != 0 {
//		return query.NewOrQuery([]query.Query{
//			query.NewAndQuery([]query.Query{
//				query.NewTermQuery(invertIdx.Iterator("NeedDeviceType", "1")),
//				query.NewOrQuery([]query.Query{
//					query.NewAndQuery([]query.Query{
//						query.NewTermQuery(invertIdx.Iterator("DeviceType", "4")),
//						query.NewTermQuery(invertIdx.Iterator("DeviceType", "5")),
//					}, nil),
//					query.NewTermQuery(invertIdx.Iterator("DeviceType", strconv.Itoa(int(cond.DeviceType)))),
//				}, nil),
//			}, nil),
//			query.NewTermQuery(invertIdx.Iterator("NeedDeviceType", "0")),
//		}, nil)
//	}
//	return nil
//}
//
////// mobileCode
////func mobileCodeQuery(idx *index.Indexer, cond *CampaignCondition) query.Query {
////	storageIdx :=  idx.GetStorageIndex()
////	return query.NewOrQuery([]query.Query{
////		query.NewAndQuery([]query.Query{
////			query.NewTermQuery(storageIdx.Iterator("MobileCode")),
////		}, []check.Checker{
////			check.NewOrChecker([]check.Checker{
////				check.NewChecker(storageIdx.Iterator("MobileCode"), cond.carrier, operation.EQ, &operations{}, false),
////				check.NewInChecker(storageIdx.Iterator("MobileCode"), cond.Carrier, &operations{}, false),
////			}),
////		}),
////	}, nil)
////}
//
//// PackageName
//// campaign.PackageName not in condition.BApp
//func packageNameQuery(idx *index.Indexer, cond *CampaignCondition) query.Query {
//	invertIdx := idx.GetInvertedIndex()
//	var q []query.Query
//	for k := range cond.BApp {
//		q = append(q, query.NewTermQuery(invertIdx.Iterator("PackageName", k)))
//	}
//	if len(q) == 0 {
//		return nil
//	}
//	return query.NewOrQuery(q, nil)
//}
//
//// adSchedule
////“-2“ in campaign.AdSchedule or ("-1-curHour" in campaign.AdSchedule ) or (curDay-curHour in campaign.AdSchedule)
//func adScheduleQuery(idx *index.Indexer, cond *CampaignCondition) query.Query {
//	invertIdx := idx.GetInvertedIndex()
//	return query.NewOrQuery([]query.Query{
//		query.NewTermQuery(invertIdx.Iterator("AdSchedule", "-2")),
//		query.NewTermQuery(invertIdx.Iterator("AdSchedule", "-1-curHour")),
//		query.NewTermQuery(invertIdx.Iterator("AdSchedule", curDay-curHour)),
//	}, nil)
//}
//
////// UserInterest
////// len(campaign.UserInterestV2) == 0 or campaign.UserInterestV2的每个二级数组中至少有一个元素在condition.DmpInterests里
////func userInterestQuery(idx *index.Indexer, cond *CampaignCondition) query.Query {
////	storageIdx := idx.GetStorageIndex()
////	return query.NewAndQuery([]query.Query{
////		query.NewTermQuery(storageIdx.Iterator("UserInterest")),
////	}, []check.Checker{
////		check.NewInChecker(storageIdx.Iterator("UserInterest"), cond.DmpInterests, &operations{}, false),
////	})
////}
//
//// deviceId
//func deviceIdQuery(idx *index.Indexer, cond *CampaignCondition) query.Query {
//	invertIdx := idx.GetInvertedIndex()
//	return query.NewOrQuery([]query.Query{
//		query.NewNotAndQuery([]query.Query{
//			query.NewTermQuery(invertIdx.Iterator("RetargetingDevice", "1")),
//			query.NewTermQuery(invertIdx.Iterator("DeviceId", "")),
//		}, nil),
//		query.NewNotAndQuery([]query.Query{
//			query.NewTermQuery(invertIdx.Iterator("DeviceId", "")),
//			query.NewTermQuery(invertIdx.Iterator("RetargetingDevice", "1")),
//		}, nil),
//	}, nil)
//}
//
//func queryDsp(idx *index.Indexer, condition interface{}) query.Query {
//	cond, ok := condition.(*CampaignCondition)
//	if !ok {
//		return nil
//	}
//	//
//	//// invert list
//	//invertIdx := idx.GetInvertedIndex()
//
//	// storage list
//	storageIdx := idx.GetStorageIndex()
//
//	//// adSchedule
//	////“-2“ in campaign.AdSchedule or ("-1-curHour" in campaign.AdSchedule ) or (curDay-curHour in campaign.AdSchedule)
//	//var adSchedule = query.NewOrQuery([]query.Query{
//	//	query.NewTermQuery(invertIdx.Iterator("AdSchedule", "-2")),
//	//	query.NewTermQuery(invertIdx.Iterator("AdSchedule", "-1-curHour")),
//	//	query.NewTermQuery(invertIdx.Iterator("AdSchedule", curDay-curHour)),
//	//}, nil)
//
//	//// endTime startTime  这个在下面合并的时候拆开了写的
//	//var endTime = query.NewOrQuery([]query.Query{
//	//	query.NewAndQuery([]query.Query{
//	//		query.NewTermQuery(storageIdx.Iterator("EndTime")),
//	//	}, []check.Checker{
//	//		check.NewChecker(storageIdx.Iterator("EndTime"), time.Now().Unix(), operation.GT, nil, false),
//	//		check.NewChecker(storageIdx.Iterator("EndTime"), 0, operation.GT, nil, false),
//	//	}),
//	//	query.NewTermQuery(storageIdx.Iterator("EndTime")),
//	//}, []check.Checker{
//	//	check.NewChecker(storageIdx.Iterator("EndTime"), 0, operation.LE, nil, false),
//	//})
//
//	var andQuery *query.AndQuery
//	var notAndQuery *query.NotAndQuery
//	var queryAnd, queryNot []query.Query
//	if campaignIdQuery(idx, cond) != nil {
//		queryNot = append(queryNot, campaignIdQuery(idx, cond))
//	}
//	if campaignTypeQuery(idx, cond) != nil {
//		queryNot = append(queryNot, campaignTypeQuery(idx, cond))
//	}
//	if ctypeQuery(idx, cond) != nil {
//		queryNot = append(queryNot, ctypeQuery(idx, cond))
//	}
//	if directQuery(idx, cond) != nil {
//		queryNot = append(queryNot, directQuery(idx, cond))
//	}
//	if industryIdQuery(idx, cond) != nil {
//		queryNot = append(queryNot, industryIdQuery(idx, cond))
//	}
//	if domainQuery(idx, cond) != nil {
//		queryNot = append(queryNot, domainQuery(idx, cond))
//	}
//	if deviceAndIpuaRetargetQuery(idx, cond) != nil {
//		queryNot = append(queryNot, deviceAndIpuaRetargetQuery(idx, cond))
//	}
//	if appCategoryQuery(idx, cond) != nil {
//		queryNot = append(queryNot, appCategoryQuery(idx, cond))
//	}
//	if appSubCategoryQuery(idx, cond) != nil {
//		queryNot = append(queryNot, appSubCategoryQuery(idx, cond))
//	}
//	queryNot = append(queryNot, subCategoryV2Query(idx, cond))
//	queryNot = append(queryNot, excludeInstalledAppsQuery(idx, cond))
//	if packageNameQuery(idx, cond) != nil {
//		queryNot = append(queryNot, packageNameQuery(idx, cond))
//	}
//	if iabCategoryQuery(idx, cond) != nil {
//		queryNot = append(queryNot, iabCategoryQuery(idx, cond))
//	}
//	if subCategoryNameQuery(idx, cond) != nil {
//		queryNot = append(queryNot, subCategoryNameQuery(idx, cond))
//	}
//	if len(queryNot) != 0 {
//		notAndQuery = query.NewNotAndQuery(queryNot, nil)
//		queryAnd = append(queryAnd, notAndQuery)
//	}
//	queryAnd = append(queryAnd, osQuery(idx, cond))
//	queryAnd = append(queryAnd, countryCodeQuery(idx, cond))
//	if iabQuery(idx, cond) != nil {
//		queryAnd = append(queryAnd, iabQuery(idx, cond))
//	}
//	if installAppsQuery(idx, cond) != nil {
//		queryAnd = append(queryAnd, installAppsQuery(idx, cond))
//	}
//	queryAnd = append(queryAnd, trafficTypeQuery(idx, cond))
//	if adtypeQuery(idx, cond) != nil {
//		queryAnd = append(queryAnd, adtypeQuery(idx, cond))
//	}
//	if effectiveCountryCodeQuery(idx, cond) != nil {
//		queryAnd = append(queryAnd, effectiveCountryCodeQuery(idx, cond))
//	}
//	if supportHttpsQuery(idx, cond) != nil {
//		queryAnd = append(queryAnd, supportHttpsQuery(idx, cond))
//	}
//	queryAnd = append(queryAnd, genderQuery(idx, cond))
//	queryAnd = append(queryAnd, devicelanguageQuery(idx, cond))
//	queryAnd = append(queryAnd, adScheduleQuery(idx, cond))
//	queryAnd = append(queryAnd, networkTypeQuery(idx, cond))
//	queryAnd = append(queryAnd, deviceModelQuery(idx, cond))
//	queryAnd = append(queryAnd, isecadvQuery(idx, cond))
//	if advertiserIdQuery(idx, cond) != nil {
//		queryAnd = append(queryAnd, advertiserIdQuery(idx, cond))
//	}
//	queryAnd = append(queryAnd, adxQuery(idx, cond))
//	if advertiserAuditQuery(idx, cond) != nil {
//		queryAnd = append(queryAnd, advertiserAuditQuery(idx, cond))
//	}
//	if idx(idx, cond) != nil {
//		queryAnd = append(queryAnd, creativeAuditQuery(idx, cond))
//	}
//	if deviceTypeQuery(idx, cond) != nil {
//		queryAnd = append(queryAnd, deviceTypeQuery(idx, cond))
//	}
//	queryAnd = append(queryAnd, userAgeQuery(idx, cond))
//	queryAnd = append(queryAnd, pkgNameQuery(idx, cond))
//	queryAnd = append(queryAnd, pkgNameQuery(idx, cond))
//	queryAnd = append(queryAnd, deviceIdQuery(idx, cond))
//
//	var checkers = []check.Checker{
//		check.NewChecker(storageIdx.Iterator("OsVersionMin"), cond.osv, operation.LE, nil, false),
//		check.NewChecker(storageIdx.Iterator("OsVersionMax"), cond.osv, operation.GE, nil, false),
//		check.NewOrChecker([]check.Checker{
//			check.NewAndChecker([]check.Checker{
//				check.NewChecker(storageIdx.Iterator("EndTime"), time.Now().Unix(), operation.GT, nil, false),
//				check.NewChecker(storageIdx.Iterator("EndTime"), 0, operation.GT, nil, false),
//			}),
//			check.NewChecker(storageIdx.Iterator("EndTime"), 0, operation.LE, nil, false),
//		}),
//		check.NewChecker(storageIdx.Iterator("ContentRating"), 12, operation.LE, nil, false),
//		check.NewOrChecker([]check.Checker{
//			check.NewChecker(storageIdx.Iterator("MobileCode"), cond.carrier, operation.EQ, &operations{}, false),
//			check.NewInChecker(storageIdx.Iterator("MobileCode"), cond.Carrier, &operations{}, false),
//		}),
//		check.NewInChecker(storageIdx.Iterator("UserInterest"), cond.DmpInterests, &operations{}, false),
//	}
//	if len(queryAnd) != 0 {
//		andQuery = query.NewAndQuery(queryAnd, checkers)
//		return andQuery
//	}
//	return nil
//
//	//resQuery := query.NewAndQuery([]query.Query{
//	//	query.NewNotAndQuery([]query.Query{
//	//		campaignIdQuery(idx, cond),
//	//		campaignTypeQuery(idx, cond),
//	//		ctypeQuery(idx, cond),
//	//		directQuery(idx, cond),
//	//		isecadvQuery(idx, cond),
//	//		industryIdQuery(idx, cond),
//	//		domainQuery(idx, cond),
//	//		deviceAndIpuaRetargetQuery(idx, cond),
//	//		appCategoryQuery(idx, cond),
//	//		appSubCategoryQuery(idx, cond),
//	//		subCategoryV2Query(idx, cond),
//	//		excludeInstalledAppsQuery(idx, cond),
//	//		packageNameQuery(idx, cond),
//	//		iabCategoryQuery(idx, cond),
//	//		subCategoryNameQuery(idx, cond),
//	//	}, nil),
//	//	osQuery(idx, cond),
//	//	countryCodeQuery(idx, cond),
//	//	iabQuery(idx, cond),
//	//	installAppsQuery(idx, cond),
//	//	trafficTypeQuery(idx, cond),
//	//	adtypeQuery(idx, cond),
//	//	effectiveCountryCodeQuery(idx, cond),
//	//	supportHttpsQuery(idx, cond),
//	//	genderQuery(idx, cond),
//	//	devicelanguageQuery(idx, cond),
//	//	adScheduleQuery(idx, cond),
//	//	networkTypeQuery(idx, cond),
//	//	deviceModelQuery(idx, cond),
//	//	advertiserIdQuery(idx, cond),
//	//	adxQuery(idx, cond),
//	//	advertiserAuditQuery(idx, cond),
//	//	creativeAuditQuery(idx, cond),
//	//	deviceTypeQuery(idx, cond),
//	//	userAgeQuery(idx, cond),
//	//	pkgNameQuery(idx, cond),
//	//	mvAppIdQuery(idx, cond),
//	//}, []check.Checker{
//	//	check.NewChecker(storageIdx.Iterator("OsVersionMin"), cond.osv, operation.LE, nil, false),
//	//	check.NewChecker(storageIdx.Iterator("OsVersionMax"), cond.osv, operation.GE, nil, false),
//	//	check.NewOrChecker([]check.Checker{
//	//		check.NewAndChecker([]check.Checker{
//	//			check.NewChecker(storageIdx.Iterator("EndTime"), time.Now().Unix(), operation.GT, nil, false),
//	//			check.NewChecker(storageIdx.Iterator("EndTime"), 0, operation.GT, nil, false),
//	//		}),
//	//		check.NewChecker(storageIdx.Iterator("EndTime"), 0, operation.LE, nil, false),
//	//	}),
//	//	check.NewChecker(storageIdx.Iterator("ContentRating"), 12, operation.LE, nil, false),
//	//	check.NewOrChecker([]check.Checker{
//	//		check.NewChecker(storageIdx.Iterator("MobileCode"), cond.carrier, operation.EQ, &operations{}, false),
//	//		check.NewInChecker(storageIdx.Iterator("MobileCode"), cond.Carrier, &operations{}, false),
//	//	}),
//	//	check.NewInChecker(storageIdx.Iterator("UserInterest"), cond.DmpInterests, &operations{}, false),
//	//})
//	//
//
//	// ******  Search  ******
//	//searcher := search.NewSearcher()
//	//searcher.Search(idx, resQuery)
//	//fmt.Println(searcher.Docs)
//	//fmt.Println(searcher.Time)
//}
//
//func fff() {
//	var j = query.NewJSONFormatter()
//	var c *CampaignCondition
//	var idx *index.Indexer
//	_, _ = j.Marshal(c)
//	_, _ = j.Unmarshal("aa", idx, c, queryDsp)
//}
