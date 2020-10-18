package main

import (
	"fmt"
	"os"
	"bufio"
	"github.com/user/goini"
)

func main() {
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
	goini.CheckErr(err)
	for _, v := range conf {
		for key,value := range v{
			fmt.Printf("%s : %s\n", key ,value)
		}
	}
}