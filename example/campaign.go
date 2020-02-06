package main

type CampaignDetail struct {
	CampaignId            int64   `bson:"campaignId,omitempty" json:"campaignId"`
	AdvertiserId          int64   `bson:"advertiserId,omitempty" json:"advertiserId"`
	Name                  string  `bson:"name,omitempty" json:"name"`
	Platform              uint8   `bson:"platform,omitempty" json:"platform"`
	Price                 float64 `bson:"price,omitempty" json:"price"`
	DeviceAndIpuaRetarget int8    `bson:"gaid_idfa_needs,omitempty" json:"gaid_idfa_needs"`
	//IabCategory  []string   `bson:"iabCategory,omitempty" json:"iabCategory"`
	IabCategory          map[string][]string `bson:"iabCategory,omitempty" json:"iabCategory"`
	CountryCode          []string            `bson:"countryCode,omitempty" json:"countryCode"`
	EffectiveCountryCode map[string]int      `bson:"effectiveCountryCode,omitempty" json:"effectiveCountryCode"` // 分地区预算控制。key：为ALL或具体国家编码，value为1或其他，为 1 表示预算充足。key["ALL"] =1 or 2;key[country="IN"] = 1 or 2
	Tag                  int                 `bson:"tag,omitempty" json:"tag"`
	Status               int8                `bson:"status,omitempty" json:"status"`
	PackageName          string              `bson:"packageName,omitempty" json:"packageName"`
	CampaignType         int8                `bson:"campaignType,omitempty" json:"campaignType"`
	Ctype                int8                `bson:"ctype,omitempty" json:"ctype"`
	AppName              string              `bson:"appName,omitempty" json:"appName"`
	IconUrl              string              `bson:"iconUrl,omitempty" json:"iconUrl"`
	AppDesc              string              `bson:"appDesc,omitempty" json:"appDesc"`
	AppSize              string              `bson:"appSize,omitempty" json:"appSize"`
	AppScore             float64             `bson:"appScore,omitempty" json:"appScore"`
	AppInstall           uint32              `bson:"appInstall,omitempty" json:"appInstall"`
	Category             uint8               `bson:"category,omitempty" json:"category"`
	//Creative     []Creative          `bson:"creative,omitempty" json:"creative"`
	MapCreative  map[string]SCreative //key=country value= creative
	OsVersionMin int                  `bson:"osVersionMinV2,omitempty" json:"osVersionMinV2"`
	OsVersionMax int                  `bson:"osVersionMaxV2,omitempty" json:"osVersionMaxV2"`
	AdSourceId   int8                 `bson:"adSourceId,omitempty" json:"adSourceId"`
	//ImageSize    []string `bson:"imageSize,omitempty" json:"imageSize"`
	PublisherId int64  `bson:"publisherId,omitempty" json:"publisherId"`
	Updated     int64  `bson:"updated,omitempty" json:"updated"`
	Domain      string `bson:"appDomain,omitempty" json:"appDomain"`
	Developer   string `bson:"developer,omitempty" json:"developer"`
	DirectPkg   string `bson:"directPackageName,omitempty" json:"directPackageName"`
	ApkUrl      string `bson:"apkUrl,omitempty" json:"apkUrl"`
	//follows are for brand campain
	System        []int            `bson:"system,omitempty" json:"system"`               //3 -m系统 5 dsp
	DeviceTypeV2  []uint8          `bson:"deviceTypeV2,omitempty" json:"deviceTypeV2"`   //4 phone 5 tablet
	CityCode      map[string][]int `bson:"cityCode,omitempty" json:"citycode"`           //城市
	NetworkTypeV2 []uint8          `bson:"networkTypeV2,omitempty" json:"networkTypeV2"` //0 all 9 wifi
	MobileCode    []string         `bson:"mobileCode,omitempty" json:"mobileCode"`       //运营商all
	//DeviceMake   []string         `bson:"devicemake,omitempty" json:"devicemake"`   //
	//DeviceModel       []string `bson:"deviceModel,omitempty" json:"devicemodel"`

	DeviceModelV3      map[string][]string `bson:"deviceModelV3,omitempty" json:"deviceModelV3"`
	TrackingUrl        string              `bson:"trackingUrl,omitempty" json:"trackingurl"`
	FrequencyCap       int                 `bson:"frequencyCap,omitempty" json:"frequencyCap"`
	AdUrlList          []string            `bson:"adUrlList,omitempty" json:"adUrlList"`
	AdvImp             []TrackUrl          `bson:"advImp,omitempty" json:"advImp"`
	UserInterest       map[string][]string `bson:"userInterest,omitempty" json:"userInterest"`
	GenderV2           []int               `bson:"genderV2,omitempty" json:"genderV2"`
	InventoryV2        InventoryV2         `bson:"inventoryV2,omitempty" json:"inventoryV2"`
	DeviceId           string              `bson:"deviceId,omitempty"  json:"deviceId"`
	TrafficType        []uint8             `bson:"trafficType,omitempty" json:"trafficType"`
	VbaTrackingLink    string              `bson:"vbaTrackingLink,omitempty" json:"vbaTrackingLink"`
	VbaConnecting      int                 `bson:"vbaConnecting,omitempty" json:"vbaConnecting"`
	RetargetingDevice  int                 `bson:"retargetingDevice,omitempty" json:"retargetingDevice"` // 1.Yes 2.No   默认 No。
	AdSchedule         map[string][]int    `bson:"adSchedule,omitempty" json:"adSchedule"`
	UserAgeV2          []int               `bson:"userAgeV2,omitempty" json:"userAgeV2"`
	StartTime          int64               `bson:"startTime,omitempty" json:"startTime"`
	EndTime            int64               `bson:"endTime,omitempty" json:"endTime"`
	BudgetFirst        bool                `bson:"budgetFirst,omitempty" json:"budgetFirst"`
	ExcludeRule        ExcludeRule         `bson:"excludeRule,omitempty" json:"excludeRule"`
	InventoryBlackList []string            `bson:"inventoryBlackList,omitempty" json:"inventoryBlackList"`
	PreviewUrl         string              `bson:"previewUrl,omitempty" json:"previewUrl"` //storeUrl
	IsCampaignCreative int                 `bson:"isCampaignCreative,omitempty" json:"isCampaignCreative"`
	OriPrice           float64             `bson:"oriPrice,omitempty" json:"oriPrice"`
	CostType           int8                `bson:"costType,omitempty" json:"costType"`
	Source             int                 `bson:"source,omitempty" json:"source"`
	//for convert Category  to string
	Direct         int    `bson:"direct,omitempty" json:"direct"`               // 单子类型，direct - 1 直单，直接客户单子；direct - 2 非直单，二手单
	ContentRating  int    `bson:"contentRating,omitempty" json:"contentRating"` // 年龄分级
	ThirdParty     string `bson:"thirdParty,omitempty" json:"thirdParty"`       // 单子类型，
	SubCategoryId  int    `bson:"subCategoryId,omitempty" json:"subCategoryId"`
	SubCategoryV2  []int  `bson:"subCategoryV2,omitempty" json:"subCategoryV2"`
	CategoryStr    string //应用市场的分类
	SubCategoryStr string //应用市场的二级分类
	PackageDevice  string `bson:"packageDevice,omitempty" json:"packageDevice"`
	TotalBudget    int64  `bson:"totalBudget,omitempty" json:"totalBudget"` //总预算
	DailyBudget    int64  `bson:"dailyBudget,omitempty" json:"dailyBudget"` //日预算
	//过滤market分类，目前仅用于京东exclude_category
	SubCategoryName []string `bson:"subCategoryName,omitempty"`
	//模板deviceLanguage定向
	InstallApps          []int                             `bson:"installApps,omitempty" json:"installApps"`
	ExcludeInstalledApps []int                             `bson:"excludeInstalledApps,omitempty" json:"excludeInstalledApps"`
	UserInterestV2       [][]int                           `bson:"userInterestV2,omitempty" json:"userInterestV2"`
	DeviceLanguage       map[string][]string               `bson:"deviceLanguage,omitempty" json:"deviceLanguage"`
	AdxWhiteBlack        map[string][]string               `bson:"adxWhiteBlack,omitempty" json:"adxWhiteBlack"`
	IsEcAdv              int                               `bson:"isEcAdv,omitempty" json:"isEcAdv"`
	RetargetVisitorType  int                               `bson:"retargetVisitorType,omitempty" json:"retargetVisitorType"`
	//for audit
	AdxInclude         []string                     `bson:"adxInclude,omitempty" json:"adxInclude"`
	IndustryId         int64                        `bson:"industryId,omitempty" json:"industryId"` //广告主行业Id
	AuditIndusty       map[string]int64             //new add 送审adx行业
	AuditAdvertiserMap map[string]string            // key=adx,value=exAdvertiserId
	AuditCreativeMap   map[string]map[string]string // key1=adx,key2="creativeId:templateId",value=exCreativeId
	//AdCreativeAuditIdx map[string]map[string]GroupCreative   //key1: country, key2: creativeSpec, val: GroupCreative
	OpenType               int8               `bson:"openType,omitempty" json:"openType"`
	CountryChannelPrice    map[string]float64 `bson:"dspCountryChanlPrice,omitempty" json:"dspCountryChanlPrice"`
	CountryChannelOriPrice map[string]float64 `bson:"dspCountryChanlOriPrice,omitempty" json:"dspCountryChanlOriPrice"`
}

