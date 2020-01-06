package query

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestString2Strings(t *testing.T) {
	s := "country= CN and   (a =1 or ( b = 1 & a!=0)) or (c in [1,2,3] and d !in [2,4,5])"
	e := NewExpression(s)
	Convey("string2string", t, func() {
		So(e.string2Strings(), ShouldNotBeNil)
		So(e.ToPostfix(e.string2Strings()), ShouldNotBeNil)
		a := e.ToPostfix(e.string2Strings())
		So(len(a), ShouldEqual, 11)
		So(e.GetValue(), ShouldNotBeNil)
		So(e.string2Strings(), ShouldNotBeNil)
		//for _, v := range e.string2Strings() {
		//	fmt.Println("***"+v+"***")
		//}
	})
	//fmt.Println(e.string2Strings())
	//fmt.Println(e.ToPostfix(e.string2Strings()))
	//a := e.ToPostfix(e.string2Strings())
	//for i, v := range a {
	//	fmt.Println(i, "->", v)
	//}
}
