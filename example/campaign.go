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


//filed: campaignId, countryCode, adCreative, time,state
/*type  TencentCreatives struct {
	CampaignId  int64  `bson:"campaignId,omitempty" json:"campaignId"`
	CountryCode string `bson:"countryCode,omitempty" json:"countryCode"`
	AdCreative []AdCreativeAudit `bson:"adCreative,omitempty" json:"adCreative"`
}*/

// GetAdvertiserId 动态获取广告主id，目前以packageName作为广告主标识，待adn有标识唯一广告主字段，再修改
func (campaign *CampaignDetail) GetAdvertiserId() interface{} {
	return campaign.PackageName
}

type CreativeAudit struct {
	CampaignId    int64                                                 `bson:"campaignId,omitempty" json:"campaignId"`
	PackageName   string                                                `bson:"packageName,omitempty" json:"packageName"`
	AuditElements map[string]map[string]map[string]CreativeAuditElement `bson:"audit,omitempty" json:"audit"` //key1=adx,key2=creativeId,key3=templateId,value=素材审核信息
	Updated       int64                                                 `bson:"updated,omitempty" json:"updated"`
}

type CreativeAuditElement struct {
	Status       string `bson:"status,omitempty" json:"status"`
	ExCreativeId string `bson:"exCreativeId,omitempty" json:"exCreativeId"`
	Category     string `bson:"category,omitempty" json:"category"`
}

type ExcludeRule struct {
	Bwtype        int64   `bson:"type,omitempty" json:"type"`
	IncludeAppIds []int64 `bson:"includeAppIds,omitempty" json:"includeAppIds"`
	ExcludeAppIds []int64 `bson:"excludeAppIds,omitempty" json:"excludeAppIds"`
	Status        int64   `bson:"status,omitempty" json:"status"`
}

/*DspTemplateConfElement  是作为模板元素结构体
Adtype 流量类型
TemplateGroup  投放的模板类型
Status		状态，1，active；2，pause
*/
type DspTemplateConfElement struct {
	Adtype string                `bson:"adtype,omitempty" json:"adtype"`               //adtype 类型
	TGroup []model.TemplateGroup `bson:"templateGroup,omitempty" json:"templateGroup"` //模板类型
	Status int                   `bson:"status,omitempty" json:"status"`               //status类型 1.active 2.pause
}

type TrackUrl struct {
	Sec int
	Url string
}

// 库存类别
type InventoryV2 struct {
	//Type        int                 `bson:"type,omitempty" json:"type"`
	IabCategory map[string][]string `bson:"iabCategory,omitempty" json:"iabCategory"`
	//AdType      []string            `bson:"adtype,omitempty" json:"adtype"`
	AdnAdtype []int    `bson:"adnAdtype,omitempty" json:"adnAdtype"`
	AppSite   []string `bson:"app_site,omitempty" json:"app_site"`
}

//1、Creative素材结构:Image, video, encard, mraid,jstag 素材map中的key为size,  endcard和基础素材的map的key为素材的类型.
type SCreative struct {
	Image             map[string][]SImage // key=string size value = 同size的多个SImage
	DyImage           map[string][]SImage
	JsImage           map[string][]SJsImage
	MraidImage        map[string][]SJsImage
	Video             map[string][]SVideo
	Common            map[string][]SCommon           //基础素材, key 基础素材类型 401-407
	GroupCreativesIdx map[string][]SGroupCreativeIdx //key表示CreativeSpec分组素材
	//素材区分ADX， adxName -> creativeSpecId -> SGroupCreativeIdx 对应的素材组集合
	AdxCreativeSpecGroupCreativesIdx map[string]map[string][]SGroupCreativeIdx
}

type SGroupCreative struct {
	GroupId   int64         //组Id
	Creatives []interface{} //分组素材 //SImage/SVideo/SJsImage/SCommon
}

type IUrlCreative interface {
	ICreative
	GetUrl() string
	SetUrl(string)
}

