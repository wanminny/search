package main

import (
	"os"
	"github.com/sirupsen/logrus"
	"bufio"
	"io"
	"strings"
	"time"
	"path/filepath"
	"log"
	"encoding/json"
)

const TIMEFORMAT = "20060102150405"

var (
	//生成的目标普通文件(.log)
	destFileDir string

	//中间的临时生成的日志目录
	tmpLogDir = ""

	line map[string]int = make(map[string]int)

)

func GetCurrentDirectory() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0])) //返回绝对路径  filepath.Dir(os.Args[0])去除最后一个元素的路径
	if err != nil {
		logrus.Fatal(err)
	}
	return strings.Replace(dir, "\\", "/", -1) //将\替换成/
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

		urlStart := strings.Index(logLine,"GET")
		urlEnd := strings.Index(logLine,"?")
		url := logLine[urlStart:urlEnd]
		//log.Println(url)
		//os.Exit(1)
		line[url]++

		//genarateFile(lineLog,deviceId)
	}

	return tip
}

func genUrlFile(url map[string]int)  {

	destFileDir := "url" + genFileTimeFormat() + ".txt"

	f, err := os.OpenFile(destFileDir, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0755)

	if err != nil{
		log.Println(err)
	}

	defer f.Close()

	contentByte,err := json.Marshal(url)

	if err != nil{
		log.Println(err)
	}

	f.Write(contentByte)
}


func findTextInFile(fileName,deviceId string) {

	f, err := os.Open(fileName)
	if err != nil {
		logrus.Fatal(err)
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
				logrus.Fatal(line, prefix, err)
			}
		}
		timeDeltaAndDeviceIdOK(line,deviceId)
	}
}

func genFileTimeFormat() string  {
	//默认当天
	return time.Now().Format(TIMEFORMAT)
}

func genarateFile(content []byte,deviceId string) {

	tmpLogDir = GetCurrentDirectory()
	destFileDir = tmpLogDir + "/" + "gen_"  + genFileTimeFormat() + ".log"

	f, err := os.OpenFile(destFileDir, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0755)

	defer f.Close()

	if err != nil {
		logrus.Fatal("生成目标文件异常:", err)
	}

	//添加换行
	c := string(content) + "\n"
	f.Write([]byte(c))

}

//GET /api/device/gpsInfo/wgs84?imeiList=1453860061,1453852979&startTime=

func main()  {

	fileName:= os.Args[2]

	condition := os.Args[1]

	findTextInFile(fileName,condition)

	//最后生成文件
	genUrlFile(line)

}
