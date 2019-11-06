package index

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestNewSimpleStorageIndex(t *testing.T) {
    Convey("Get", t, func() {
		s := NewSimpleStorageIndex()
    	So(s.Get("fieldName", 1), ShouldBeNil)
    	So(s.Add("fieldName", 1, 1), ShouldBeNil)
    	So(s.Del("fieldName", 1), ShouldBeTrue)
    	So(s.Iterator("fieldName"), ShouldBeNil)
	})
}
