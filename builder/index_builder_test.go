package builder

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestNewIndexBuilder(t *testing.T) {
	Convey("new index builder", t, func() {
		ib := NewIndexBuilder(&MongoIndexManagerOps{})
		So(ib, ShouldNotBeNil)
	})

}
