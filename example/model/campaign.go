package model

const (
	CampaignSourceOfferSync     = 0
	CampaignSourceAdnPortal     = 1
	CampaignSourceSSPlatform    = 2
	CampaignSourceSSAdvPlatform = 3
)

type CampaignInfo struct {
	CampaignId         int64                        `bson:"campaignId,omitempty" json:"campaignId,omitempty"`
	AdvertiserId       *int32                       `bson:"advertiserId,omitempty" json:"advertiserId,omitempty"`
	TrackingUrl        string                       `bson:"trackingUrl,omitempty" json:"trackingUrl,omitempty"`
	TrackingUrlHttps   string                       `bson:"trackingUrlHttps,omitempty" json:"trackingUrlHttps,omitempty"`
	DirectUrl          string                       `bson:"directUrl,omitempty" json:"directUrl,omitempty"`
	Price              *float64                     `bson:"price,omitempty" json:"price,omitempty"`
	OriPrice           *float64                     `bson:"oriPrice,omitempty" json:"oriPrice,omitempty"`
	CityCodeV2         map[string][]int32           `bson:"cityCodeV2,omitempty" json:"cityCodeV2,omitempty"`
	Status             int32                        `bson:"status,omitempty" json:"status,omitempty"`
	Network            *int32                       `bson:"network,omitempty" json:"network,omitempty"`
	PreviewUrl         *string                      `bson:"previewUrl,omitempty" json:"previewUrl,omitempty"`
	PackageName        string                       `bson:"packageName,omitempty" json:"packageName,omitempty"`
	CampaignType       *int32                       `bson:"campaignType,omitempty" json:"campaignType,omitempty"`
	Ctype              *int32                       `bson:"ctype,omitempty" json:"ctype,omitempty"`
	AppSize            *string                      `bson:"appSize,omitempty" json:"appSize,omitempty"`
	Tag                *int32                       `bson:"tag,omitempty" json:"tag,omitempty"`
	AdSourceId         *int32                       `bson:"adSourceId,omitempty" json:"adSourceId,omitempty"`
	PublisherId        *int64                       `bson:"publisherId,omitempty" json:"publisherId,omitempty"`
	PreClickCacheTime  *int32                       `bson:"preClickCacheTime,omitempty" json:"preClickCacheTime,omitempty"`
	FrequencyCap       *int32                       `bson:"frequencyCap,omitempty" json:"frequencyCap,omitempty"`
	DirectPackageName  string                       `bson:"directPackageName,omitempty" json:"directPackageName,omitempty"`
	SdkPackageName     string                       `bson:"sdkPackageName,omitempty" json:"sdkPackageName,omitempty"`
	AdvImp             *[]AdvImp                    `bson:"advImp,omitempty" json:"advImp,omitempty"`
	AdUrlList          *[]string                    `bson:"adUrlList,omitempty" json:"adUrlList,omitempty"`
	JumpType           int32                        `bson:"jumpType,omitempty" json:"jumpType,omitempty"`
	VbaConnecting      *int32                       `bson:"vbaConnecting,omitempty" json:"vbaConnecting,omitempty"`
	VbaTrackingLink    *string                      `bson:"vbaTrackingLink,omitempty" json:"vbaTrackingLink,omitempty"`
	RetargetingDevice  *int32                       `bson:"retargetingDevice,omitempty" json:"retargetingDevice,omitempty"`
	SendDeviceidRate   *int32                       `bson:"sendDeviceidRate,omitempty" json:"sendDeviceidRate,omitempty"`
	Endcard            *map[string]EndCard          `bson:"endcard,omitempty" json:"endcard,omitempty"`
	Loopback           *LoopBack                    `bson:"loopback,omitempty" json:"loopback,omitempty"`
	BelongType         *int32                       `bson:"belongType,omitempty" json:"belongType,omitempty"`
	ConfigVBA          *ConfigVBA                   `bson:"configVBA,omitempty" json:"configVBA,omitempty"`
	AppPostList        *AppPostList                 `bson:"appPostList,omitempty" json:"appPostList,omitempty"`
	BlackSubidListV2   map[string]map[string]string `bson:"blackSubidListV2,omitempty" json:"blackSubidListV2,omitempty"`
	BtV4               *BtV4                        `bson:"btV4,omitempty" json:"btV4,omitempty"`
	OpenType           *int32                       `bson:"openType,omitempty" json:"openType,omitempty"`
	SubCategoryName    *[]string                    `bson:"subCategoryName,omitempty" json:"subCategoryName,omitempty"`
	IsCampaignCreative *int32                       `bson:"isCampaignCreative,omitempty" json:"isCampaignCreative,omitempty"`
	CostType           *int32                       `bson:"costType,omitempty" json:"costType,omitempty"`
	Source             *int32                       `bson:"source,omitempty" json:"source,omitempty"`
	JUMPTYPECONFIG     map[string]int32             `bson:"JUMP_TYPE_CONFIG,omitempty" json:"JUMP_TYPE_CONFIG,omitempty"`
	ChnID              *int                         `bson:"chnId,omitempty" json:"chnId,omitempty"`
	ThirdParty         string                       `bson:"thirdParty,omitempty" json:"thirdParty,omitempty"`
	TCQF               *TCQF                        `bson:"tcqf,omitempty" json:"tcqf,omitempty"`
	JUMPTYPECONFIGV2   map[string]int32             `bson:"JUMP_TYPE_CONFIG_2,omitempty" json:"JUMP_TYPE_CONFIG_2,omitempty"`
	Updated            int                          `bson:"updated,omitempty" json:"updated,omitempty"`
	MChanlPrice        *[]MChanlPrice               `bson:"mChanlPrice,omitempty" json:"mChanlPrice,omitempty"`
	Category           *int32                       `bson:"category,omitempty" json:"category,omitempty"`
	WxAppId            *string                      `bson:"wxAppId,omitempty" json:"wxAppId,omitempty"`
	WxPath             *string                      `bson:"wxPath,omitempty" json:"wxPath,omitempty"`
	BindId             *string                      `bson:"bindId,omitempty" json:"bindId,omitempty"`
	DeepLink           *string                      `bson:"deepLink,omitempty" json:"deepLink,omitempty"`
	ApkVersion         *string                      `bson:"apkVersion,omitempty" json:"apkVersion,omitempty"`
	ApkMd5             *string                      `bson:"apkMd5,omitempty" json:"apkMd5,omitempty"`
	ApkUrl             *string                      `bson:"apkUrl,omitempty" json:"apkUrl,omitempty"`
	Platform           *int32                       `bson:"platform,omitempty" json:"platform,omitempty"`
	OsVersionMinV2     *int                         `bson:"oVersionMinV2,omitempty" json:"osVersionMinV2,omitempty"`
	OsVersionMaxV2     *int                         `bson:"osVersionMaxV2,omitempty" json:"osVersionMaxV2,omitempty"`
	StartTime          *int                         `bson:"startTime,omitempty" json:"startTime,omitempty"`
	EndTime            *int                         `bson:"endTime,omitempty" json:"endTime,omitempty"`
	NetWorkType        []int32                      `bson:"netWorkType,omitempty" json:"netWorkType,omitempty"`
	// Creative 相关的字段
	BasicCrList           *BasicCrList                        `bson:"basicCrList,omitempty" json:"basicCrList,omitempty"`
	ReadCreative          *int                                `bson:"readCreative,omitempty" json:"readCreative,omitempty"`
	CreateSrc             int                                 `bson:"createSrc,omitempty" json:"createSrc,omitempty"`
	FakeCreative          map[string]map[string]*FakeCreative `bson:"fakeCreative,omitempty" json:"fakeCreative,omitempty"`
	AlacRate              *int                                `bson:"alacRate,omitempty" json:"alacRate,omitempty"`
	AlecfcRate            *int                                `bson:"alecfcRate,omitempty" json:"alecfcRate,omitempty"`
	Mof                   *int                                `bson:"mof,omitempty" json:"mof,omitempty"`
	MCountryChanlPrice    map[string]float64                  `bson:"mCountryChanlPrice,omitempty" json:"mCountryChanlPrice,omitempty"`       // M  按 国家 + 渠道的  given price
	MCountryChanlOriPrice map[string]float64                  `bson:"mCountryChanlOriPrice,omitempty" json:"mCountryChanlOriPrice,omitempty"` // M  按 国家 + 渠道的  receive price
	NeedToNotice3s        int                                 `bson:"needToNotice3s,omitempty" json:"needToNotice3s,omitempty"`               // 判断是否需要通知3s

}

