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

   ![](../pic/search_tree.png)

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
