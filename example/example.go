package main

import (
	"fmt"
	index2 "github.com/Mintegral-official/juno/index"
)

func main() {
	// 建立索引
	index := index2.NewIndex("")

	//// 查询
	//query := juno.NewQuery("")// 构建查询
	//juno.NewQuery("")
	//q := NewQuery(NewAndExpress(
	//	NewEqExpression("country", "us"),
	//	NewRangeExpression("price", 1, 20)，
	//	NewOrExpress(
	//	NewEqExpression("country", "us"),
	//	NewInExpression("packageName", "package1", "package2")
	//),
	//))
	//	searchResult := index.Search(query)
	//	fmt.Println(searchResult)
	fmt.Println(index)
	//fmt.Println(index.Build())
}