func (c *CampaignInfo) IsSSPlatform() bool {
	return c.CreateSrc == CampaignSourceSSPlatform || c.CreateSrc == CampaignSourceSSAdvPlatform
}

type FakeCreative struct {
	Id     string `bson:"id,omitempty" json:"id,omitempty"`
	Rate   int    `bson:"rate,omitempty" json:"rate,omitempty"`
	Name   string `bson:"name,omitempty" json:"name,omitempty"`
	AdType string `bson:"adType,omitempty" json:"adType,omitempty"`
}

type AppPostList struct {
	Include []string `bson:"include,omitempty" json:"include"`
	Exclude []string `bson:"exclude,omitempty" json:"exclude"`
}

type JumpParam struct {
	B2t              *int32 `bson:"b2t,omitempty" json:"b2t"`
	B2tStatus        *int32 `bson:"b2tStatus,omitempty" json:"b2tStatus"`
	NoDeviceId       *int32 `bson:"noDeviceId,omitempty" json:"noDeviceId"`
	NoDeviceIdStatus *int32 `bson:"noDeviceIdStatus,omitempty" json:"noDeviceIdStatus"`
}

type ReduceRuleItem struct {
	Priority *int32 `bson:"priority,omitempty" json:"priority"`
	Install  *int32 `bson:"install,omitempty" json:"install"`
	Status   *int32 `bson:"status,omitempty" json:"status"`
	Start    *int64 `bson:"start,omitempty" json:"start"`
}

