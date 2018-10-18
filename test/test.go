package main

import (
	"os"
	"bufio"
	"io"
	"log"
	"strings"
)

var (
	destFileDir = ""

	deviceId = "Device"
)

func genarateFile(content []byte) {

	destFileDir = "/gopath/src/gobible/logmanager/cli/log/gen-11"  + ".log"
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
	log.Println(logLine,conds)
	//os.Exit(0)
	tip := true
	for _, content := range conds {
		if !strings.Contains(logLine, content) {
			//log.Println(222)
			tip = false
			break
		}
	}
	//满足全部条件
	log.Println(555,tip)
	if tip {
		log.Println(logLine)
		genarateFile(lineLog)
	}

	return tip
}


func main()  {


	//fullFilename := "/Users/itfanr/Documents/test.txt"
	//fmt.Println("fullFilename =", fullFilename)
	//var filenameWithSuffix string
	//filenameWithSuffix = path.Base(fullFilename) //获取文件名带后缀
	//fmt.Println("filenameWithSuffix =", filenameWithSuffix)
	//
	//
	//
	//util.Copy("/gopath/src/gobible/logmanager/cli/my/xx.txt","/gopath/src/gobible/logmanager/cli/my/tcp_stdout.log-20181015")


	realName := "/gopath/src/gobible/logmanager/cli/tcp_stdout.log-2018101512"

	f, err := os.Open(realName)
	if err != nil {
		log.Println(err)
	}
	reader := bufio.NewReader(f)
	n := 0
	for {
		line, prefix, err := reader.ReadLine()
		n++

		//if n == 5 {
		//	log.Println("quit.")
		//	os.Exit(1)
		//}
		log.Println(string(line))

		if err != nil {
			if err == io.EOF {
				log.Println("处理文件结束了,ok !")
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
	log.Println(n)

}

