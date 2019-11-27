package model

import (
	"github.com/Mintegral-official/juno/conf"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestNewMongo(t *testing.T) {
	cfg := &conf.MongoCfg{
		URI:            "mongodb://localhost:27017",
		DB:             "new_adn",
		Collection:     "campaign",
		ConnectTimeout: 10000,
		ReadTimeout:    20000,
	}
	Convey("mongo", t, func() {
		m, e := NewMongo(cfg)
		if e != nil {
			So(m, ShouldBeNil)
		}
		So(m, ShouldBeNil)
		So(e, ShouldNotBeNil)
		//f, e := m.Find()
		//So(f, ShouldBeNil)
		//So(e, ShouldBeNil)

	})

}
