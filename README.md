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
   1. 数值型（int64）
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
见example.go

## 设计

搜索引擎主要分为2个部分

1. 索引
2. 查询

### 一.索引

#### index接口
```go
type Index interface {
    Add(docInfo *document.DocInfo) error  // 新增文档 
    Del(docInfo *document.DocInfo) error  // 删除文档 
    GetDataType(fieldName string) document.FieldType  // 获取field类型
    Dump(filename string) error  // 将索引Dump到磁盘 TODO
    Load(filename string) error  // 从磁盘加载索引 TODO
    DebugInfo() *debug.Debug  // 调试信息
}
```

#### 文档接口
```go
type DocInfo struct {
    Id     DocId      // id
    Fields []*Field   // 属性信息
}

type Field struct {
    Name      string         // 名称
    IndexType IndexType      // 索引类型 1：倒排  2：正排  3：both
    Value     interface{}   // value值
    ValueType FieldType     // value类型
}

```

json结构
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

#### 倒排 invertIndex

```go
type InvertedIndexer struct {
    data   sync.Map         // Key： FieldName + Value（一个字符串、一个数值、也可以是一个数值范围(TODO)）value: InvertList,目前实现为skiplist
    aDebug *debug.Debug    // 调试信息  在debug模式下执行
}

// invert index 接口
type InvertedIndex interface {
    Add(fieldName string, id document.DocId) error   // 新增
    Del(fieldName string, id document.DocId) bool    // 删除
    Update(fieldName string, ids []document.DocId)   // 更新某一条invert list
    Iterator(name, value string) datastruct.Iterator // 迭代器
    Count() int                                      // 统计倒排链的个数
    DebugInfo() *debug.Debug                         // debug信息接口
}

// eg: inverted_index_impl_test.go

```

#### 正排 storageIndex

```go
// 同invert index

type StorageIndexer struct {
    data   sync.Map              // 正排分字段存储，结构为map<fieldname, <docid, value>>
    aDebug *debug.Debug
}

type StorageIndex interface {
    Get(filedName string, id document.DocId) interface{}
    Add(fieldName string, id document.DocId, value interface{}) error
    Del(fieldName string, id document.DocId) bool
    Iterator(fieldName string) datastruct.Iterator
    Count() int
    DebugInfo() *debug.Debug
}

// eg: storage_index_impl_test.go

```

### 查询

#### 查询语法

查询是类sql语法，有表达式组成（可嵌套），表达式有 and, or, not等操作

支持 =, >=, >, <=,<, !=, range(暂不支持), in, !in

查询语法支持三种格式  string,  json, go struct

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

### 1. query 接口

```go
type Query interface {
    Next() (document.DocId, error)                   
    Current() (document.DocId, error)
    GetGE(id document.DocId) (document.DocId, error)
    DebugInfo() *debug.Debug
}
```

1. termQuery
```go

type TermQuery struct {
    iterator datastruct.Iterator
    debugs   *debug.Debugs
}

// 最小的query单元，通过不同的termQuery的组合来实现查询的并集、差集、交集等操作
```

2. andQuery  并集
```go
// 每个query中都出现的值
type AndQuery struct {
    queries  []Query               // 多个不同的query
    checkers []check.Checker       // 正排的过滤接口
    curIdx   int               
    debugs   *debug.Debugs
}

sl := datastruct.NewSkipList(datastruct.DefaultMaxLevel)

sl.Add(document.DocId(1), [1]byte{})
sl.Add(document.DocId(3), [1]byte{})
sl.Add(document.DocId(6), [1]byte{})
sl.Add(document.DocId(10), [1]byte{})

sl1 := datastruct.NewSkipList(datastruct.DefaultMaxLevel)

sl1.Add(document.DocId(1), [1]byte{})
sl1.Add(document.DocId(4), [1]byte{})
sl1.Add(document.DocId(6), [1]byte{})
sl1.Add(document.DocId(9), [1]byte{})

Convey("Next1", t, func() {
        a := NewAndQuery([]Query{NewTermQuery(sl.Iterator())}, nil)
        v, e := a.Current()
        //fmt.Println(v, e)
        So(v, ShouldEqual, 1)
        So(e, ShouldBeNil)

        v, e = a.Next()
        // fmt.Println(v, e)
        So(v, ShouldEqual, 3)
        So(e, ShouldBeNil)

        v, e = a.Next()
        // fmt.Println(v, e)
        So(v, ShouldEqual, 6)
        So(e, ShouldBeNil)

        v, e = a.Next()
        // fmt.Println(v, e)
        So(v, ShouldEqual, 10)
        So(e, ShouldBeNil)
    })

// eg: and_query_test.go

```