type ITagCodeCreative interface {
	ICreative
	GetTagCode() string
	SetTagCode(string)
}

type IVideoCreative interface {
	ICreative
	GetVideoLength() int
	SetVideoLength(int)
	GetVideoResolution() string
	SetVideoResolution(string)
	GetWidth() int
	SetWidth(int)
	GetHeight() int
	SetHeight(int)
	GetBitRate() int
	SetBitRate(int)
	GetScreenShot() string
	SetScreenShot(string)
}

type ICreative interface {
	GetMime() []string
	SetMime([]string)
	GetAttribute() []uint8
	SetAttribute([]uint8)
	GetResolution() string
	SetResolution(string)
	GetResourceType() uint8
	SetResourceType(uint8)
	GetShowType() uint8
	SetShowType() uint8
	GetTsCreativeId() int64
	SetTsCreativeId(int64)
	GetCreativeType() int
	SetCreativeType() int
	GetFMd5() string
	SetFMd5(string)
}

//1、静态动态图片结构
type SImage struct {
	SubResourceType uint8
	FormatType      uint8
	Url             string
	Mime            []string
	Attribute       []uint8
	Resolution      string
	ResourceType    uint8
	ShowType        uint8
	TsCreativeId    string
	CreativeId      int64
	CreativeType    int
	CommonField     SCommonField
}

//2、jstag，mraid图片结构
type SJsImage struct {
	SubResourceType uint8
	FormatType      uint8
	TagCode         string
	Mime            []string
	Attribute       []uint8
	Resolution      string
	ResourceType    uint8
	ShowType        uint8
	TsCreativeId    string
	CreativeId      int64
	CreativeType    int
	CommonField     SCommonField
}

//3、视频图片结构
type SVideo struct {
	Url             string
	VideoLength     int
	VideoResolution string
	Width           int
	Height          int
	Bitrate         int
	ScreenShot      string
	Mime            []string
	Attribute       []uint8
	SubResourceType uint8
	ResourceType    uint8
	ShowType        uint8
	TsCreativeId    string
	CreativeId      int64
	CreativeType    int
	FormatType      uint8
	CommonField     SCommonField
}

type SCommonField struct {
	FMd5         string
	PackageName  string
	CountryCode  string
	Source       int64
	TemplateType int64
	Score        int64
	Orientation  int64
	Clarity      int64
	VideoSize    int64
	Utime        int64
	Ctime        int64
	Cname        string
}

//5、基础素材结构
type SCommon struct {
	Value           string
	Mime            []string
	Attribute       []uint8
	SubResourceType uint8
	ResourceType    uint8
	ShowType        uint8
	TsCreativeId    string
	CreativeId      int64
	CreativeType    int
	FormatType      uint8
	CommonField     SCommonField
}

// Image,JsImage,MraidImage 可以对应多个resourceType
// 参考：
// 字段       -> resourceType
// Image      -> SubResourceTypeImage,SubResourceTypeStaticEndCard
// JsImage    -> SubResourceTypeImage,SubResourceTypeDynamicEndCard,SubResourceTypeRichMediaShow,SubResourceTypeTplPhotoGallery,SubResourceTypeTplAnimation,SubResourceTypeTplVideo
// MraidImage -> SubResourceTypeRichMediaShow,DeprecatedMraidCreative,SubResourceTypeTplFallingObjectGame
type SRankCreativeIdx struct {
	//同一SRankCreativeIdx可能有多个ResourceType，SubResourceType不做获取素材的标识
	//建议使用FormatType区分
	SubResourceType uint8
	FormatType      uint8 //image,dynamic_image,js,mraid,video,common
	Creativeids     []int64
	Image           map[int64]*SImage //key为图片的CreativeId
	DyImage         map[int64]*SImage
	JsImage         map[int64]*SJsImage
	MraidImage      map[int64]*SJsImage
	Video           map[int64]*SVideo
	Common          map[int64]*SCommon        //基础素材,key creativeid
	GroupCreative   map[int64]*SGroupCreative //素材组
	PickRType       uint8                     //选取的image类型
	PickRTypes      map[uint8]bool
	//ResourceType    uint8                    // tencent_adx 不初始化此类型
}

