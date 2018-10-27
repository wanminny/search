package search

import (
	"os"
	"bufio"
	"io"
	"strings"
	"time"
	"gobible/logmanager/cli/http/services/util"
	"strconv"
	"math/rand"
	"github.com/sirupsen/logrus"
	"gobible/logmanager/cli/http/config"
	"log"
	"gobible/logmanager/cli/http/cache/redis"
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
	gzipOK  = make(chan struct{}, 1)

	//end
	end  = make(chan int)

	// 需要解压的文件的扩展名
	extName = ".gz"

	//要查找的文件名称的前缀
	prefix = "tcp_stdout.log-"

	//根据开始 结束时间 返回名称
	dirs = make([]string, 0)

	dirsMap = make(map[string]bool)

	//中间的临时生成的日志目录
	tmpLogDir = "tmp-log-dir"

	//当前的目录的拷贝目录
	copyDirTar = "copy-dir-tar"

	//当前的文件的拷贝目录
	copyFileTar = "copy-file-tar"

	//download 目录
	DownloadDir = ""

	//下载域名
	DownLoadDomain = ""

	//随机数字目录
	randInt64 int64 = 9876543210

	ZipDirSignal = make(chan string)
)

type FileInfo struct {
	//文件名称目录
	Name string
	//是否是压缩文件
	gz bool

	fullName string

	empty bool
}

var (
	fileMaps  = make(map[string]FileInfo)
)

func randSeed()  {
	rand.Seed(time.Now().UnixNano())
}

func init()  {
	log.SetFlags(log.Llongfile | log.Ltime)
}

func initDir()  {

	copyDirTar =  config.TmpTransferDir + "/"+"copyDir" + genCopyDirTar()
	copyFileTar = config.TmpTransferDir + "/"+ "copyFileDir" + genCopyFileTar()
	tmpLogDir = config.TmpTransferDir + "/"+ "tmpLog" + genTmpLogDir()
	logrus.Println(copyFileTar,copyDirTar,tmpLogDir)
}


func genCopyDirTar() string {
	return strconv.FormatInt(rand.Int63n(randInt64),10)
}

func genCopyFileTar() string {
	return strconv.FormatInt(rand.Int63n(randInt64),10)
}

func genTmpLogDir() string {
	return strconv.FormatInt(rand.Int63n(randInt64),10)
}


func genarateFile(content []byte,deviceId string) {

	//防止异常
	keyWords := strings.Trim(deviceId," ")

	// linux上面如果不加全路径会导致出错？
	destFileDir = util.GetCurrentDirectory() + "/" + tmpLogDir + "/" + "gen_" + keyWords + "_" + genFileTimeFormat() + ".log"
	//destFileDir = "log/gen-"+currentTimeFormat()+".log"

	f, err := os.OpenFile(destFileDir, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0755)
	defer f.Close()

	if err != nil {
		logrus.Println("生成目标文件异常:", err)
		return
	}
	//添加换行
	c := string(content) + "\n"
	f.Write([]byte(c))

}

//压缩文件
func gzipFile(deviceId,ZipResultDir,hashKey string) {

	<-gzipOK

	//deviceId 是查找原因 之前是准备仅仅定义设备
	fileName := deviceId + "_" + currentTimeFormatZip()+".zip"

	if ZipResultDir == config.ZipResultDir{
		destGzipFileDir = util.GetCurrentDirectory() + "/" + ZipResultDir + "/" + fileName
	} else{

		length := len(ZipResultDir)
		//防止最后一个/ab/c/d/ 最后带上了/
		if ZipResultDir[length-1:] == "/"{
			ZipResultDir = ZipResultDir[:length-1]
		}
		destGzipFileDir = ZipResultDir + "/" + fileName
	}

	wantZipDir := util.GetCurrentDirectory() + "/" + tmpLogDir

	//压缩文件
	util.ZipDir(wantZipDir, destGzipFileDir)


	//成功以后需要处理 状态和下载地址
	hashValue := map[string]interface{}{
		"download": DownLoadDomain + "/" + fileName ,
		"status": config.RedisStatus(config.RedisStatusOK),
	}
	redis.HMSet(hashKey,hashValue)

	//结束程序 信号
	end <- 1
}

func genFileTimeFormat() string  {
	//默认当天
	return time.Now().Format(TIMEFORMAT)
}

// 生成的压缩文件 便于区分
func currentTimeFormatZip() string {
	//默认当天
	return time.Now().Format(TIMEFORMATZIP)
}

func timeDeltaAndDeviceIdOK(lineLog []byte,deviceId string) bool {

	logLine := string(lineLog)
	//将condition 按照 空格拆分
	conds := strings.Split(strings.Trim(deviceId, " "), " ")
	//logrus.Println(conds)
	tip := true
	for _, content := range conds {
		if !strings.Contains(logLine, content) {
			tip = false
			break
		}
	}
	//满足全部条件
	if tip {
		genarateFile(lineLog,deviceId)
	}

	return tip
}

func mkDirs()  {

	//if !util.PathExist(util.GetCurrentDirectory() + "/" + ZipResultDir){
	//	err :=os.Mkdir(ZipResultDir,0755)
	//	if err != nil{
	//		logrus.Println(err)
	//	}
	//}

	err :=os.Mkdir(copyDirTar,0755)
	if err != nil{
		logrus.Println(err)
	}
	err = os.Mkdir(copyFileTar,0755)
	if err != nil{
		logrus.Println(err)
	}

	err = os.Mkdir(tmpLogDir,0755)
	if err != nil{
		logrus.Println(err)
	}

}

func delDirs()  {

	////最后是否要删除
	err :=  os.RemoveAll(copyDirTar)
	if err != nil{
		logrus.Println(err)
	}
	err =  os.RemoveAll(copyFileTar)
	if err != nil{
		logrus.Println(err)
	}

	err =  os.RemoveAll(tmpLogDir)
	if err != nil{
		logrus.Println(err)
	}

}


func getDestFileDir() string {
	return util.GetCurrentDirectory() +"/"+ copyFileTar + "/"
}


func findTextInFile(fileName,deviceId string) {

	f, err := os.Open(fileName)
	if err != nil {
		logrus.Println("findTextInFile 打开文件异常:",err)
	}
	defer f.Close()

	reader := bufio.NewReader(f)

	for {
		line, prefix, err := reader.ReadLine()

		if err != nil {
			if err == io.EOF {
				logrus.Println("处理文件 压缩文件 结束了,ok !")
				break
			} else {
				logrus.Println("ReadLine 异常:",line, prefix, err)
			}
		}
		timeDeltaAndDeviceIdOK(line,deviceId)
	}
}