package debug

import (
	"encoding/json"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func Test_Debug(t *testing.T) {
	Convey("debug", t, func() {
		debug := &Debug{
			Name: "debug1",
			Msg:  []string{"msg1", "msg2"},
			Node: []*Debug{
				{
					Name: "msg2",
					Msg:  []string{"msg3", "msg4"},
					Node: nil,
				},
			},
		}

		j, err := json.Marshal(debug)
		So(err, ShouldBeNil)
		So(j, ShouldNotBeNil)
		So(len(debug.Node), ShouldEqual, 1)
		So(len(debug.Msg), ShouldEqual, 2)

		debug.AddDebug(&Debug{
			Name: "123",
			Msg:  nil,
			Node: nil,
		})
		debug.AddDebugMsg("msg11", "msg22")
		So(len(debug.Node), ShouldEqual, 2)
		So(len(debug.Msg), ShouldEqual, 4)
		So(debug, ShouldNotBeNil)
		So(debug.String(), ShouldNotBeNil)
	})
}
