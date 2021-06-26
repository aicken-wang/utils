package main

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

func ReadFiles(filepath string) {
	f, err := os.Open(filepath)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	bufio := bufio.NewReader(f)
	index := 0
	urlBase := ""
	for {
		buf, err := bufio.ReadBytes('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			panic(err)
		}
		input := string(buf[:len(buf)-1])
		strLen := len(input)
		if strLen == 0 {
			continue
		}
		index++
		if index == 1 {
			urlBase = input
			continue
		}
		url := fmt.Sprintf("%s%s", urlBase, input)
		fmt.Printf("第 %d 本, url:%s\n", index-1, url)
		time.Sleep(time.Millisecond * 200)
		downloadPDF(url, input)
	}
}
func downloadPDF(url, filename string) {
	res, err := http.Get(url)
	if err != nil {
		fmt.Println("下载失败", err)
		return
	}
	defer res.Body.Close()
	contentLen, ok := res.Header["Content-Length"]
	if ok {
		fmt.Println("contentLen:", contentLen)
	}

	fd, err := os.Create("./2021/06/" + filename)
	if err != nil {
		fmt.Println("创建文件失败", err)
		return
	}
	defer fd.Close()
	io.Copy(fd, res.Body)
}
func main() {
	ReadFiles("./link/202106.txt")
}