type LoopBack struct {
	Domain *string `bson:"domain,omitempty" json:"domain,omitempty"`
	Key    *string `bson:"key,omitempty" json:"key,omitempty"`
	Value  *string `bson:"value,omitempty" json:"value,omitempty"`
	Rate   *int32  `bson:"rate,omitempty" json:"rate,omitempty"`
}

type EndCard struct {
	Urls             *[]EndCardUrls          `bson:"urls,omitempty" json:"urls,omitempty"`
	Status           *int32                  `bson:"status,omitempty" json:"status,omitempty"`
	Orientation      *int32                  `bson:"orientation,omitempty" json:"orientation,omitempty"`
	VideoTemplateUrl *[]VideoTemplateUrlItem `bson:"videoTemplateUrl,omitempty" json:"videoTemplateUrl,omitempty"`
	EndcardProtocal  *int                    `bson:"endcardProtocol,omitempty" json:"endcardProtocol,omitempty"`
	EndcardRate      *map[string]int         `bson:"endcardRate,omitempty" json:"endcardRate,omitempty"`
	EndcardType      *int32                  `bson:"endcardType,omitempty" json:"endcardType,omitempty"`
}

type EndcardItem struct {
	Url             string `bson:"url,omitempty" json:"url,omitempty"`
	UrlV2           string `bson:"url_v2,omitempty" json:"url_v2,omitempty"`
	Orientation     int32  `bson:"orientation,omitempty" json:"orientation,omitempty"`
	ID              int32  `bson:"id,omitempty" json:"id,omitempty"`
	EndcardProtocal int
	EndcardRate     map[string]int
}

type VideoTemplateUrlItem struct {
	ID           *int32  `bson:"id,omitempty" json:"id,omitempty"`
	URL          *string `bson:"url,omitempty" json:"url,omitempty"`
	URLZip       *string `bson:"url_zip,omitempty" json:"url_zip,omitempty"`
	Weight       *int32  `bson:"weight,omitempty" json:"weight,omitempty"`
	PausedURL    *string `bson:"paused_url,omitempty" json:"paused_url,omitempty"`
	PausedURLZip *string `bson:"paused_url_zip,omitempty" json:"paused_url_zip,omitempty"`
}

type EndCardUrls struct {
	Id     *int32  `bson:"id,omitempty" json:"id,omitempty"`
	Url    *string `bson:"url,omitempty" json:"url,omitempty"`
	Weight *int32  `bson:"weight,omitempty" json:"weight,omitempty"`
	UrlV2  *string `bson:"url_v2,omitempty" json:"url_v2,omitempty"`
}

type AdvImp struct {
	Sec *int32  `bson:"sec,omitempty" json:"sec,omitempty"`
	Url *string `bson:"url,omitempty" json:"url,omitempty"`
}