3. orQuery 交集
```go
// 满足条件的所有值
type OrQuery struct {
    checkers []check.Checker
    h        Heap              // 底层采用二叉堆实现
    debugs   *debug.Debugs
}

sl := datastruct.NewSkipList(datastruct.DefaultMaxLevel)

sl.Add(document.DocId(1), 1)
sl.Add(document.DocId(3), 2)
sl.Add(document.DocId(6), 2)
sl.Add(document.DocId(10), 1)


Convey("Next1", t, func() {
    a := NewOrQuery([]Query{NewTermQuery(sl.Iterator())}, []check.Checker{
        check.NewChecker(sl.Iterator(), 1, operation.EQ, nil),
    })
    v, e := a.Next()
    // fmt.Println(v, e)
    So(v, ShouldEqual, 1)
    So(e, ShouldBeNil)
    v, e = a.Next()
    // fmt.Println(v, e)
    So(v, ShouldEqual, 10)
    So(e, ShouldBeNil)

    v, e = a.Next()
    // fmt.Println(v, e)
    So(v, ShouldEqual, 0)
    So(e, ShouldNotBeNil)
})

// eg: or_query_test.go

```

4. notAndQuery 差集
```go
// 在第一个query中出现，在其他query中没有出现的值
type NotAndQuery struct {
    queries  []Query
    checkers []check.Checker
    curIdx   int
    debugs   *debug.Debugs
}

// eg:not_and_query_test.go
```

5. sqlQuery 类sql语法查询 -  将符合条件的字符串转换成query语句进行查询
```go
type SqlQuery struct {
    Node       *datastruct.TreeNode   //构建二叉树
    Stack      *datastruct.Stack      // 在后续遍历时用过迭代的方式，所以选用stack
    Expression *Expression            // 表达式解析
    e          operation.Operation    // 是否需要自定义的表达式
    transfer   bool                   // 是否需要转换  campaign <-> condition Operation对象中的value既可以是campaign 也可以是condition
}
func NewSqlQuery(str string, e operation.Operation, transfer bool) *SqlQuery {}  // 创建sql query对象
func (sq *SqlQuery) LRD(idx *index.Indexer) Query {}  // 构建query查询表达式
```

#### 2.check过滤接口

1. = != > < >= <=
```go
func NewChecker(si datastruct.Iterator, value interface{}, op operation.OP, e operation.Operation, transfer bool) *CheckerImpl {}
```
si : 要进行操作的storageIdx
value : 过滤条件
op : 操作符
e : 如有特殊的过滤逻辑，可自定义Operation接口实现
transfer: campaign <-> condition Operation对象中的value既可以是campaign 也可以是condition
2. in
```go
func NewInChecker(si datastruct.Iterator, value interface{}, e operation.Operation, transfer bool) *InChecker {}
```
si : 要进行操作的storageIdx
value : 过滤条件 只支持[]int []int32 []int64 []float32 []float64 []string 其他类型可以通过自定义Operation实现
e : 支持复杂的In操作  eg：两个slice的包含、交集关系等
transfer: campaign <-> condition Operation对象中的value既可以是campaign 也可以是condition
3. !in 和in操作类似
4. and or
```go
func NewOrChecker(c []Checker) *OrChecker {}
func NewAndChecker(c []Checker) *AndChecker {}
```
and和or的操作 能组合多个不同的check条件进行过滤操作

## 自定义（只针对正排）
= != > < >= <= in !in的操作
```go

type operation struct {
    value interface{}
}


func (o *operation) Equal(value interface{}) bool {
    // your logic
    return true
}

func (o *operation) Less(value interface{}) bool {
    // your logic
    return true
}

func (o *operation) In(value interface{}) bool {
    // your logic
    return true
}

func (o *operation) SetValue(value interface{}) {
    o.value = value
}
```
实现Operation接口（operastion.go）可以自定义相关的操作符的操作

