package index

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestNewSimpleStorageIndex(t *testing.T) {
	s := NewSimpleStorageIndex()
    Convey("Get", t, func() {
    	So(s.Get("filename", 1), ShouldBeNil)
	})
}
