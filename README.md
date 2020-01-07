# Juno(朱诺)
朱诺号木星探测器是目前人类是制造出最快的宇宙飞行器。
这里，朱诺是一个通用的易用的高性能的内存型广告检索引擎

## 目标

1. 通用性： 能试用广告检索的大部分情况
2. 易用性： 可以极低的代价从0搭建搜索引擎
3. 高性能： 本身搜索性能20ms内，单机QPS>1-2K
4. 插件化，可扩展： 检索各模块都是接口的形式，可以根据需求轻松定制

## 主要特性

1. 支持倒排索引
   1. 数值型（int）
   2. 字符串型 (string)
   3. 切片 ([]int64 []string)
2. 正排索引
   1. 数值型（int, double）
   2. 字符串型
   3. set集合
   4. List
   5. KV
3. 查询支持多索引查询、布尔查询、范围查询、集合查询

## 示例

```go
// 建立索引
// build index
b, e := builder.NewMongoIndexBuilder(&builder.MongoIndexManagerOps{
    URI:            "mongodb://127.0.0.1:27017",
    IncInterval:    5,
    BaseInterval:   7200,
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
b.Build(ctx, "indexName")
idx := b.GetIndex()

// invert list
invertIdx := idx.GetInvertedIndex()

// storage list
storageIdx := idx.GetStorageIndex()
	

// 正排
if1 := idx.GetStorageIndex().Iterator("AdvertiserId")
if2 := idx.GetStorageIndex().Iterator("Platform")
if3 := idx.GetStorageIndex().Iterator("Price")

//倒排
if1 := idx.GetInvertIndex().Iterator("Platform", "1")
if2 := idx.GetInvertIndex().Iterator("AdvertiserId", "100")

// query查询
q := query.NewOrQuery([]query.Query{
    query.NewOrQuery([]query.Query{
        query.NewTermQuery(invertIdx.Iterator("Platform", "1")),
    }, nil),
    query.NewOrQuery([]query.Query{
        query.NewTermQuery(invertIdx.Iterator("AdvertiserId", "457")),
    }, nil),
    query.NewOrQuery([]query.Query{
        query.NewTermQuery(invertIdx.Iterator("DeviceTypeV2", "4")),
        query.NewTermQuery(invertIdx.Iterator("DeviceTypeV2", "5")),
    }, nil),
    query.NewAndQuery([]query.Query{
        query.NewAndQuery([]query.Query{
            query.NewTermQuery(storageIdx.Iterator("Price")),
        }, []check.Checker{
            check.NewInChecker(storageIdx.Iterator("Price"),
                2.3, 1.4, 3.65, 2.46, 2.5),
        }),
        query.NewAndQuery([]query.Query{
            query.NewTermQuery(storageIdx.Iterator("AdvertiserId")),
        }, []check.Checker{
            check.NewNotChecker(storageIdx.Iterator("AdvertiserId"), int64(647), int64(658), int64(670)),
        })}, nil)},
    nil,
)
// 查询
r1 := search.NewSearcher()
r1.Search(tIndex, q)

//sql 查询  or and in !in (只支持小写)
a := "AdvertiserId=457 or Platform=1 or (Price in [2.3, 1.4, 3.65, 2.46, 2.5] and AdvertiserId !in [647, 658, 670])"
sq := query.NewSqlQuery(a)
m := sq.LRD(tIndex)
r2 := search.NewSearcher()
r2.Search(tIndex, m)

```

## 设计

搜索引擎主要分为2个部分

1. 查询
2. 索引

### 一、查询

#### 查询语法

查询是类sql语法，有表达式组成（可嵌套），表达式有 and, or, not等操作

支持 =, >=, >, <=,<, !=, range(暂不支持), in, !in

查询语法支持三种格式  string,  json, go stuct


```json
{
    "and": [
        {
            "=": {
                "field": "country",
                "value": "US"
            }
        },
        {
            "range": {
                "field": "price",
                "value": [
                    1,
                    20
                ]
            }
        },
        {
            "or": [
                {
                    "=": {
                        "field": "platform",
                        "value": "ios"
                    }
                },
                {
                    "in": {
                        "field": "packageName",
                        "value": [
                            "package1",
                            "package2"
                        ]
                    }
                }
            ]
        }
    ]
}
```

#### 查询执行过程

1. 构建查询语法树

   ![](pic/search_tree.png)

2. 执行语法树

   1. 语法树本身可以抽象成一个迭代器，迭代的过程就是对倒排链查找的过程

3. 过滤

### 二、索引

#### 索引接口

index接口：

```go
type Index interface {
	Add(docInfo *document.DocInfo) error  // 新增文档
	Del(docInfo *document.DocInfo) error  // 删除文档
	GetDataType(fieldName string) document.FieldType  // 获取field类型
	Dump(filename string) error  // 将索引Dump到磁盘
	Load(filename string) error  // 从磁盘加载索引
	DebugInfo() *debug.Debug  // 调试信息
}
```

DocInfo:  json结构

```json
{
    "Id": 12345,
    "Fields": [
        {
           "FieldName": "Field1",
           "value":"value",
           "indexType":0
        },
        {
           "FieldName": "Field2",
           "value":"value",
           "indexType":1
        },
        {
           "FieldName": "Field3",
           "value":"value",
           "indexType":2
        }
    ]
}
```

#### 索引内存结构

##### 倒排

倒排索引是一个Key, InvertList结构

- Key： FieldName + Value（一个字符串、一个数值、也可以是一个数值范围）
- InvertList： 是一个有序的集合的接口，，可以是数组、跳表、排序树等
- Value:  一个字符串、一个数值、也可以是一个数值范围

###### 倒排接口：

倒排索引可以有不同的实现方式，只要满足下面的接口，都可以称之为倒排索引

```go
// InvertList 倒排结构的接口，仅负责查询，不负责索引更新
type InvertIndex interface {
	Add(fieldName string, id document.DocId) error
	Del(fieldName string, id document.DocId) bool
	Iterator(fieldName string) datastruct.Iterator
}
```

##### 正排索引

正排分字段存储，结构为map<fieldname, <docid, value>>

###### 正排接口

```go
// 按字段存储正排信息
	Get(filedName string, id document.DocId) interface{}
	Add(fieldName string, id document.DocId, value interface{}) error
	Del(fieldName string, id document.DocId) bool
	Iterator(fieldName string) datastruct.Iterator
```



#### 索引构建

索引构建模块能方便的将数据源中的数据构建成索引，同时能感知数据源的变化，并将变化同步至索引中


索引构建模块会支持多种数据源，如文件、mongo、mysql等

