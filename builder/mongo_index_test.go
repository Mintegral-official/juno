package builder

import (
	. "github.com/smartystreets/goconvey/convey"
	"go.mongodb.org/mongo-driver/bson"
	"testing"
	"time"
)

type campaignParser struct {
}

type Data struct {
	upTime int64
}

type campaignInfo struct {
	CampaignId   int64  `bson:"campaignId,omitempty" json:"campaignId,omitempty"`
	AdvertiserId *int32 `bson:"advertiserId,omitempty" json:"advertiserId,omitempty"`
	Platform     *int32 `bson:"platform,omitempty" json:"platform,omitempty"`
	Uptime       int64  `bson:"updated,omitempty"`
}

func (c *campaignParser) Parse(bytes []byte, userData interface{}) *ParserResult {
	ud, ok := userData.(*Data)
	if !ok {
		return nil
	}
	campaign := &campaignInfo{}
	if err := bson.Unmarshal(bytes, &campaign); err != nil {
		return nil
	}
	if ud.upTime < campaign.Uptime {
		ud.upTime = campaign.Uptime
	}
	return nil
}

func TestNewMongoIndexBuilder(t *testing.T) {
	Convey("mongo index", t, func() {
		mib, e := NewMongoIndexBuilder(&MongoIndexManagerOps{
			URI:            "mongodb://127.0.0.1:27017",
			IncInterval:    5,
			BaseInterval:   120,
			IncParser:      &campaignParser{},
			BaseParser:     &campaignParser{},
			BaseQuery:      bson.M{"status": 1},
			IncQuery:       bson.M{"updated": bson.M{"$gt": time.Now().Unix() - int64(5*time.Second)}},
			DB:             "new_adn",
			Collection:     "campaign",
			ConnectTimeout: 10000,
			ReadTimeout:    20000,
		})
		So(mib, ShouldBeNil)
		So(e, ShouldNotBeNil)
	})
}
