# 开发 web 服务程序
## 概述
开发简单 web 服务程序 cloudgo，了解 web 服务器工作原理。    
**任务目标**
1. 熟悉 go 服务器工作原理
2. 基于现有 web 库，编写一个简单 web 应用类似 cloudgo。
3. 使用 curl 工具访问 web 程序
4. 对 web 执行压力测试
## 任务要求
**基本要求**
1. 编程 web 服务程序类似 cloudgo 应用。
    - 支持静态文件服务
    - 支持简单 js 访问
    - 提交表单，并输出一个表格（必须使用模板）
2. 使用 curl 测试，将测试结果写入 README.md
3. 使用 ab 测试，将测试结果写入 README.md。并解释重要参数。

**扩展要求**     
选择以下一个或多个任务，以博客的形式提交。
1. 通过源码分析、解释一些关键功能实现
2. 选择简单的库，如 mux 等，通过源码分析、解释它是如何实现扩展的原理，包括一些 golang 程序设计技巧。

## 实验环境
- windows10
- go1.13

## 实验准备
- 安装需要的库
go get xxx
```
github.com/codegangsta/negroni
github.com/gorilla/mux   
github.com/unrolled/render
```
### 使用的库说明
1. net/http      
http包提供 HTTP 客户端和服务器实现。       
包括三个关键类型：    
   - Handler接口：所有请求的处理器、路由ServeMux都满足该接口。
    ```go
    type Handler interface {
        ServeHTTP(ResponseWriter, *Request)
    }
    ```

   - ServeMux结构体：HTTP请求的多路转接器（路口），它负责将每一个接收到URL与注册模式的列表进行匹配，并调用和URL最匹配的模式的处理器。它内部用一个map来保存所有处理器Handler
   - HandlerFunc适配器：HandlerFunc默认实现了ServeHTTP接口，具有`func(ResponseWriter, *Request)` 函数签名的函数f都可以通过HandlerFunc(f)显式转为 `Handler` 作为 HTTP的服务处理函数
    ```go
    type HandlerFunc func(ResponseWriter, *Request)
  	
  	// ServeHTTP calls f(w, r).
  	func (f HandlerFunc) ServeHTTP(w ResponseWriter, r *Request) {
  		f(w, r)
  	}
    ```

2. negroni                
在 Go 语言里，Negroni 是一个很地道的 Web 中间件，它是一个具备微型、非嵌入式、鼓励使用原生 net/http 库特征的中间件。            
   - `negroni.Classic()`提供一些默认的中间件，这些中间件在多数应用都很有用。       
   - 得益于`negroni.Handler`这个接口，Negroni 提供双向的中间件机制。       
