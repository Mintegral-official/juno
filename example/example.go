package main

import "fmt"

func main()  {
	// 建立索引
	index := Index.NewIndex("")

	// 查询
	query := // 构建查询
	q := NewQuery(NewAndExpress(
		NewEqExpression("country", "us"),
		NewRangeExpression("price", 1, 20)，
		NewOrExpress(
		NewEqExpression("country", "us"),
		NewInExpression("packageName", "package1", "package2")
	),
))
	searchResult := index.Search(query)
	fmt.Println(searchResult)
}


