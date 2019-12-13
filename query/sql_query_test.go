package query

import (
	"github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestString2Strings(t *testing.T) {
	s := "country= CN &   (a =1 | ( b = 1 & a!=0)) | (c @ [1,2,3] & d # [2,4,5])"
	convey.Convey("string2string", t, func() {
		convey.So(String2Strings(s), convey.ShouldNotBeNil)
		convey.So(ToPostfix(String2Strings(s)), convey.ShouldNotBeNil)
		a := ToPostfix(String2Strings(s))
		convey.So(len(a), convey.ShouldEqual, 11)
	})
	//fmt.Println(String2Strings(s))
	//fmt.Println(ToPostfix(String2Strings(s)))
	//a := ToPostfix(String2Strings(s))
	//for i, v := range a {
	//	fmt.Println(i, "->", v)
	//}
}