type SRankTemplate struct {
	Id             string           //模板id, 父级id
	Weight         float64          //模板权重，基于单子的。
	CreativeSec    string           //广告位支持的素材规格。
	AdElemType     []int32          //模板元素类型
	CreativeGroups []SCreativeGroup //一个模板对应多个素材组,每个素材组有多个素材
}

type SCreativeGroup struct {
	Id        int64         //素材组id， 送审每组素材唯一id.
	Creatives []interface{} // interface为素材结构体： SImage, SJstag, SVideo, SComon
}

func (srci *SRankCreativeIdx) IsNotEmpty() bool {
	return len(srci.Creativeids) > 0
}

type CreativeIdx struct {
	CamapaignIdx map[int64]map[string]*AdCreative  //map 的key: campaigned, countrycode
	PackageIdx   map[string]map[string]*AdCreative //map 的key: packagename, countrycode
}

type AdCreative struct {
	PackageName string `bson:"packageName,omitempty" json:"packageName"`
	Status      int64  `bson:"status,omitempty" json:"status"`
	CampaignId  int64  `bson:"campaignId,omitempty" json:"campaignId"`
	CountryCode string `bson:"countryCode,omitempty" json:"countryCode"`
	Content     map[string]map[string]map[string]interface{} `bson:"content,omitempty" json:"content"`
	Updated int64 `bson:"updated,omitempty" json:"updated"`
	//Content     string `bson:"content,omitempty" json:"content"`
	//Json        SContent `bson:"json,omitempty" json:"json"`
	//Json map[string]interface{} `bson:"json,omitempty" json:"json"`
}

type AdvertiserAudit struct {
	AdvertiserId        int64                             `bson:"advertiserId,omitempty"`
	PackageName         string                            `bson:"packageName,omitempty"`
	AdvertiserAuditInfo map[string]AdvertiserAuditElement `bson:"audit,omitempty"` //key=adx,value=auditElement
	Updated             int64                             `bson:"updated"`
}

func (adv *AdvertiserAudit) GetAdvertiserId() interface{} {
	return adv.PackageName
}

type AdvertiserAuditElement struct {
	Status         string `bson:"status,omitempty"`
	ExAdvertiserId string `bson:"exAdvertiserId,omitempty"`
	ExCategory     string `bosn:"extCat,omitempty"`
}

type Creative struct {
	CreativeId   int64     `bson:"creativeId,omitempty" json:"creativeId"`
	Lang         int       `bson:"lang,omitempty" json:"lang"`
	Type         int8      `bson:"type,omitempty" json:"type"`
	Width        int16     `bson:"width,omitempty" json:"width"`
	Height       int16     `bson:"height,omitempty" json:"height"` //update int16 for mobfox
	ImageSize    string    `bson:"imageSize,omitempty" json:"imageSize"`
	Name         string    `bson:"name,omitempty" json:"name"`
	ImageUrl     string    `bson:"imageUrl,omitempty" json:"imageUrl"`
	Status       int8      `bson:"status,omitempty" json:"status"`
	TextVideo    TextVideo `bson:"textVideo,omitempty" json:"textVideo"`       //创意中关于视频的存储结构
	VideoUrl     string    `bson:"videoUrl,omitempty" json:"videoUrl"`         // 视频的url
	ResourceType uint8     `bson:"resourceType,omitempty" json:"resourceType"` // 枚举类型1:Image 2:Video 3:Mraid
	Mime         []string  `bson:"mime,omitempty" json:"mime"`                 //以下为etl新增字段
	Attribute    []uint8   `bson:"attribute,omitempty" json:"attribute"`       //creative_attribute_list
	TemplateType int       `bson:"templateType,omitempty" json:"templateType"` //rich meadia 1:animation 2:video  3:falling object game  4:photo gallery
	TagCode      string    `bson:"tagCode,omitempty" json:"tagCode"`           //img-jstag   video-js tag  / mraid tag
	ShowType     uint8     `bson:"showType,omitempty" json:"showType"`
}
type TextVideo struct {
	VideoLength     int    `bson:"videoLength,omitempty" json:"videoLength"`         //视频长度 seconds
	VideoResolution string `bson:"videoResolution,omitempty" json:"VideoResolution"` //分辨率
}

