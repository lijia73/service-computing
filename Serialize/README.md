# serialize
将多种类型数据包括结构数据格式化为json字符流
- int
- uint
- 字符串
- 结构体
- 数组/切片

## 使用方法
### 安装
1. go get -u gitee.com/li-jia666/serialize
2. import "gitee.com/li-jia666/serialize"

``func JsonMarshal(v interface{}) ([]byte, error)``

### 简单的使用案例

```go
package main

import (
	"fmt"
	"os"
	"gitee.com/li-jia666/serialize"
)

type ColorGroup struct {
	ID     int	`mytag:"color_id"`
	Name   string	`mytag:"color_name"`
	Colors []string	`mytag:"colors"`
}

func main() {
    group := ColorGroup{
        ID:     1,
        Name:   "Reds",
        Colors: []string{"Crimson", "Red", "Ruby", "Maroon"},
    }
    b, err := serialize.JsonMarshal(group)
    if err != nil {
        fmt.Println("error:", err)
    }
    os.Stdout.Write(b)
}
```