package index

import (
	"github.com/Mintegral-official/juno/conf"
	"github.com/Mintegral-official/juno/document"
	"github.com/Mintegral-official/juno/helpers"
	"github.com/Mintegral-official/juno/model"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

var doc1 = &document.DocInfo{
	Id: 0,
	Fields: []*document.Field{
		{
			Name:      "field1",
			IndexType: 1,
			Value:     nil,
		},
		{
			Name:      "field2",
			IndexType: 0,
			Value:     nil,
		},
		{
			Name:      "field1",
			IndexType: 2,
			Value:     nil,
		},
	},
}

var doc2 = &document.DocInfo{
	Id: 1,
	Fields: []*document.Field{
		{
			Name:      "field1",
			IndexType: 0,
			Value:     nil,
		},
		{
			Name:      "field2",
			IndexType: 1,
			Value:     nil,
		},
		{
			Name:      "field1",
			IndexType: 0,
			Value:     nil,
		},
	},
}

var doc3 = &document.DocInfo{
	Id: 2,
	Fields: []*document.Field{
		{
			Name:      "field1",
			IndexType: 0,
			Value:     nil,
		},
		{
			Name:      "field2",
			IndexType: 0,
			Value:     nil,
		},
		{
			Name:      "field1",
			IndexType: 1,
			Value:     nil,
		},
	},
}

func TestNewIndex(t *testing.T) {
	Convey("NewIndex", t, func() {
		So(NewIndex("index"), ShouldNotBeNil)
	})

	Convey("Add", t, func() {
		index := NewIndex("index")
		So(index.Add(nil), ShouldEqual, helpers.DocumentError)
		So(index.Add(doc1), ShouldBeNil)
		So(index.Add(doc2), ShouldBeNil)
		So(index.Add(doc3), ShouldBeNil)
		if1 := index.invertedIndex.Iterator("field1")
		c := 0
		for if1.HasNext() {
			if if1.Next() != nil {
				c++
			}
		}
		So(c, ShouldEqual, 3)

		if2 := index.invertedIndex.Iterator("field2")
		c = 0
		for if2.HasNext() {
			if if2.Next() != nil {
				c++
			}
		}
		So(c, ShouldEqual, 2)
		sf1 := index.storageIndex.Iterator("field1")
		c = 0
		for sf1.HasNext() {
			if sf1.Next() != nil {
				c++
			}
		}
		So(c, ShouldEqual, 2)
		sf2 := index.storageIndex.Iterator("field2")
		c = 0
		for sf2.HasNext() {
			if sf2.Next() != nil {
				c++
			}
		}
		So(c, ShouldEqual, 1)
		So(index.Del(doc1), ShouldBeNil)
	})
}

func TestInterface(t *testing.T) {
	var a, b interface{}
	i := int64(1)
	a = &i
	b = &i

	if b == a {
		// fmt.Println("xxxxxxxxxxxxxxx")
	}

}

func TestNewIndex2(t *testing.T) {
	cfg := &conf.MongoCfg{
		URI:            "mongodb://localhost:27017",
		DB:             "new_adn",
		Collection:     "campaign",
		ConnectTimeout: 10000,
		ReadTimeout:    20000,
	}

	Convey("mongoIndex", t, func() {
		mon, err := model.NewMongo(cfg)
		pkgNames := make([]string, 10)
		osVersionCodeV2 := ""
		timestamp := -1
		category := -1
		networkId := -1
		directMarket := -1
		networkType := -1
		So(err, ShouldNotBeNil)
		// fmt.Println(mon)
		r, e := mon.Find()
		So(e, ShouldNotBeNil)
		//fmt.Println(e)
		for i := 0; i < len(r); i++ {
			if !helpers.In(r[i].PackageName, pkgNames) {
				continue
			}
			if !r[i].IsSSPlatform() {
				continue
			}
			if int(*r[i].AdvertiserId) == 919 || int(*r[i].AdvertiserId) == 976 {
				continue
			}
			if !helpers.In(osVersionCodeV2, r[i].OsVersionMinV2, r[i].OsVersionMaxV2) {
				continue
			}
			if !helpers.In(timestamp, r[i].StartTime, r[i].EndTime) {
				continue
			}
			if !helpers.Equal(category, "UNKNOWN") || !helpers.Equal(category, int(*r[i].Category)) {
				continue
			}
			if !helpers.Equal(networkId, r[i].Network) {
				continue
			}
			if !helpers.Equal(directMarket, "UNKNOWN") || !helpers.Equal(directMarket, "NO_LIMIT") ||
				!(helpers.Equal(directMarket, "NO_APK") && (helpers.Equal(r[i].CampaignType, "GooglePlay") ||
					helpers.Equal(r[i].CampaignType, "OTHER"))) || !(helpers.Equal(directMarket, "ONLY_APK") &&
				helpers.Equal(r[i].CampaignType, "APK")) || !(helpers.Equal(directMarket, "ONLY_GP") &&
				helpers.Equal(r[i].CampaignType, "GooglePlay")) {
				continue
			}
			if !helpers.In("UNKNOWN", r[i].NetWorkType) || !helpers.In(networkType, r[i].NetWorkType) ||
				!(helpers.In("RECALL_IF_UNKNOWN", r[i].NetWorkType) && helpers.Equal(networkType, "UNKNOWN")) {
				continue
			}

		}
	})

}