[negroni-github](https://github.com/urfave/negroni/blob/master/translations/README_zh_CN.md)

3. mux                  
实现了一个请求路由和分发的 Go 框架。       
`mux.Router`根据已注册路由列表匹配传入请求，并调用与URL或其他条件匹配的路由的处理程序。
[Golang 第三方库学习 · mux](https://www.imooc.com/article/45868)
4. render        
一个软件包，提供了轻松呈现JSON，XML，文本，二进制数据和HTML模板的功能。       
[render-github](https://github.com/unrolled/render)
## 运行说明
进入项目根目录
```
go run main.go -p 9090
```
## 设计说明
### 静态文件服务支持
`service.go`
Go 的`net/http`包中提供了静态文件的服务，`ServeFile`和`FileServer` 等函数。    
- 首先在服务器上创建目录`assets`作为静态文件虚拟根目录。       
- 一条语句就实现了 `mx.PathPrefix("/").Handler(http.FileServer(http.Dir(webRoot + "/assets/")))` 静态文件服务。    
    - `http.Dir`是类型。将字符串转为`http.Dir`类型，这个类型实现了 `FileSystem`接口，把路径字符串映射到文件系统。（Dir 不是函数）
    - `http.FileServer()`是函数，返回`Handler`接口，该接口处理`http`请求，访问`root`的文件请求。
    - `mx.PathPrefix`添加前缀路径路由。 

### 支持js访问
**后端部分**                      
添加服务：`apitest.go`
```go
package service

import (
	"net/http"

	"github.com/unrolled/render"
)

func apiTestHandler(formatter *render.Render) http.HandlerFunc {

	return func(w http.ResponseWriter, req *http.Request) {
		formatter.JSON(w, http.StatusOK, struct {
			ID      string `json:"id"`
			Content string `json:"content"`
		}{ID: "8675309", Content: "Hello from Go!"})
	}
}
```
**前端部分**                          
描述页面结构：`index.html`
```html
<html>
<head>
  <link rel="stylesheet" href="css/main.css"/>
  <script src="http://code.jquery.com/jquery-latest.js"></script>
  <script src="js/hello.js"></script>
</head>
<body>
  <img src="images/cng.png" height="48" width="48"/>
  Sample Go Web Application!!
      <div>
          <p class="greeting-id">The ID is </p>
          <p class="greeting-content">The content is </p>
      </div>
</body>
</html>
```
定义页面行为：`hello.js`                         
```js
$(document).ready(function() {
    $.ajax({
        url: "/api/test"
    }).then(function(data) {
       $('.greeting-id').append(data.id);
       $('.greeting-content').append(data.content);
    });
});
```
### 使用模板输出
用到了库`github.com/unrolled/render`呈现HTML模板
- `formatter`构建，指定了模板的目录，模板文件的扩展名   
修改`server.go`中的代码   

```
	formatter := render.New(render.Options{
		Directory:  "templates",
		Extensions: []string{".html"},
		IndentJSON: true,
	})
```
- `apiTableHandler`使用了模板，接收表单请求，使用表单请求体数据填充模板
```go
func apiTableHandler(formatter *render.Render) http.HandlerFunc {

	return func(w http.ResponseWriter, req *http.Request) {
		req.ParseForm()
		formatter.HTML(w, http.StatusOK, "table", struct {
			Username string 
			Password string 
		}{Username: req.Form["username"][0], Password: req.Form["password"][0]})
	}
}
```

-  模板`table.html` 在`templates`目录中
```html
<html>
  <body>
    <table border="1">
      <tr>
        <th>Username</th>
        <th>Password</th>
      </tr>
      <tr>
        <td>{{.Username}}</td>
        <td>{{.Password}}</td>
      </tr>
    </table>
  </body>
</html>
```

- 提交表单的页面`inputtable.html`
```html
<html>
  <body>
    <form action="/table" method="post">
      Userame :&nbsp&nbsp&nbsp
      <input type="text" name="username" /><br /><br />
      Password :&nbsp&nbsp
      <input type="text" name="password" /><br /><br />
      <input type="submit" value="Submit" />
    </form>
  </body>
</html>
```

## 效果展示
### 支持静态文件服务
![](https://gitee.com/li-jia666/service-computing/raw/master/Cloudgo/img/1.png)
### 支持简单 js 访问
![](https://gitee.com/li-jia666/service-computing/raw/master/Cloudgo/img/2.png)
![](https://gitee.com/li-jia666/service-computing/raw/master/Cloudgo/img/5.png)
### 提交表单，并输出一个表格（必须使用模板）
![](https://gitee.com/li-jia666/service-computing/raw/master/Cloudgo/img/3.png)
![](https://gitee.com/li-jia666/service-computing/raw/master/Cloudgo/img/4.png)

## 测试
### curl
#### 支持静态文件服务
`curl -v http://localhost:9090/`
![](https://gitee.com/li-jia666/service-computing/raw/master/Cloudgo/img/11.png)
#### 支持简单 js 访问
`curl -v http://localhost:9090/api/test`
![](https://gitee.com/li-jia666/service-computing/raw/master/Cloudgo/img/12.png)
#### 提交表单，并输出一个表格（必须使用模板）
`curl -v -d "username=sysu;password=123" http://localhost:9090/table`
![](https://gitee.com/li-jia666/service-computing/raw/master/Cloudgo/img/13.png)
### ab
#### 支持静态文件服务
ab -n 1000 -c 100 http://localhost:9090/
```
D:\Apache24\bin>ab -n 1000 -c 100 http://localhost:9090/
This is ApacheBench, Version 2.3 <$Revision: 1879490 $>
Copyright 1996 Adam Twiss, Zeus Technology Ltd, http://www.zeustech.net/
Licensed to The Apache Software Foundation, http://www.apache.org/

Benchmarking localhost (be patient)
Completed 100 requests
Completed 200 requests
Completed 300 requests
Completed 400 requests
Completed 500 requests
Completed 600 requests
Completed 700 requests
Completed 800 requests
Completed 900 requests
Completed 1000 requests
Finished 1000 requests


Server Software:
Server Hostname:        localhost
Server Port:            9090

Document Path:          /
Document Length:        171 bytes

Concurrency Level:      100
Time taken for tests:   1.711 seconds
Complete requests:      1000
Failed requests:        0
Total transferred:      334000 bytes
HTML transferred:       171000 bytes
Requests per second:    584.53 [#/sec] (mean)
Time per request:       171.079 [ms] (mean)
Time per request:       1.711 [ms] (mean, across all concurrent requests)
Transfer rate:          190.66 [Kbytes/sec] received

Connection Times (ms)
              min  mean[+/-sd] median   max
Connect:        0    2  34.2      1    1083
Processing:     6   59   8.7     60      77
Waiting:        1   34  17.6     34      75
Total:          6   61  35.3     61    1141

Percentage of the requests served within a certain time (ms)
  50%     61
  66%     61
  75%     61
  80%     62
  90%     62
  95%     68
  98%     74
  99%     76
 100%   1141 (longest request)
 ```

#### 支持简单js访问
ab -n 1000 -c 100 http://localhost:9090/api/test

```
D:\Apache24\bin>ab -n 1000 -c 100 http://localhost:9090/api/test
This is ApacheBench, Version 2.3 <$Revision: 1879490 $>
Copyright 1996 Adam Twiss, Zeus Technology Ltd, http://www.zeustech.net/
Licensed to The Apache Software Foundation, http://www.apache.org/

Benchmarking localhost (be patient)
Completed 100 requests
Completed 200 requests
Completed 300 requests
Completed 400 requests
Completed 500 requests
Completed 600 requests
Completed 700 requests
Completed 800 requests
Completed 900 requests
Completed 1000 requests
Finished 1000 requests


Server Software:
Server Hostname:        localhost
Server Port:            9090

Document Path:          /api/test
Document Length:        53 bytes

Concurrency Level:      100
Time taken for tests:   0.897 seconds
Complete requests:      1000
Failed requests:        0
Total transferred:      176000 bytes
HTML transferred:       53000 bytes
Requests per second:    1115.32 [#/sec] (mean)
Time per request:       89.660 [ms] (mean)
Time per request:       0.897 [ms] (mean, across all concurrent requests)
Transfer rate:          191.70 [Kbytes/sec] received

Connection Times (ms)
              min  mean[+/-sd] median   max
Connect:        0    1  10.1      0     318
Processing:     5   54   8.6     57      62
Waiting:        1   30  15.9     30      62
Total:          5   55  13.2     57     372

Percentage of the requests served within a certain time (ms)
  50%     57
  66%     57
  75%     58
  80%     58
  90%     59
  95%     61
  98%     62
  99%     62
 100%    372 (longest request)
 ```

#### 提交表单，并输出一个表格（必须使用模板）
在当前目录下创建一个文件`post.txt`

编辑文件`post.txt`写入
username=sysu;password=123

ab -n 1000 -c 100 -p post.txt -T 'application/x-www-form-urlencoded' http://localhost:9090/table

```
D:\Apache24\bin>ab -n 1000 -c 100 -p post.txt -T 'application/x-www-form-urlencoded' http://localhost:9090/table`
This is ApacheBench, Version 2.3 <$Revision: 1879490 $>
Copyright 1996 Adam Twiss, Zeus Technology Ltd, http://www.zeustech.net/
Licensed to The Apache Software Foundation, http://www.apache.org/

Benchmarking localhost (be patient)
Completed 100 requests
Completed 200 requests
Completed 300 requests
Completed 400 requests
Completed 500 requests
Completed 600 requests
Completed 700 requests
Completed 800 requests
Completed 900 requests
Completed 1000 requests
Finished 1000 requests


Server Software:
Server Hostname:        localhost
Server Port:            9090

Document Path:          /table`
Document Length:        19 bytes

Concurrency Level:      100
Time taken for tests:   0.920 seconds
Complete requests:      1000
Failed requests:        0
Non-2xx responses:      1000
Total transferred:      176000 bytes
Total body sent:        186000
HTML transferred:       19000 bytes
Requests per second:    1087.50 [#/sec] (mean)
Time per request:       91.954 [ms] (mean)
Time per request:       0.920 [ms] (mean, across all concurrent requests)
Transfer rate:          186.91 [Kbytes/sec] received
                        197.53 kb/s sent
                        384.45 kb/s total

Connection Times (ms)
              min  mean[+/-sd] median   max
Connect:        0    1  10.4      1     328
Processing:     2   56  11.6     56      81
Waiting:        1   32  17.2     32      72
Total:          2   56  15.6     57     384

Percentage of the requests served within a certain time (ms)
  50%     57
  66%     57
  75%     57
  80%     58
  90%     70
  95%     76
  98%     79
  99%     80
 100%    384 (longest request)
 ```
#### 重要参数
**命令参数**

ab命令最基本的参数是-n和-c：              
-n 执行的请求数量                     
-c 并发请求个数                     

其他参数：                    
-t 测试所进行的最大秒数                   
-p 包含了需要POST的数据的文件                 
-T POST数据所使用的Content-type头信息              
-k 启用HTTP KeepAlive功能，即在一个HTTP会话中执行多个请求，默认时，不启用KeepAlive功能         

**结果参数**

Server Software: 服务器软件版本              
Server Hostname: 请求的URL                    
Server Port: 请求的端口号                    
Document Path: 请求的服务器的路径                  
Document Length: 页面长度 单位是字节                   
Concurrency Level: 并发数                     
Time taken for tests: 一共使用了的时间                   
Complete requests: 总共请求的次数                 
Failed requests: 失败的请求次数                                     
Total transferred: 总共传输的字节数 http头信息                      
HTML transferred: 实际页面传递的字节数                    
Requests per second: 每秒多少个请求                     
Time per request: 平均每个用户等待多长时间                     
Time per request: 服务器平均用多长时间处理               
Transfer rate: 传输速率                     
Connection Times: 传输时间统计                 
Percentage of the requests served within a certain time: 确定时间内服务请求占总数的百分比                 

## 源码分析
`DefaultServeMux`与`gorilla/mux`对比阅读

### DefaultServeMux
注册路由（DefautServeMux），即把一个模式（url）和对应的处理函数（handler）注册到DefautServeMux中。                  
因为http包中已经内置了路由：`DefaultServeMux`，所以可以直接调用`http.HandleFunc`函数注册一个处理器函数（handler）和对应的模式（pattern）到`DefaultServeMux`中。（当然，也可以自定义一个路由器，如：mux := http.NewServeMux()）                   
```go
func HandleFunc(pattern string, handler func(ResponseWriter, *Request)) {
    DefaultServeMux.HandleFunc(pattern, handler)
}
```
`DefaultServeMux`其实就是一个默认的`ServeMux`（多路复用路由器）的实例
```go
// NewServeMux allocates and returns a new ServeMux.
func NewServeMux() *ServeMux { return new(ServeMux) }

// DefaultServeMux is the default ServeMux used by Serve.
var DefaultServeMux = &defaultServeMux

var defaultServeMux ServeMux
```
`ServeMux`的结构如下:
```go
type ServeMux struct {
    mu    sync.RWMutex  // 并发处理涉及到的锁
    m     map[string]muxEntry // map的key(string)为一些url模式，value是一个muxEntry
    hosts bool // 判断是否在任意的规则下带有host信息
}
```
`muxEntry`的结构如下：
```go
type muxEntry struct {
    explicit bool // 判断是否精确匹配
    h        Handler // 路由器匹配的Handler
    pattern  string  // 路由器匹配的url模式
}
```
`DefaultServeMux`中的`HandleFunc(pattern, handler)`方法实际是定义在ServeMux中的：
```go
func (mux *ServeMux) HandleFunc(pattern string, handler func(ResponseWriter, *Request)) {
	mux.Handle(pattern, HandlerFunc(handler))
}
```
其中，`HandlerFunc`默认实现了`ServeHTTP`接口，具有`func(ResponseWriter, *Request)` 函数签名的函数f都可以通过HandlerFunc(f)显式转为 `Handler` 作为HTTP的服务处理函数

除此之外，上面的函数还调用了`ServeMux`的`Handle`方法(`mux.Handle`)，将`pattern`和`handler`函数做了一个map映射。mux中的Handle函数如下：
```go
// Handle registers the handler for the given pattern.
// If a handler already exists for pattern, Handle panics.
func (mux *ServeMux) Handle(pattern string, handler Handler) {
    mux.mu.Lock()
    defer mux.mu.Unlock()
    if pattern == "" {
        panic("http: invalid pattern " + pattern)
    }
    if handler == nil {
        panic("http: nil handler")
    }
    if mux.m[pattern].explicit { // 如果pattern与handler匹配
        panic("http: multiple registrations for " + pattern)
    }

    if mux.m == nil {   // 如果map不存在，则建立map
        mux.m = make(map[string]muxEntry)
    }
    // 判断pattern与handler是否精确匹配
    mux.m[pattern] = muxEntry{explicit: true, h: handler, pattern: pattern}

    if pattern[0] != '/' {
        mux.hosts = true
    }

    n := len(pattern)
    if n > 0 && pattern[n-1] == '/' && !mux.m[pattern[0:n-1]].explicit {

        path := pattern
        if pattern[0] != '/' {
            path = pattern[strings.Index(pattern, "/"):]
        }
        url := &url.URL{Path: path}
        mux.m[pattern[0:n-1]] = muxEntry{
            h: RedirectHandler(url.String(),StatusMovedPermanently), pattern: pattern
        }
	}
}    
```
分析以上代码，可知该函数的功能为把一个URL模式(pattern)和与其匹配的处理函数(handler)绑定到`muxEntry`的`map`上，这个map就相当于一个pattern和handler的匹配表，存储在`ServeMux`结构中。前面提到了`DefaultServeMu`x是`ServeMux`的一个实例，因此，调用`HandleFunc(pattern, handler)`方法最终`pattern`和其对应的`handler`绑定到了`DefautServeMux`中，这就代表路由注册完毕了。

路由查找的过程实际上就是遍历路由表的过程，返回最长匹配请求路径的路由信息，找不到则返回NotFoundHandler。如果路径以`xxx/`结尾，则只要满足`/xxx/*` 就符合该路由。
```go
  func (mux *ServeMux) handler(host, path string) (h Handler, pattern string) {
          mux.mu.RLock()
          defer mux.mu.RUnlock()
          if h == nil {
                  h, pattern = mux.match(path)
          }
          if h == nil {
                  h, pattern = NotFoundHandler(), ""
          }
          return
  }

  func (mux *ServeMux) match(path string) (h Handler, pattern string) {
          var n = 0
          for k, v := range mux.m {
                  if !pathMatch(k, path) {
                  continue
                  }
                  //找出匹配度最长的
                  if h == nil || len(k) > n {
                  n = len(k)
                  h = v.h
                  pattern = v.pattern
                  }
          }
          return
  }

  func pathMatch(pattern, path string) bool {
  n := len(pattern)
          if pattern[n-1] != '/' {
                  return pattern == path
          }
          return len(path) >= n && path[0:n] == pattern
  }
  ```

#### 总结
`DefaultServeMux`提供的路由处理器虽然简单易上手，但是存在很多不足，比如：
- 不支持正则路由
- 只支持路径匹配，不支持按照Method，header，host等信息匹配，所以也就没法实现RESTful架构

### gorilla/mux
为此，我们可以使用第三方库`gorilla/mux`提供的更加强大的路由处理器（mux 代表 HTTP request multiplexer，即 HTTP 请求多路复用器），和`http.ServeMux`实现原理一样，gorilla/mux 提供的路由器实现类`mux.Router`也会匹配用户请求与系统注册的路由规则，然后将用户请求转发过去。

`mux.Router`主要具备以下特性：

- 实现了`http.Handler`接口，所以和`http.ServeMux`完全兼容；
- 可以基于 URL 主机、路径、前缀、scheme、请求头、请求参数、请求方法进行路由匹配；
- URL 主机、路径、查询字符串支持可选的正则匹配；
- 支持构建或反转已注册的 URL 主机，以便维护对资源的引用；
- 支持路由嵌套（类似 Laravel 中的路由分组），以便不同路由可以共享通用条件，比如主机、路径前缀等。

Router实现了`http.Handler`接口，所以可以被注册处理请求。
```go
type Router struct {
	//当没有路由匹配时使用的可配置Handler。
	NotFoundHandler http.Handler

	// //当请求方法与路由不匹配时使用的可配置Handler。
	MethodNotAllowedHandler http.Handler

	// 要匹配的路由，按顺序。
	routes []*Route

	// Routes by name for URL building.
	namedRoutes map[string]*Route

	//如果为true，请在处理请求后不要清除请求上下文。
	//
	//不推荐使用：无效，因为上下文存储在请求本身上。
	KeepContext bool

	//找到匹配项后调用的中间件数组
	middlewares []middleware

	// 与`Route`共享的配置
	routeConf
}
```
路由信息存放在Route数组中,Route 存储匹配请求和路径的信息。
```go
type Route struct {
	// Request handler for the route.
	handler http.Handler
	// If true, this route never matches: it is only used to build URLs.
	buildOnly bool
	// The name used to build URLs.
	name string
	// Error resulted from building a route.
	err error

	// "global" reference to all named routes
	namedRoutes map[string]*Route

	// config possibly passed in from `Router`
	routeConf
}
```

当请求到来时，`ServeHTTP`调度在匹配的路由中注册的处理函数。
```go
func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if !r.skipClean {
		path := req.URL.Path
		if r.useEncodedPath {
			path = req.URL.EscapedPath()
		}
		// Clean path to canonical form and redirect.
		if p := cleanPath(path); p != path {

			// Added 3 lines (Philip Schlump) - It was dropping the query string and #whatever from query.
			// This matches with fix in go 1.2 r.c. 4 for same problem.  Go Issue:
			// http://code.google.com/p/go/issues/detail?id=5252
			url := *req.URL
			url.Path = p
			p = url.String()

			w.Header().Set("Location", p)
			w.WriteHeader(http.StatusMovedPermanently)
			return
		}
	}
	var match RouteMatch
    var handler http.Handler
    //从router中找出符合要求的handler
	if r.Match(req, &match) {
		handler = match.Handler
		req = requestWithVars(req, match.Vars)
		req = requestWithRoute(req, match.Route)
	}
    //如果找到的不符合要求
	if handler == nil && match.MatchErr == ErrMethodMismatch {
		handler = methodNotAllowedHandler()
	}
    //如果找不到，执行没有找到handler的处理方式
	if handler == nil {
		handler = http.NotFoundHandler()
	}

	handler.ServeHTTP(w, req)
}

```
`Match`函数会按顺序遍历数组，找到第一个匹配的路由，并执行对应的处理函数，如果找不到则执行`NotFoundHandler`。
```go
func (r *Router) Match(req *http.Request, match *RouteMatch) bool {
	for _, route := range r.routes {
        //调用route的Match方法
		if route.Match(req, match) {
			// Build middleware chain if no error was found
			if match.MatchErr == nil {
				for i := len(r.middlewares) - 1; i >= 0; i-- {
					match.Handler = r.middlewares[i].Middleware(match.Handler)
				}
			}
			return true
		}
	}

	if match.MatchErr == ErrMethodMismatch {
		if r.MethodNotAllowedHandler != nil {
			match.Handler = r.MethodNotAllowedHandler
			return true
		}

		return false
	}

	// Closest match for a router (includes sub-routers)
	if r.NotFoundHandler != nil {
		match.Handler = r.NotFoundHandler
		match.MatchErr = ErrNotFound
		return true
	}

	match.MatchErr = ErrNotFound
	return false
}
```
每一个Route中包含一个matcher数组，是所有限定条件的集合，用来匹配：请求头、方法（GET/POST）等。    
matcher是一个返回bool值的接口。添加路由限定条件就是往matcher数组中增加一个限定函数。         

以方法匹配为例简单看一下实现：
```go
// methodMatcher matches the request against HTTP methods.
type methodMatcher []string

//methodMatcher就是取出r.Method然后判断该方式是否是设定的Method
func (m methodMatcher) Match(r *http.Request, match *RouteMatch) bool {
   return matchInArray(m, r.Method)
}

// matchInArray returns true if the given string value is in the array.
func matchInArray(arr []string, value string) bool {
    for _, v := range arr {
     if v == value {
         return true
     }
   }
   return false
}

// Methods adds a matcher for HTTP methods.
// It accepts a sequence of one or more methods to be matched, e.g.:
// "GET", "POST", "PUT".
func (r *Route) Methods(methods ...string) *Route {
	for k, v := range methods {
		methods[k] = strings.ToUpper(v)
	}
	return r.addMatcher(methodMatcher(methods))
}
```
方法匹配（匹配GET/POST等）的方式就是在支持的方法中寻找本次请求的Request的Method。


Route.Match()会遍历matcher数组，只有数组中所有的元素都返回true时则说明此请求满足该路由的限定条件。
```go
func (r *Route) Match(req *http.Request, match *RouteMatch) bool {
	if r.buildOnly || r.err != nil {
		return false
	}

	var matchErr error

	// Match everything.
	for _, m := range r.matchers {
		if matched := m.Match(req, match); !matched {
			if _, ok := m.(methodMatcher); ok {
				matchErr = ErrMethodMismatch
				continue
			}

			// Ignore ErrNotFound errors. These errors arise from match call
			// to Subrouters.
			//
			// This prevents subsequent matching subrouters from failing to
			// run middleware. If not ignored, the middleware would see a
			// non-nil MatchErr and be skipped, even when there was a
			// matching route.
			if match.MatchErr == ErrNotFound {
				match.MatchErr = nil
			}

			matchErr = nil
			return false
		}
	}

	if matchErr != nil {
		match.MatchErr = matchErr
		return false
	}

	if match.MatchErr == ErrMethodMismatch && r.handler != nil {
		// We found a route which matches request method, clear MatchErr
		match.MatchErr = nil
		// Then override the mis-matched handler
		match.Handler = r.handler
	}

	// Yay, we have a match. Let's collect some info about it.
	if match.Route == nil {
		match.Route = r
	}
	if match.Handler == nil {
		match.Handler = r.handler
	}
	if match.Vars == nil {
		match.Vars = make(map[string]string)
	}

	// Set variables.
	r.regexp.setMatch(req, match, r)
	return true
}
```

Router中同样定义了`HandleFunc(path, handler)`方法,用于创建路由匹配并定义URL路径的匹配规则。
```go
func (r *Router) HandleFunc(path string, f func(http.ResponseWriter,
	*http.Request)) *Route {
	return r.NewRoute().Path(path).HandlerFunc(f)
}
```
它使用`NewRoute()`创建了一个新的路由Route，接着调用`Path()`为Route添加路由规则。
```go
func (r *Route) Path(tpl string) *Route {
	r.err = r.addRegexpMatcher(tpl, regexpTypePath)
	return r
}
```
`Path()`调用了`addRegexpMatcher()`，根据传入的路径字符串tql创建正则表达式并调用`newRouteRegexp()`。`newRouteRegexp()`解析了tql，返回一个`routeRegexp`，用于匹配主机、路径或查询字符串，最后加入matchers数组中。

```go
// addRegexpMatcher adds a host or path matcher and builder to a route.
func (r *Route) addRegexpMatcher(tpl string, typ regexpType) error {
	if r.err != nil {
		return r.err
	}
	r.regexp = r.getRegexpGroup()
	if typ == regexpTypePath || typ == regexpTypePrefix {
		if len(tpl) > 0 && tpl[0] != '/' {
			return fmt.Errorf("mux: path must start with a slash, got %q", tpl)
		}
		if r.regexp.path != nil {
			tpl = strings.TrimRight(r.regexp.path.template, "/") + tpl
		}
	}
	rr, err := newRouteRegexp(tpl, typ, routeRegexpOptions{
		strictSlash:    r.strictSlash,
		useEncodedPath: r.useEncodedPath,
	})
	if err != nil {
		return err
	}
	for _, q := range r.regexp.queries {
		if err = uniqueVars(rr.varsN, q.varsN); err != nil {
			return err
		}
	}
	if typ == regexpTypeHost {
		if r.regexp.path != nil {
			if err = uniqueVars(rr.varsN, r.regexp.path.varsN); err != nil {
				return err
			}
		}
		r.regexp.host = rr
	} else {
		if r.regexp.host != nil {
			if err = uniqueVars(rr.varsN, r.regexp.host.varsN); err != nil {
				return err
			}
		}
		if typ == regexpTypeQuery {
			r.regexp.queries = append(r.regexp.queries, rr)
		} else {
			r.regexp.path = rr
		}
	}
	r.addMatcher(rr)
	return nil
}

// matcher types try to match a request.
type matcher interface {
	Match(*http.Request, *RouteMatch) bool
}

// addMatcher adds a matcher to the route.
func (r *Route) addMatcher(m matcher) *Route {
	if r.err == nil {
		r.matchers = append(r.matchers, m)
	}
	return r
}

```

最后是`HandlerFunc()`，它给特定的路由匹配了请求响应函数`YourHandler()`
```go

// HandlerFunc sets a handler function for the route.
func (r *Route) HandlerFunc(f func(http.ResponseWriter, *http.Request)) *Route {
	return r.Handler(http.HandlerFunc(f))
}

// Handler sets a handler for the route.
func (r *Route) Handler(handler http.Handler) *Route {
	if r.err == nil {
		r.handler = handler
	}
	return r
}
```

#### 总结
`Mux`完全兼容`http.ServerMux`，相比之下，有几个有点是原生`ServerMux`不具备的：
- 支持正则路由
- 支持按照Method，header，host等信息匹配.