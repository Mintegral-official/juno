# Juno(朱诺)
朱诺号木星探测器是目前人类是制造出最快的宇宙飞行器。
这里，朱诺是一个通用的易用的高性能的内存型广告检索引擎

## 目标

1. 通用性： 能试用广告检索的大部分情况
2. 易用性： 可以极低的代价从0搭建搜索引擎
3. 高性能： 本身搜索性能20ms内，单机QPS>1-2K
4. 插件化，可扩展： 检索各模块都是接口的形式，可以根据需求轻松定制

## 示例

### 通过代码构建索引

```go
func main() {
  // build index
	idx := index.NewIndex("default")
	_ = idx.Add(&document.DocInfo{
		Id: 1,
		Fields: []*document.Field{
			{Name: "field1", IndexType: document.InvertedIndexType, Value: int64(1), ValueType: document.IntFieldType},
			{Name: "field2", IndexType: document.InvertedIndexType, Value: "abc", ValueType: document.StringFieldType},
		},
	})

  // search
	s := search.NewSearcher()
	s.Search(idx, query.NewTermQuery(idx.GetInvertedIndex().Iterator("field1", "1")))
	fmt.Println(s.Docs)
}
```



###直接从mongo读数据建立索引

需要用户实现 MongoParser接口， 将mongo的一条记录转化成DocInfo

```go
type ParserResult struct {
	DataMod DataMod
	Value   *document.DocInfo
}

type MongoParser interface {
	Parse([]byte, interface{}) *ParserResult
}

// example2 build index with mongo
func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// build index
	b, e := builder.NewMongoIndexBuilder(&builder.MongoIndexManagerOps{
		URI:            "mongodb://13.250.108.190:27017",
		IncInterval:    5,
		BaseInterval:   120,
		IncParser:      &CampaignParser{},
		BaseParser:     &CampaignParser{},
		BaseQuery:      bson.M{"status": 1},
		IncQuery:       bson.M{"updated": bson.M{"$gte": time.Now().Unix() - 5, "$lte": time.Now().Unix()}},
		DB:             "new_adn",
		Collection:     "campaign",
		ConnectTimeout: 10000,
		ReadTimeout:    20000,
		UserData:       &UserData{},
		Logger:         logrus.New(),
		OnBeforeInc: func(userData interface{}) interface{} {
			ud, ok := userData.(*UserData)
			if !ok {
				return nil
			}
			incQuery := bson.M{"updated": bson.M{"$gte": ud.upTime - 5, "$lte": time.Now().Unix()}}
			return incQuery
		},
	})
	if e != nil {
		fmt.Println(e)
		return
	}
	if e := b.Build(ctx, "indexName"); e != nil {
		fmt.Println("build error", e.Error())
	}

	tIndex := b.GetIndex()

```



