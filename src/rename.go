package main

import (
	"flag"
	"strings"
	"log"
	"os"
	"fmt"
	"path/filepath"
)

func main() {
	var dir string
	flag.StringVar(&dir,"d","","-d directory path")
	var pattern string
	flag.StringVar(&pattern, "r","","-r replace")
	flag.Parse()
	if dir == "" || pattern == "" {
		fmt.Println("rename.exe -d ./ -r 'xxxx'")
		flag.Usage()
		return
	}

	pathfiles := make([]string,0)
	filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			pathfiles = append(pathfiles, path)
		}
		return nil
	})
	// 遍历文件路径，修改文件名
	for _, path := range pathfiles {
		// 替换字符串
		newPath := strings.Replace(path,pattern,"",-1)
		// Ext后缀名
		//filepath.Join(filepath.Dir(path),strings.Replace(path,pattern,"",-1) + filepath.Ext(path))
		log.Println("new ",newPath)
		os.Rename(path,newPath)
	}

}
