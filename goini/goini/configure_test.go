package goini

import (
	"reflect"
	"fmt"
	"testing"
	"os"
	"bufio"
)
func TestWatch(t *testing.T) {
	filepath:="./conf/conf.ini"
	var mylistener ListenFunc = func(inifile string) {
	}
	var change = func(filepath string) {
		file, err := os.OpenFile(filepath, os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			fmt.Println("文件打开失败", err)
		}
		//及时关闭file句柄
		defer file.Close()
		write := bufio.NewWriter(file)
		write.WriteString("\r\n"+anno+"a=b")
		write.Flush()
	}
	go change(filepath)
	conf, _ := Watch(filepath, mylistener)
	
	var testconf configuration
	var start map[string]string = make(map[string]string)
	start["app_mode"] = "development"
	var paths map[string]string = make(map[string]string)
	paths["data"] = "/home/git/grafana"
	var server map[string]string = make(map[string]string)
	server["protocol"] = "http"
	server["http_port"] = "9999"
	server["enforce_domain"] = "true"
	testconf = append(testconf, start)
	testconf = append(testconf, paths)
	testconf = append(testconf, server)

	if !reflect.DeepEqual(testconf, conf) {
		t.Errorf("expected %+v but got %+v", testconf, conf)
	}
}

func ExampleWatch() {
	filepath:="./conf/conf.ini"
	var change = func(filepath string) {
		file, err := os.OpenFile(filepath, os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			fmt.Println("文件打开失败", err)
		}
		//及时关闭file句柄
		defer file.Close()
		write := bufio.NewWriter(file)
		write.WriteString("\r\n"+anno+"a=b")
		write.Flush()
	}
	var mylistener ListenFunc =func (inifile string){
	}
	go change(filepath)
	conf,err:=Watch(filepath,mylistener)
	for _, v := range conf {
		for key,value := range v{
			fmt.Printf("%s : %s\n", key ,value)
		}
	}
	CheckErr(err)
	// Unordered output: 
	// app_mode : development 
	// data : /home/git/grafana 
	// protocol : http 
	// http_port : 9999 
	// enforce_domain : true
}