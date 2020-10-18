//读配置文件程序包
package goini

import (
	"bufio"
	"io"
	"os"
	"strings"
	"fmt"
	"time"
)

//自定义一个错误类型
type myError struct {
	etime time.Time
	info  string
}
 
//实现error接口
func (e *myError) Error() string {
	return fmt.Sprintf("time %s : %s",e.etime.Format("2006-01-02 15:04:05"),e.info)
}


var (
    fi os.FileInfo
    err error
)

type Listener interface  { listen(inifile string)  }

type ListenFunc func(string)

//监听自函数运行以来发生的一次配置文件变化
func (l ListenFunc) listen(inifile string){
	fi, err=os.Stat(inifile)
	if err != nil {
		myerr:=&myError{time.Now(),err.Error()}
		fmt.Println(myerr.Error())
		return
	}
	modinit:=fi.ModTime().Unix()//获取文件的修改时间
	for{
		fi, err=os.Stat(inifile)
		if err != nil {
			myerr:=&myError{time.Now(),err.Error()}
			fmt.Println(myerr.Error())
			return 
		}
		mod:=fi.ModTime().Unix()
		if modinit!=mod{ //如修改时间变化，说明文件被修改，跳出循环(只监听一次变化)
			break
		}
	}
	l(inifile); //调用接收者函数，开发者在这个函数内自己决定如何处理配置变化
}

type configuration []map[string]string

//监听自函数运行以来发生的一次配置文件变化并返回最新的配置文件解析内容
func Watch(filename string,listener Listener) (configuration, error){
	listener.listen(filename) //使用方法listen监听配置文件是否被修改，如果修改，函数流程往下
	con:=GetConfig(filename) 
	err:=con.Analyse() //读文件，提取出配置信息
	var result configuration //返回一组key，values对
	for _, v := range con.conflist {
		for _, value := range v {
			result = append(result, value)
		}
	}
	return result,err
}

type Config struct {
	filepath string                         //配置文件的路径
	conflist []map[string]map[string]string //配置信息的切片
}

//创建一个空的配置文件结构
func GetConfig(filepath string) *Config {
	c := new(Config)
	c.filepath = filepath

	return c
}

//读取配置文件，提取参数
func (c *Config) Analyse() error {
	file, err := os.Open(c.filepath)
	if err != nil {
		myerr:=&myError{time.Now(),err.Error()}
		return myerr
	}
	defer file.Close()
	var data map[string]map[string]string=make(map[string]map[string]string)
	var section string="start"
	data[section] = make(map[string]string)

	buf := bufio.NewReader(file)
	for {
		l, err := buf.ReadString('\n')
		line := strings.TrimSpace(l)
		if err != nil {
			if err != io.EOF {
				myerr:=&myError{time.Now(),err.Error()}
				return myerr
			}
			if len(line) == 0 {
				break
			}
		}
		switch {
		case len(line) == 0:
		case string(line[0]) == "#":
		case line[0] == '[' && line[len(line)-1] == ']':
			section = strings.TrimSpace(line[1 : len(line)-1])
			data = make(map[string]map[string]string)
			data[section] = make(map[string]string)
		default:
			i := strings.IndexAny(line, "=")
			if i == -1 {
				continue
			}
			value := strings.TrimSpace(line[i+1 : len(line)])
			data[section][strings.TrimSpace(line[0:i])] = value
			if c.uniquappend(section) == true {
				c.conflist = append(c.conflist, data)
			}
		}
	}
	return nil
}

//错误处理
func CheckErr(err error) {
	if err != nil {
		fmt.Println(err.Error())
	}
}

//过滤重复元素
func (c *Config) uniquappend(conf string) bool {
	for _, v := range c.conflist {
		for k, _ := range v {
			if k == conf {
				return false
			}
		}
	}
	return true
}