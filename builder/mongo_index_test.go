package builder

import (
	. "github.com/smartystreets/goconvey/convey"
	"go.mongodb.org/mongo-driver/bson"
	"testing"
	"time"
)

type campaignParser struct {
}

type campaignInfo struct {
	CampaignId   int64  `bson:"campaignId,omitempty" json:"campaignId,omitempty"`
	AdvertiserId *int32 `bson:"advertiserId,omitempty" json:"advertiserId,omitempty"`
	Platform     *int32 `bson:"platform,omitempty" json:"platform,omitempty"`
}

func (c *campaignParser) Parse(bytes []byte) (*ParserResult, error) {
	campaign := &campaignInfo{}
	if err := bson.Unmarshal(bytes, &campaign); err != nil {
		return nil, err
	}
	return nil, nil
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
