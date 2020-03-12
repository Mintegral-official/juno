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



### mongo读数据建立索引

juno内置支持从mongo读取数据建立索引，支持两种模式，全量与增量

只需要用户实现 MongoParser接口， 将mongo的一条记录转化成DocInfo即可

```go
// mongo解析结果
type ParserResult struct {
	DataMod DataMod
	Value   *document.DocInfo
}

// mongo解析器
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
		IncInterval:    5,  // 增量间隔
		BaseInterval:   120,  // 全量间隔
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

        // 获取构建的索引
	tIndex := b.GetIndex()
}
```



## 查询语法

juno目前支持类sql，go struct两种查询语法，同时支持debug模式，可以获取特定文档没召回的原因

保留字段： where search index

### 类SQL查询语法 示例

> 基本query: 支持=,!=, >, >=, <, <=, in, has等操作符
>
> 1. campaignId = 1
> 2. campaignId in [1, 2, 3]
> 3. adid has [1, 2, 3]
> 4. price > 10
> 5. price < 100
> 6. campain != 5
>
> 复核query: 基本query的组合, 支持 and, or, not
> 1. campainId = 1 && price > 10
> 2. adid has [1, 2, 3] && (price > 10 || os = 1)
> 3. adid has [1, 2, 3] || (not campaignId = 5)
> 4. adid has not [1, 2, 3]
>
> 自定义函数：
> func(fieldName, query) bool
>
> 1. func1(price, 100)
> 2. campainId = 1 && price > 10 && func(price, 100) 
> 3. regex_func(fieldName, "xxx")
>
> 文档过滤原因
>
> Query: {adv = 1 && price > 10 | business1} && {adv = 1 && price > 10 && func(price, 100) | business2}  docid in [1,2]
>
> 返回结果：1: business1=true, business=false;2:business1=fasle,business2=false

### go struct 查询

```go
q := query.NewOrQuery([]query.Query{
			query.NewOrQuery([]query.Query{
				query.NewTermQuery(invertIdx.Iterator("Platform", "1")),
			}, nil),
			query.NewOrQuery([]query.Query{
				query.NewTermQuery(invertIdx.Iterator("AdvertiserId", "457")),
			}, nil),
			/* special example */
			query.NewOrQuery([]query.Query{
				query.NewTermQuery(storageIdx.Iterator("DeviceTypeV2")),
			}, []check.Checker{
				check.NewInChecker(storageIdx.Iterator("DeviceTypeV2"), devi, nil, false),
			}),
			query.NewAndQuery([]query.Query{
				query.NewAndQuery([]query.Query{
					query.NewTermQuery(storageIdx.Iterator("Price")),
				}, []check.Checker{
					check.NewInChecker(storageIdx.Iterator("Price"), pi, nil, false),
				}),
				query.NewAndQuery([]query.Query{
					query.NewTermQuery(storageIdx.Iterator("AdvertiserId")),
				}, []check.Checker{
					check.NewNotChecker(storageIdx.Iterator("AdvertiserId"), ai, nil, false),
				})}, nil)},
			nil,
		)

r := search.Search(index, q)
```

### 文档过滤原因

* [详见](./docs/replay.md)



## 性能

## 

## 未来特性

1. 多数据源构建索引
2. 索引持久化
