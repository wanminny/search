package util

import (
	"path"
	"strings"
	"os"
	"log"
	"path/filepath"
)

func GetCurrentDirectory() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0])) //返回绝对路径  filepath.Dir(os.Args[0])去除最后一个元素的路径
	if err != nil {
		log.Fatal(err)
	}
	return strings.Replace(dir, "\\", "/", -1) //将\替换成/
}


func GetFileName(fullFilename string)  string {

	//fullFilename := "/Documents/test.txt"
	//fmt.Println("fullFilename =", fullFilename)
	var filenameWithSuffix string
	filenameWithSuffix = path.Base(fullFilename) //获取文件名带后缀
	//fmt.Println("filenameWithSuffix =", filenameWithSuffix)

	var fileSuffix string
	fileSuffix = path.Ext(filenameWithSuffix) //获取文件后缀
	//fmt.Println("fileSuffix =", fileSuffix)

	var filenameOnly string
	filenameOnly = strings.TrimSuffix(filenameWithSuffix, fileSuffix)//获取文件名
	//fmt.Println("filenameOnly =", filenameOnly)

	return filenameOnly
}

//文件是否已存在
func PathExist(path string) bool {
	_, err := os.Stat(path)
	if err != nil && os.IsNotExist(err) {
		return false
	}
	return true
}