type ChanlPrice struct {
	Chanl string  `bson:"chanl,omitempty" json:"chanl"` //渠道
	Price float64 `bson:"price,omitempty" json:"price"` //价格
}

//filed: campaignId, countryCode, adCreative, time,state
type TencentCreatives struct {
	CampaignId  int64              `bson:"campaignId,omitempty" json:"campaignId"`
	CountryCode string             `bson:"countryCode,omitempty" json:"countryCode"`
	PackageName string             `bson:"packageName,omitempty" json:"packageName"`
	AdCreative  []SAdCreativeAudit `bson:"adCreative,omitempty" json:"adCreative"`
	//status = 1 有效uint8
}

type SAdCreativeAudit struct {
	CreativeSpec  int32               `bson:"creativeSpec,omitempty" json:"creativeSpec"`
	CreativeGroup []SGroupCreativeIdx `bson:"creativeGroup,omitempty" json:"creativeGroup"`
}

type SGroupCreativeIdx struct {
	Id        string         `bson:"id,omitempty" json:"id"`
	Attribute string         `bson:"attribute,omitempty" json:"attribute"`
	Industry  string         `bson:"industry,omitempty" json:"industry"`
	Creatives []SCreativeAdt `bson:"creatives,omitempty" json:"creatives"`
	Status    uint8          `bson:"status,omitempty" json:"status"` // status 枚举值： 0准备中，1待审核，2通过，3拒绝，4信息变更，5送审失败，6预审通过，11重新送审，44删除。 其中先审后投2是通过，先投后审2和6通过
}

type SCreativeAdt struct {
	Id           int64  `bson:"id,omitempty" json:"id"`
	Resolution   string `bson:"rs,omitempty" json:"rs"` //素材尺寸,如果素材的类型common 则rs为creativeType
	ResourceType uint8  `bson:"rt,omitempty" json:"rt"` //image, dyimage, video, jstag, mraid 素材和format, common素材和resourceType对应
	// ResourceType 的详细值， 见http://confluence.mobvista.com/pages/viewpage.action?pageId=16696246 的15节
}

type CampaignAdxCreativeAudit struct {
	CampaignId  int64  `bson:"campaignId,omitempty" json:"campaignId"`
	CountryCode string `bson:"countryCode,omitempty" json:"countryCode"`
	PackageName string `bson:"packageName,omitempty" json:"packageName"`
	//Campaign对应的adx素材审核对象
	AdxCreativeAuditSlice []AdxCreativeAudit `bson:"adx,omitempty" json:"adx"`
}

//adx素材审核对象
type AdxCreativeAudit struct {
	//adx 名字，例如 toutiao
	Name     string `bson:"name,omitempty" json:"name"`
	Industry int64  `bson:"industry,omitempty" json:"industry"`
	//每个adx对应多个AdCreative,即多个广告位unit
	CreativeSpecSlice []CreativeSpec `bson:"adCreative,omitempty" json:"adCreative"`
}

//每个广告位unit和对应的所有素材组,国内通用，始于头条ADX，
type CreativeSpec struct {
	Id            string              `bson:"creativeSpec,omitempty" json:"creativeSpec"`
	CreativeGroup []SGroupCreativeIdx `bson:"creativeGroup,omitempty" json:"creativeGroup"`
}
