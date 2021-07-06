package parseFiles

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"testing"
	"time"
)

func ParseFile(pathFile string) (files []string, err error) {

	file_bytes, e := ioutil.ReadFile(pathFile)
	if e != nil {
		err = fmt.Errorf("error:%v", e)
		return
	}
	splits := []string{".pdf", ".epub", ".mobi"}
	lines := strings.Split(string(file_bytes), "\n")
	files = make([]string, 0)
	for _, line := range lines {
		var i int
		for i = 0; i < len(splits); i++ {
			if strings.Contains(line, splits[i]) {
				idx := strings.Index(line, splits[i])
				line = line[:idx+len(splits[i])]
				fmt.Printf("line: %s \n", line)
				files = append(files, line)
			}
		}
		if i == len(splits)-1 {
			fmt.Printf("没有匹配到书名:%s\n", line)
			time.Sleep(time.Second)
		}

	}
	return

}

func TestParseFile(t *testing.T) {
	t.Log("解析PDF epub mobi电子书名称")
	lines, err := ParseFile("./test.txt")
	if err != nil {
		return
	}
	ioutil.WriteFile("./new.txt", []byte(strings.Join(lines, "\n")), os.ModeAppend)
}

// go test ./parse_files_test.go
