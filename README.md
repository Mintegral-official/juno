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