type Creative struct {
	CampaignId     *int64      `bson:"campaignId,omitempty" json:"campaignId,omitempty"`
	CreativeId     *int64      `bson:"creativeId,omitempty" json:"creativeId,omitempty"`
	Lang           *int32      `bson:"lang,omitempty" json:"lang,omitempty"`
	Type           *int32      `bson:"type,omitempty" json:"type,omitempty"`
	Width          *int32      `bson:"width,omitempty" json:"width,omitempty"`
	Height         *int32      `bson:"height,omitempty" json:"height,omitempty"`
	ImageSize      *string     `bson:"imageSize,omitempty" json:"imageSize,omitempty"`
	ImageSizeId    *int32      `bson:"imageSizeId,omitempty" json:"imageSizeId,omitempty"`
	Name           *string     `bson:"name,omitempty" json:"name,omitempty"`
	TextVideo      interface{} `bson:"textVideo,omitempty" json:"textVideo,omitempty"`
	VideoUrlEncode *string     `bson:"videoUrlEncode,omitempty" json:"videoUrlEncode,omitempty"`
	VideoUrl       *string     `bson:"videoUrl,omitempty" json:"videoUrl,omitempty"`
	ImageUrl       *string     `bson:"imageUrl,omitempty" json:"imageUrl,omitempty"`
	Comment        *string     `bson:"comment,omitempty" json:"comment,omitempty"`
	CreativeCta    *string     `bson:"creativeCta,omitempty" json:"creativeCta,omitempty"`
	Status         *int32      `bson:"status,omitempty" json:"status,omitempty"`
	Tag            *int32      `bson:"tag,omitempty" json:"tag,omitempty"`
	Created        *int64      `bson:"created,omitempty" json:"created,omitempty"`
	ResourceType   *int32      `bson:"resourceType,omitempty" json:"resourceType,omitempty"`
	Mime           *[]string   `bson:"mime,omitempty" json:"mime,omitempty"`
	Attribute      *[]int32    `bson:"attribute,omitempty" json:"attribute,omitempty"`
	TemplateType   interface{} `bson:"templateType,omitempty" json:"templateType,omitempty"`
	TagCode        *string     `bson:"tagCode,omitempty" json:"tagCode,omitempty"`
	ShowType       *int32      `bson:"showType,omitempty" json:"showType,omitempty"`
}

type TCQF struct {
	SubIds map[string]SubInfo `bson:"subIds,omitempty" json:"subIds"`
}

type BtV2 struct {
	SubIds  interface{}        `bson:"subIds,omitempty" json:"subIds"`
	BtClass map[string]BtClass `bson:"btClass,omitempty" json:"btClass"`
}

type BtClass struct {
	Percent   float64 `bson:"percent,omitempty" json:"percent"`
	CapMargin int32   `bson:"capMargin,omitempty" json:"capMargin"`
	Status    int32   `bson:"status,omitempty" json:"status"`
}

type BtV4 struct {
	SubIds  map[string]SubInfo `bson:"subIds,omitempty" json:"subIds"`
	BtClass map[string]BtClass `bson:"btClass,omitempty" json:"btClass"`
}

type SubInfo struct {
	Rate        int                   `bson:"rate,omitempty" json:"rate"`
	PackageName string                `bson:"packageName,omitempty" json:"packageName"`
	DspSubIds   map[string]DspSubInfo `bson:"dspSubIds,omitempty" json:"dspSubIds"`
}

type SubInfoe struct {
	Rate        int     `bson:"rate,omitempty" json:"rate"`
	PackageName string  `bson:"packageName,omitempty" json:"packageName"`
	DspSubIds   []int32 `bson:"dspSubIds,omitempty" json:"dspSubIds"`
}

type DspSubInfo struct {
	Rate        int    `bson:"rate,omitempty" json:"rate"`
	PackageName string `bson:"packageName,omitempty" json:"packageName"`
}

type ConfigVBA struct {
	UseVBA       int `bson:"useVBA,omitempty" json:"useVBA"`
	FrequencyCap int `bson:"frequencyCap,omitempty" json:"frequencyCap"`
	Status       int `bson:"status,omitempty" json:"status"`
}

type MChanlPrice struct {
	Chanl *string  `bson:"chanl,omitempty" json:"chanl"`
	Price *float64 `bson:"price,omitempty" json:"price"`
}

type ExtCreativeNew struct {
	PlayWithoutVideo int    `json:"pwv,omitempty"` // playable_ads_without_video
	VideoEndType     int    `json:"vet,omitempty"` // VideoEndType
	TemplateGroupId  *int   `json:"t_group,omitempty"`
	EndScreenId      string `json:"es_id,omitempty"`
	IsCreativeNew    bool   `json:"is_new,omitempty"`
}

type BasicCrList struct {
	AppName   string  `bson:"401,omitempty" json:"401"`
	AppDesc   string  `bson:"402,omitempty" json:"402"`
	AppRate   float64 `bson:"403,omitempty" json:"403"`
	CtaButton string  `bson:"404,omitempty" json:"404"`
	AppIcon   string  `bson:"405,omitempty" json:"405"`
	NumRating int     `bson:"406,omitempty" json:"406"`
}
