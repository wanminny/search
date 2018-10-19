package search

import (
	"os"
	"bufio"
	"io"
	"log"
	"strings"
	"time"
	"gobible/logmanager/cli/http/services/util"
)


const TIMEFORMAT = "20060102"
const TIMEFORMATZIP = "20060102150405"
const DIR = "/gopath/src/gobible/logmanager/cli"

var (
	startTime string //开始时间
	endTime   string //结束时间
	fileName  string //要查找的文件名称
	deviceId  string //设备id

	directory string //要复制的目录

	help bool //显示帮助

	//默认多少天以前 两天前
	delta = "-48h"

	//生成的目标普通文件(.log)
	destFileDir string

	//生成的目标的压缩文件(gzip)
	destGzipFileDir string

	//谁否可以压缩文件了
	gzipOK chan struct{} = make(chan struct{}, 1)

	//end
	end chan int = make(chan int)

	// 需要解压的文件的扩展名
	extName = ".gz"

	//要查找的文件名称的前缀
	prefix = "tcp_stdout.log-"

	//根据开始 结束时间 返回名称
	dirs = make([]string, 0)

	dirsMap = make(map[string]bool)

	//目标文件夹
	//destDir = "/tar"

	//当前的目录的拷贝目录
	copyDirTar = "./copy-dir-tar"

	//当前的文件的拷贝目录
	copyFileTar = "./copy-file-tar"

	fileMaps  = make(map[string]FileInfo)
)

type FileInfo struct {
	//文件名称目录
	Name string
	//是否是压缩文件
	gz bool

	fullName string

	empty bool
}

// 创建文件夹;处理完成后并删除！
func mkdirs()  {

	//先删除 后创建
	//err :=  os.Remove(copyDirTar)
	//if err != nil{
	//	log.Println(err)
	//}
	//err =  os.Remove(copyFileTar)
	//if err != nil{
	//	log.Println(err)
	//}

	err :=os.Mkdir(copyDirTar,0755)
	if err != nil{
		log.Println(err)
	}
	err = os.Mkdir(copyFileTar,0755)
	if err != nil{
		log.Println(err)
	}

	////最后是否要删除
	err =  os.Remove(copyDirTar)
	if err != nil{
		log.Println(err)
	}
	err =  os.Remove(copyFileTar)
	if err != nil{
		log.Println(err)
	}
}

func genarateFile(content []byte) {

	destFileDir = "log/gen-" + currentTimeFormat() + ".log"
	//destFileDir = "log/gen-"+currentTimeFormat()+".log"

	f, err := os.OpenFile(destFileDir, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0755)

	defer f.Close()

	//log.Println(destFileDir)

	if err != nil {
		log.Fatal("生成目标文件异常:", err)
	}

	//添加换行
	c := string(content) + "\n"
	f.Write([]byte(c))

}

func timeDeltaAndDeviceIdOK(lineLog []byte) bool {

	logLine := string(lineLog)
	//将condition 按照 空格拆分
	conds := strings.Split(strings.Trim(deviceId, " "), " ")
	//log.Println(conds)
	tip := true
	for _, content := range conds {
		if !strings.Contains(logLine, content) {
			tip = false
			break
		}
	}
	//满足全部条件
	if tip {
		genarateFile(lineLog)
	}

	return tip
}

func findTextInFile(fileName string) {

	f, err := os.Open(fileName)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	reader := bufio.NewReader(f)

	for {
		line, prefix, err := reader.ReadLine()

		if err != nil {
			if err == io.EOF {
				log.Println("处理文件 压缩文件 结束了,ok !")
				break
				////通知可以压缩文件了
				//gzipOK <- struct{}{}
				//go gzipFile()
				//<-end
				//return
			} else {
				log.Fatal(line, prefix, err)
			}
		}
		//log.Println(string(line),prefix)
		timeDeltaAndDeviceIdOK(line)
	}
}


func getDestDir() string {

	return util.GetCurrentDirectory() + "/"+ copyDirTar
}


func getDestFileDir() string {

	return util.GetCurrentDirectory() +"/"+ copyFileTar[2:] + "/"
}

func currentTimeFormat() string {
	//默认当天
	return time.Now().Format(TIMEFORMAT)
}


// 生成的压缩文件 便于区分
func currentTimeFormatZip() string {
	//默认当天
	return time.Now().Format(TIMEFORMATZIP)
}

