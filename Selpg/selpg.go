package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"

	flag "github.com/spf13/pflag"
)


var (
	startPage = flag.IntP("start_page", "s", -1, "start page")
	endPage = flag.IntP("end_page", "e", -1, "end page")
	pageLen = flag.IntP("page_len", "l", 72, "line number per page")
	pageBreak = flag.BoolP("use_page_break", "f", false, "pages devided by page break \\f")
	printDest = flag.StringP("print_dest", "d", "", "specify the printer")
	pageType string
)

func main() {
	handleError(parseArgs())
	handleError(readAndWrite())
}

func parseArgs() error {
	flag.Parse() // 使用pflag提供的方法解析命令行参数

	if *pageBreak { // 判断文件分页格式
		pageType = "f"
	} else {
		pageType = "l"
	}

	if *startPage == -1 || *endPage == -1 { // 判断参数是否完整（起始页和结束页）
		printUsage()
		return errors.New("The arguments are not enough")
	}

	if *startPage <= 0 || *endPage <= 0 { // 判断页码是否合理
		return errors.New("The page number can not be negative")
	}

	if *startPage > *endPage { // 判断页码是否合理
		return errors.New("The start page cannot be greater than end page")
	}

	if pageType == "l" && *pageLen <= 0 { // 判断页长是否合理
		return errors.New("The line number per page can not be negative")
	}

	if pageType == "f" && *pageLen != 72 { // 判断是否同时设置两种分页格式
		return errors.New("-f and -l linePerPage cannot be set at the same time")
	}

	return nil
}

func readAndWrite() error {
	var reader *bufio.Reader // 定义一个输入流
	var writer *bufio.Writer // 定义一个输出流

	if flag.NArg() == 0 { // 根据[FILE]参数给输入流赋值
		reader = bufio.NewReader(os.Stdin)
	} else {
		input, err := os.Open(flag.Arg(0))
		if err != nil {
			return err
		}
		defer input.Close()
		reader = bufio.NewReader(input)
	}

	if len(*printDest) == 0 { // 根据-d参数给输出流赋值
		writer = bufio.NewWriter(os.Stdout)
	} else {
		
		cmd := exec.Command("lp", "-d"+*printDest)
		output, err := cmd.StdinPipe()
		if err != nil {
			return err
		}
		defer output.Close()
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Start(); err != nil {
			return err
		}
		writer = bufio.NewWriter(output)
	}
	defer writer.Flush()

	var pageSpliter byte // 获取分页格式标识符
	if pageType == "f" {
		pageSpliter = '\f'
	} else {
		pageSpliter = '\n'
	}

	pages, lines := 1, 0
	for {
		sub, err := reader.ReadBytes(pageSpliter) // 根据标识符读取一段内容
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}
		if pageType == "f" { // 进行行/页相应处理
			pages++
		} else {
			lines++
			if lines > *pageLen {
				lines = 1
				pages++
			}
		}
		if pages >= *startPage && pages <= *endPage {
			if _, err := writer.Write(sub); err != nil {
				return err
			}
		} else if pages > *endPage{
			break
		}
	}

	return nil
}

func printUsage() { // 打印使用方法
	fmt.Fprintln(os.Stderr, "Usage: selpg [-s startPage] [-e endPage] [-l linePerPage | -f] [-d destination] input_file >output_file 2>error_file")
	flag.PrintDefaults()
	os.Exit(2)
}

func handleError(err error) { // 错误处理
	if err != nil {
		if _, err2 := fmt.Fprintf(os.Stderr, "%s\n", err.Error()); err2 != nil {
			panic(err2)
		}
		os.Exit(1)
	}
}