#### Debug 在debug模式下执行
1. debug接口
```go
// CurNum, NextNum, GetNum 表示函数调用次数，后续性能测试使用
type Debugs struct {
	DebugInfo *Debug // debug info
	CurNum    int    // Current() transfer times
	NextNum   int    // Next() transfer times
	GetNum    int    // GetGE() transfer times
}
type Debug struct {
	Name string   `json:"name"`
	Msg  []string `json:"msg"`
	Node []*Debug `json:"node"`
}
// name 表示的是某个query或者check的名字。例如：NewAndQuery 对应的名字是: AndQuery
// Msg 一个string切片，存放的每一个元素表示的是过滤信息，比如某个id没有在某条query里面出现，check条件等
// Node节点， 对应query之间的嵌套关系，AndQuery里面的[]Query放在node里面，递归保存query
//eg: NewAndQuery([]Query{NewTermQuery{}, NewAndQuery(NewTermQuery{},NewTermQuery)},
//               []checker{NewChecker{},NewChecker{},NewAndChecker(NewChecker{})})
{
    "name":"AndQuery",
    // 这个表示3这个id在check中的条件的情况，
    // 第一个为true,说明在第一个check中是不会被过滤的
    // 第二个为false,说明在第二个check中会被过滤掉
    "msg":["3 check: [true,false,{\"name\":\"AndCheck\",\"msg\":[\"true\"],\"node\":null}]"],
    "node":[
        {
            "name":"TermQuery",
            "msg":[],
            "node":null
        },
        {
            "name":"AndQuery",
            "msg":[],
            "node":[
                {
                    "name":"TermQuery",
                    "msg":[],
                    "node":null
                },
                {
                    "name":"TermQuery",
                    "msg":[],
                    "node":null
                }
            ]
        }
    ]
}

// Query中新增两个接口 ：
// 1.SetDebug(isDebug ...int): query调用接口,query.SetDebug(1) 表示开启debug模式
// 2.UnsetDebug(): query.UnsetDebug() 关闭debug模式
```

