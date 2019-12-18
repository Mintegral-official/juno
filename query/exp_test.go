package query

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestString2Strings(t *testing.T) {
	s := "country= CN &   (a =1 | ( b = 1 & a!=0)) | (c @ [1,2,3] & d # [2,4,5])"
	e := NewExpression(s)
	Convey("string2string", t, func() {
		So(e.string2Strings(), ShouldNotBeNil)
		So(e.ToPostfix(e.string2Strings()), ShouldNotBeNil)
		a := e.ToPostfix(e.string2Strings())
		So(len(a), ShouldEqual, 11)
		So(e.GetValue(), ShouldNotBeNil)
	})
	//fmt.Println(e.string2Strings())
	//fmt.Println(e.ToPostfix(e.string2Strings()))
	//a := e.ToPostfix(e.string2Strings())
	//for i, v := range a {
	//	fmt.Println(i, "->", v)
	//}
}
