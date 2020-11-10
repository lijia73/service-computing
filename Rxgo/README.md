# RxGo
- Rx是一个函数库，让开发者可以利用可观察序列和LINQ风格查询操作符来编写异步和基于事件的程序，使用Rx，开发者可以用Observables表示异步数据流。
- 在[pmlpml/RxGo](https://github.com/pmlpml/rxgo)的基础上增加了一组新的操作filtering。
该库的基本组成：

`rxgo.go` 给出了基础类型、抽象定义、框架实现、Debug工具等

`generators.go` 给出了 sourceOperater 的通用实现和具体函数实现

`transforms.go` 给出了 transOperater 的通用实现和具体函数实现

`filtering.go` 给出了 filOperater 的通用实现和具体函数实现
## 使用方法
### 安装
1. go get -u gitee.com/li-jia666/rxgo
2. import "gitee.com/li-jia666/rxgo"


### 简单的使用案例

```go
package main

import (
	"fmt"
	"time"
	RxGo "gitee.com/li-jia666/rxgo"
)

func main() {
    fmt.Println("Debounce:")
    RxGo.Just(1,2,3,4,5,6).Map(func(x int) int {
		switch x {
		case 1:
			time.Sleep(0 * time.Millisecond)
		case 2:
			time.Sleep(250 * time.Millisecond)
		case 3:
			time.Sleep(300 * time.Millisecond)
		case 4:
			time.Sleep(100 * time.Millisecond)
		case 5:
			time.Sleep(260 * time.Millisecond)
		case 6:
			time.Sleep(50 * time.Millisecond)
		}
		return x
	}).Debounce(250 * time.Millisecond).Subscribe(func(x int) {
		fmt.Print(x)
    })
    fmt.Println()

    fmt.Println("Distinct:")
    RxGo.Just(1, 2, 1, 1, 2, 3, 4, 4).Distinct().Subscribe(func(x int) {
		fmt.Print(x)
    })
    fmt.Println()

    fmt.Println("ElementAt:")
        for i:=0;i<6;i++{
            RxGo.Just(18,12,21,33,15,66).ElementAt(i).Subscribe(func(x int) {
                fmt.Printf("%d:%d\n",i,x)
            })
        }
    fmt.Println()

    fmt.Println("First:")
    RxGo.Just(18,12,21,33,15,66).First().Subscribe(func(x int) {
            fmt.Print(x)
        })
    fmt.Println()
    
    fmt.Println("IgnoreElements:")
    RxGo.Just(18,12,21,33,15,66).IgnoreElements().Subscribe(func(x int) {
            fmt.Print(x)
        })
    fmt.Println()
    
    fmt.Println("Last:")
    RxGo.Just(18,12,21,33,15,66).Last().Subscribe(func(x int) {
            fmt.Print(x)
        })
    fmt.Println()

    fmt.Println("Sample:")
    RxGo.Just(1,2,3,4,5,6).Map(func(x int) int {
                switch x {
                case 1:
                    time.Sleep(0 * time.Millisecond)
                case 2:
                    time.Sleep(10 * time.Millisecond)
                case 3:
                    time.Sleep(5 * time.Millisecond)
                case 4:
                    time.Sleep(20 * time.Millisecond)
                case 5:
                    time.Sleep(20 * time.Millisecond)
                case 6:
                    time.Sleep(50 * time.Millisecond)
                }
                return x
            }).Sample(25 * time.Millisecond).Subscribe(func(x int) {
                fmt.Print(x)
            })
    fmt.Println()
    

    fmt.Println("Skip:")
    RxGo.Just(18,12,21,33,15,66).Skip(3).Subscribe(func(x int) {
            fmt.Print(x)
        })
        
    fmt.Println()

    fmt.Println("SkipLast")
    RxGo.Just(18,12,21,33,15,66).SkipLast(3).Subscribe(func(x int) {
            fmt.Print(x)
        })
    fmt.Println()


    fmt.Println("Take:")
    RxGo.Just(18,12,21,33,15,66).Take(2).Subscribe(func(x int) {
            fmt.Print(x)
        })
    fmt.Println()

    fmt.Println("TakeLast:")
    RxGo.Just(18,12,21,33,15,66).TakeLast(2).Subscribe(func(x int) {
            fmt.Print(x)
        })
    fmt.Println()
}
```