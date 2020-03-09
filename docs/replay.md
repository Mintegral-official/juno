# Juno过滤原因回放

juno检索引擎支持过滤原因回放

## 将query可序列化成json，记录日志

```go
// 序列化接口
func (q* Query) Marshal()  ([]byte, error)
```

json结构：

```json
{
    "and": {
        "not": [
            {
                "=": [
                    "Os",
                    "2"
                ]
            },
            {
                "or": [
                    {
                        "or": [
                            {
                                "=": [
                                    "PackageName",
                                    "ru.bestfeeds.VK-Feed"
                                ]
                            }
                        ]
                    },
                    {
                        "=": [
                            "DeviceAndIpuaRetarget",
                            "1"
                        ]
                    }
                ]
            }
        ]
    },
    "and_check": {
        "check": [
            "OsVersionMin",
            9030300,
            2,
            0,
            false,
            "<="
        ]
    }
}
```





## 调用Replay接口， 获取文档过滤原因

```go
func Replay(index *index.Index, jsonQuery []byte, ids []document.DocId)
```

返回文档过滤原因

```json
       
{
    "12345":{
        "And":{
            "result":true,
            "node":[
                {
                    "term":{
                        "key":"DeviceAndIpuaRetarget",
                        "value":"2",
                        "result":true
                    }
                }
            ]
        },
        "and_check":{
            "field":"NeedMobileCode",
            "value":0,
            "op":">",
            "result":true,
            "queryValue":0
        }
    },
    "123467":{

    }
}
```