```go
// invertIndex 新增方法，通过id获取对应的fieldName,
func (i *InvertedIndexer) GetValueById(id document.DocId) []string {}

var doc1 = &document.DocInfo{
	Id: 0,
	Fields: []*document.Field{
		{
			Name:      "field1",
			IndexType: 1,
			Value:     1,
			ValueType: document.IntFieldType,
		},
		{
			Name:      "field2",
			IndexType: 0,
			Value:     "2",
			ValueType: document.StringFieldType,
		},
		{
			Name:      "field1",
			IndexType: 2,
			Value:     "1",
			ValueType: document.StringFieldType,
		},
	},
}

var doc2 = &document.DocInfo{
	Id: 1,
	Fields: []*document.Field{
		{
			Name:      "field1",
			IndexType: 0,
			Value:     "1",
			ValueType: document.StringFieldType,
		},
		{
			Name:      "field2",
			IndexType: 1,
			Value:     "2",
			ValueType: document.StringFieldType,
		},
		{
			Name:      "field1",
			IndexType: 0,
			Value:     "2",
			ValueType: document.StringFieldType,
		},
	},
}

// eg：
// docId = 0, 则返回[field2_2, field1_1]
// docId = 1, 则返回[field1_1, field1_2]
// 调用方法
var idx Indexer
invertIdx := idx.GetInvertIndex()
invertIdx.GetValueById(docId)


	    ss := index.NewIndex("")
	    s1 := ss.GetInvertedIndex()
	    s2 := ss.GetStorageIndex()
		q := query.NewAndQuery([]query.Query{
			query.NewTermQuery(s1.Iterator("fieldName", "1")),
			query.NewTermQuery(s1.Iterator("fieldName", "2")),
			query.NewTermQuery(s1.Iterator("fieldName", "3")),
		}, []check.Checker{
			check.NewChecker(s2.Iterator("fieldName"), 3, operation.EQ, nil, false),
			check.NewChecker(s2.Iterator("fieldName"), 3, operation.EQ, nil, false),
			check.NewAndChecker([]check.Checker{
				check.NewChecker(s2.Iterator("fieldName"), 3, operation.EQ, nil, false),
				check.NewChecker(s2.Iterator("fieldName"), 3, operation.EQ, nil, false),
			}),
		})

		q.SetDebug(1) // 设置debug调试信息

		se := NewSearcher()

        //  ss是index索引  q是query语句
		fmt.Println(se.Debug(ss, q).String()) // debug查询
		// se.DebugInfo(ss, q, ids)  ids：指定查找的id列表
		{
		    "node":{
		    "10":[
		        // 不一定所有的id都会走check
		        ["field:fieldName_1","reason: found id"],
		        ["field:fieldName_2","reason: not found"],
		        ["fieldName_1"]
		    ],
		    "3":[
		        // 对应的是check的结果，包含：
		        // fieldname, 正排的那个值 storageIdx.Iterator("Name"), Name
		        // condition传进去的value，
		        // 操作符(=, !=, >, <)等操作，
		        // operation：判断这个operation是自定义的还是使用默认的，transfer暂不用考虑
		        ["and check result: false","FieldName: fieldName\tvalue: 3\tOP: =\tdefined operation: false\ttransfer: false\tcheck result: true","FieldName: fieldName\tvalue: 3\tOP: =\tdefined operation: false\ttransfer: false\tcheck result: true","{\"node\":{\"3\":[[\"and check result: false\",\"FieldName: fieldName\\tvalue: 3\\tOP: =\\tdefined operation: false\\ttransfer: false\\tis checked: true\",\"FieldName: fieldName\\tvalue: 3\\tOP: =\\tdefined operation: false\\ttransfer: false\\tis checked: true\"]]}}\tcheck result: false"],
		        // 中间表示的是query的过滤原因，在field条件下，这个id是否存在
		        ["field:fieldName_1","reason: found id"],
		        ["field:fieldName_2","reason: found id"],
		        ["field:fieldName_3","reason: found id"],
		        // 表示的是倒排建立索引的每一个id对应的全部field值，建立索引时候的name_value组合
		        ["fieldName_5","fieldName_6","fieldName_1","fieldName_2","fieldName_3","fieldName_4"]
		    ],
		    "4":[
		        ["and check result: false","FieldName: fieldName\tvalue: 3\tOP: =\tdefined operation: false\ttransfer: false\tcheck result: true","FieldName: fieldName\tvalue: 3\tOP: =\tdefined operation: false\ttransfer: false\tcheck result: true","{\"node\":{\"3\":[[\"and check result: false\",\"FieldName: fieldName\\tvalue: 3\\tOP: =\\tdefined operation: false\\ttransfer: false\\tis checked: true\",\"FieldName: fieldName\\tvalue: 3\\tOP: =\\tdefined operation: false\\ttransfer: false\\tis checked: true\"],[\"and check result: false\",\"FieldName: fieldName\\tvalue: 3\\tOP: =\\tdefined operation: false\\ttransfer: false\\tis checked: true\",\"FieldName: fieldName\\tvalue: 3\\tOP: =\\tdefined operation: false\\ttransfer: false\\tis checked: true\"]],\"4\":[[\"and check result: false\",\"FieldName: fieldName\\tvalue: 3\\tOP: =\\tdefined operation: false\\ttransfer: false\\tis checked: true\",\"FieldName: fieldName\\tvalue: 3\\tOP: =\\tdefined operation: false\\ttransfer: false\\tis checked: true\"]]}}\tcheck result: false"],
		        ["field:fieldName_1","reason: found id"],
		        ["field:fieldName_2","reason: found id"],
		        ["field:fieldName_3","reason: found id"],
		        ["fieldName_6","fieldName_1","fieldName_2","fieldName_3","fieldName_4","fieldName_5"]
		    ],
		    "6":[
		        ["field:fieldName_1","reason: found id"],
		        ["field:fieldName_2","reason: found id"],
		        ["field:fieldName_3","reason: not found"],
		        ["fieldName_1","fieldName_2","fieldName_4","fieldName_6"]
		    ]}
        }
```

