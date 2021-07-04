package main

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"
)

type MyFileInfo struct {
	IsDir    bool   // 是否为目录
	Folder   string // 目录名称
	FullPath string // 文件全路径

}

const (
	PathSeparator string = string(os.PathSeparator)
)

// 获取当前目录下的所有文件
func GetFiles(folder string, info string) (files []MyFileInfo, err error) {
	rFiles, err := ioutil.ReadDir(folder)
	if err != nil {
		err = fmt.Errorf("GetFiles ReadDir failed %v", err)
		return
	}
	for _, file := range rFiles {
		if file.IsDir() {
			nextFolder := folder + PathSeparator + file.Name()
			s, e := GetFiles(nextFolder, file.Name())
			if err != nil {
				err = fmt.Errorf("reverse dir error:%v", e)
			}
			files = append(files, s...)
		} else {
			fullFilePath := folder + PathSeparator + file.Name()
			fInfo := MyFileInfo{}
			if len(info) == 0 {
				fInfo.IsDir = false
				fInfo.Folder = ""
				fInfo.FullPath = fullFilePath
			} else {
				fInfo.IsDir = true
				fInfo.Folder = info
				fInfo.FullPath = fullFilePath
			}
			files = append(files, fInfo)
		}
	}
	return
}
func IsFile(filePath string) bool {
	f, e := os.Stat(filePath)
	if e != nil {
		return false
	}
	return !f.IsDir()
}

func Copy(src, dest string) error {
	// Gather file information to set back later.
	si, err := os.Lstat(src)
	if err != nil {
		return err
	}

	// Handle symbolic link.
	if si.Mode()&os.ModeSymlink != 0 {
		target, err := os.Readlink(src)
		if err != nil {
			return err
		}
		// NOTE: os.Chmod and os.Chtimes don't recoganize symbolic link,
		// which will lead "no such file or directory" error.
		return os.Symlink(target, dest)
	}

	sr, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sr.Close()

	dw, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer dw.Close()

	if _, err = io.Copy(dw, sr); err != nil {
		return err
	}

	// Set back file information.
	if err = os.Chtimes(dest, si.ModTime(), si.ModTime()); err != nil {
		return err
	}
	return os.Chmod(dest, si.Mode())
}
func RemoveFile(fullPath string) error {
	err := os.Remove(fullPath)
	if err != nil {
		fmt.Println("删除文件失败:",err)
	}
	fmt.Println(fullPath," 已经删除")
	return err
}
func CreateDir(dir string) error {
	fInfo, err := os.Stat(dir)
	if err != nil {
		err := os.MkdirAll(dir, os.ModePerm)
		if err != nil {
			fmt.Printf("创建目录失败err:[%s], stat ret fileInfo:[%v]\n", err, fInfo)
			return errors.New(fmt.Sprintf("Create dir:[%s] failed, error:%s", dir, err))
		}
		fmt.Printf("创建目录:%s\n", dir)
	}
	return nil
}
func DeleteDupFile(file MyFileInfo, filter []string) {
	if IsFile(file.FullPath) {
		for _, f := range filter {
			if strings.Contains(file.FullPath, f) {
				fmt.Println("file ",file.FullPath)
				RemoveFile(file.FullPath)
			}
		}
	}
}
func MoveFile(file MyFileInfo, descFolder string, filter []string) {
	if IsFile(file.FullPath) {
		lastIndex := strings.LastIndex(file.FullPath, PathSeparator)
		srcFilename := string(file.FullPath[lastIndex+1:])
		for _, f := range filter {
			if strings.Contains(file.FullPath, f) {
				//RemoveFile(f)
				goto END
			}
		}
		if file.IsDir {
			var desc string
			if strings.HasSuffix(descFolder, PathSeparator) {
				desc = descFolder + file.Folder
			} else {
				desc = descFolder + PathSeparator + file.Folder
			}
			CreateDir(desc)
			fmt.Println("需要创建目录:", desc)
		}

		err := Copy(file.FullPath, descFolder+PathSeparator+file.Folder+PathSeparator+srcFilename)
		if err != nil {
			fmt.Println("文件copy 失败", err)
		} else {
			fmt.Println("拷贝文件成功")
		}
	}
END:
	fmt.Println("end")
}
func GetFullAllFiles(folder string) (files []MyFileInfo, err error) {
	return GetFiles(folder, "")
}
func main() {

	files, err := GetFullAllFiles("D:\\天翼云盘下载\\go项目")
	if err != nil {
		fmt.Println(err)
	}
	Filter := []string{"(1)", "(2)"}
	for _, MyFileInfo := range files {
		DeleteDupFile(MyFileInfo,  Filter)
		continue
		MoveFile(MyFileInfo, "D:/天翼云盘下载/tmp", Filter)
		//fmt.Println(MyFileInfo.Folder," ",MyFileInfo.FullPath)
	}
	fmt.Println("hello world ")

}
