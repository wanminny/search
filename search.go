package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"time"
	"io/ioutil"
	"path"
	"gobible/logmanager/cli/util"
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
	copyDirTar = "copy-dir-tar"

	//当前的文件的拷贝目录
	copyFileTar = "copy-file-tar"

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

func genarateFile(content []byte) {

	keyWords := deviceId
	destFileDir = "log/gen-" + keyWords + currentTimeFormat() + ".log"
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

//压缩文件
func gzipFile() {
	<-gzipOK

	//destGzipFileDir = "/tar/" + currentTimeFormat() + ".tar.gz"
	destGzipFileDir = deviceId + currentTimeFormatZip()+".zip"

	//Compress(destGzipFileDir,destFileDir)
	//log.Println(GetCurrentDirectory())

	wantZipDir := util.GetCurrentDirectory() + "/log"

	//压缩文件
	ZipDir(wantZipDir, destGzipFileDir)

	//结束程序 信号
	end <- 1
}

func init() {
	//log.SetFlags(log.Llongfile | log.Ltime)
}

func initParameter() {

	flag.BoolVar(&help, "h", false, "show help:")

	flag.StringVar(&startTime, "start", "", "format as: 2016-01-02,默认是当天.")

	flag.StringVar(&endTime, "end", "", "format as: 2016-01-02,默认当前时间的48h前.")

	flag.StringVar(&deviceId, "c", "", "which condition to search.")

	flag.StringVar(&directory, "dir", "", "如果文件不存在需要复制的目录(解压文件)")

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

func endTimeFormat() string {

	now := time.Now()
	endT, err := time.ParseDuration(delta)
	if err != nil {
		log.Fatal(err)
	}
	return now.Add(endT).Format(TIMEFORMAT)
}

//格式化使用方式：
func usage() {

	fmt.Fprint(os.Stdout, `
search version : 0.0.1 
Usage of ./search: -c  "RFID:3d xxxx " -start "20180925" -end "20180926" -dir /gopath/src/gobible/logmanager/cli
  -c string
    	which condition to search 必填参数.
  -dir string
    	 如果文件不存在需要复制的目录(来解压文件).
  -end string
    	format as: 20160102,默认当前时间的48h前. 必填参数 `+
		`
  -h	show help: (default true)
  -start string
    	format as: 20160102,必填参数. 
`)
	fmt.Println()
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

func parameterCheck()  {

	if len(deviceId) == 0 {
		usage()
		log.Fatal("缺少查找条件参数")
	}
	if len(startTime) == 0 {
		usage()
		log.Fatal("缺少开始时间参数")
	}
	if len(endTime) == 0 {
		usage()
		log.Fatal("缺少结束时间参数")
	}

	if len(directory) == 0 {
		usage()
		log.Fatal("没有指定要查找的目录.")
	}

}

func timeCheckAndStoreDirs()  {

	dirs := make([]string, 0)
	ts, err := time.Parse(TIMEFORMAT, startTime)
	if err != nil {
		log.Fatal("解析开始时间格式错误:", err, startTime)
	}
	te, err := time.Parse(TIMEFORMAT, endTime)
	if err != nil {
		log.Fatal("解析结束时间格式错误:", err, endTime)
	}
	if te.Before(ts) {
		log.Fatal("日期不合法,结束日期比开始日期还早哦.")
	}
	if ts.Equal(te) { // 日期相等
		dirs = append(dirs, startTime)
	} else {
		// 日期大于前者
		dirs = append(dirs, startTime)
		//log.Println(dirs)
		ts = ts.Add(time.Hour * 24)
		for te.After(ts) || te.Equal(ts){
			dirs = append(dirs, ts.Format(TIMEFORMAT))
			ts = ts.Add(time.Hour * 24)
		}
	}
	log.Println(dirs)
}

func mkdirs()  {

	//先删除 后创建
	err :=  os.Remove(copyDirTar)
	if err != nil{
		log.Println(err)
	}
	err =  os.Remove(copyFileTar)
	if err != nil{
		log.Println(err)
	}

	err =os.Mkdir(copyDirTar,0755)
	if err != nil{
		log.Println(err)
	}
	err = os.Mkdir(copyFileTar,0755)
	if err != nil{
		log.Println(err)
	}

	////最后是否要删除？
	//err =  os.Remove(copyDirTar)
	//if err != nil{
	//	log.Println(err)
	//}
	//err =  os.Remove(copyFileTar)
	//if err != nil{
	//	log.Println(err)
	//}


}

func getGlobalDirsName()  {

	ts, err := time.Parse(TIMEFORMAT, startTime)
	if err != nil {
		log.Fatal("解析开始时间格式错误:", err, startTime)
	}
	te, err := time.Parse(TIMEFORMAT, endTime)
	if err != nil {
		log.Fatal("解析结束时间格式错误:", err, endTime)
	}
	if te.Before(ts) {
		log.Fatal("日期不合法,结束日期比开始日期还早哦.")
	}
	if ts.Equal(te) { // 日期相等
		dirs = append(dirs, startTime)
	} else {
		// 日期大于前者
		dirs = append(dirs, startTime)
		//log.Println(dirs)
		ts = ts.Add(time.Hour * 24)
		for te.After(ts) || te.Equal(ts){
			dirs = append(dirs, ts.Format(TIMEFORMAT))
			ts = ts.Add(time.Hour * 24)
		}
	}
	log.Println(dirs)
}


func getDestDir() string {

	//return util.GetCurrentDirectory() + "/"+ copyDirTar
	return util.GetCurrentDirectory() + "/" + copyDirTar + "/"

}


func getDestFileDir() string {

	return util.GetCurrentDirectory() +"/"+ copyFileTar + "/"
}



func main() {

	initParameter()
	flag.Parse()
	if help {
		usage()
		return
	}
	parameterCheck()
	getGlobalDirsName()
	mkdirs()
	for _, dirv := range dirs {
		realName := util.GetCurrentDirectory() + "/" + prefix + dirv
		realNameGz := util.GetCurrentDirectory() + "/" + prefix + dirv + extName
		//log.Println(realName, realNameGz, util.PathExist(realName))

		//如果不带扩展名gz的文件存在 认为该文件在当前目录 【如果同时存在既有压缩的又有不压缩的；则以不压缩的为准；gz的忽略
		// 如果没有非压缩的文件 则走下面的流程 以有压缩的为准】
		if (util.PathExist(realName)){
			//dirsMap[dirv] =  true
			fileMaps[dirv] = struct {
				Name     string
				gz       bool
				fullName string
				empty    bool
			}{Name: dirv, gz:false , fullName: realName, empty:false }
			continue
		}

		// 如果不带扩展名gz的文件不存在；看是否本地有带gz的文件存在如果存在后续就不需要在重复处理了；
		if util.PathExist(realNameGz) {
			//dirsMap[dirv] =  true
			fileMaps[dirv] = struct {
				Name     string
				gz       bool
				fullName string
				empty    bool
			}{Name: dirv, gz: true, fullName: realNameGz, empty:false }
			continue
		}else{

			//既没有非压缩的也没有压缩的 情况 ===>  需要到指定目录去copy unzip and search ;
			//dirsMap[dirv] =  false
			fileMaps[dirv] = struct {
				Name     string
				gz       bool
				fullName string
				empty    bool
			}{empty:true}
		}
	}

	//fmt.Printf("%#v\n",fileMaps)
	//os.Exit(1)

	// 获取 不满足 条件的slice
	unsatisfy := make([]string,0)
	for k,v := range fileMaps{
		if v.empty {
			unsatisfy = append(unsatisfy,k)
		}
	}
	fmt.Printf("%#v",unsatisfy)

	for k,v := range fileMaps{
		//先处理本地的所有的已知的文件；【包括两种情况 1.是有没有扩展名称的； 2.一种是有扩展名称的】
		realNameIte := util.GetCurrentDirectory() + "/" + prefix + k
		realNameGzIt := util.GetCurrentDirectory() + "/" + prefix + k + extName
		if !v.empty{
			//一种情况是 非压缩
			if !v.gz {
				//文本文件
				f, err := os.Open(realNameIte)
				if err != nil {
					//文件已经存在
					log.Println(err)
				}
				reader := bufio.NewReader(f)
				for {
					line, prefix, err := reader.ReadLine()

					if err != nil {
						if err == io.EOF {
							log.Println("非压缩 处理文件结束了,ok !")
							break
							//通知可以压缩文件了
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
			} else{   ////一种情况是压缩 v.gz == true

				//先拷贝在解压;然后在查找；
				destFile := util.GetCurrentDirectory() +"/"+ copyFileTar + "/" + util.GetFileName(realNameGzIt)+ ".gz"
				util.SimpleCopyFile(destFile,realNameGzIt)

				UnGzipFile(destFile,util.GetCurrentDirectory() +"/"+ copyFileTar +"/" + util.GetFileName(realNameGzIt)) //xxx.gz
				//log.Println(fullName,util.GetFileName(filenameFullName),99)
				findTextInFile(getDestFileDir() + util.GetFileName(realNameGzIt))
			}
		}else{

		}
	}

	// 先把前面满足条件的一次性跑完；
	// 这里专门处理不满足条件的；  即这里是所有结构体为空的情况empty= true
	for _,v := range unsatisfy{

		realNameGzIt := util.GetCurrentDirectory() + "/" + prefix + v + extName

		//是当前目录没有的；需要去指定目录 处理的
		log.Println("开始===>去指定目录中查找.")
		//如果文件不存在 则复制指定目录的文件过来让后变量素有文件
		if len(directory) == 0 {
			usage()
			log.Fatal("没有指定要查找的目录.")
		}else{
			//直接将制定目录的.gz文件解压到指定文件然后查找处理
			//遍历所有的目录
			util.Copy(directory,getDestDir())
			log.Println(getDestDir())
			dirv := getDestDir()
			files, err := ioutil.ReadDir(dirv)
			if err != nil {
				log.Fatal(err)
				//continue
			}
			for _, v := range files {
				filenameFullName := path.Base(v.Name())
				fullName := dirv + v.Name()
				ext := path.Ext(filenameFullName)

				if ext == extName {
					//文件名称是满足格式的压缩文件才需要处理
					log.Println(util.GetFileName(filenameFullName),44)
					inSliceFileName := util.GetFileName(filenameFullName)[len(prefix):]
					if ok,err :=util.Contain(inSliceFileName,unsatisfy); err != nil {
						log.Println(ok,err)
						continue
					}
					log.Println(filenameFullName,ext,fullName,66)
					//先解压文件；
					UnGzipFile(fullName,util.GetCurrentDirectory() + "/" + copyDirTar + "/" + util.GetFileName(realNameGzIt)) //xxx.gz

					//UnGzipFile(fullName,getDestDir() + util.GetFileName(realNameGzIt)) //xxx.gz
					//log.Println(fullName,util.GetFileName(filenameFullName),99)
					log.Println(fullName,util.GetCurrentDirectory() + "/" + copyDirTar + "/" + util.GetFileName(realNameGzIt),77)
					log.Println(dirv + util.GetFileName(filenameFullName),55)
					findTextInFile(dirv + util.GetFileName(filenameFullName))
				}
			}
		}
	}

	gzipOK <- struct{}{}
	go gzipFile()
	<-end

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
