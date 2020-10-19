# goini
使用goini更简单的读取go的ini配置文件


## 使用方法

ini配置文件格式样列

```ini
; possible values : production,development
app_mode = development

[paths]
; Path to where grafana can store temp files, sessions, and the sqlite3 db (if that is used)
data = /home/git/grafana

[server]
; Protocol (http or https)
protocol = http

; The http port  to use
http_port = 9999

; Redirect to correct domain if host header does not match domain
; Prevents DNS rebinding attacks
enforce_domain = true
```
`goini.Watch(filename string,listener Listener) (configuration, error)`
接收文件名以及`listener`接口作为参数，返回key-value式样的配置解析结果与自定义错误。其功能为监听自函数运行以来发生的一次配置文件变化并返回最新的配置文件解析内容。

```
var mylistener goini.ListenFunc =func (inifile string){
	fmt.Println("Deal with configure change here")
}
conf,err:=goini.Watch(./example1,mylistener)
```
简单的使用案例

```go
	var mylistener goini.ListenFunc =func (inifile string){
		fmt.Println("Deal with configure change here")
	}
	//在配置文件中增加一行
	var add = func(filepath string,content string) {
		file, err := os.OpenFile(filepath, os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			fmt.Println("文件打开失败", err)
		}
		defer file.Close()
		write := bufio.NewWriter(file)
		write.WriteString("\r\n")
		write.WriteString(content)
		write.Flush()
	}
	go add(os.Args[1],os.Args[2])
	conf,err:=goini.Watch(os.Args[1],mylistener)
	for _, v := range conf {
		for key,value := range v{
			fmt.Printf("%s : %s\n", key ,value)
		}
	}
	goini.CheckErr(err)
```