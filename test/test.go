package main

import (
	"os"
	"bufio"
	"io"
	"log"
	"strings"
	"fmt"
	"path"
	"crypto/md5"
	"encoding/hex"
	"flag"
	"time"
	"runtime"
	"reflect"
	"net/http"
	"strconv"
	"encoding/json"
	"net/http/httptest"
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

func mkdirs()  {

	//newpath := filepath.Join(".", "copy-dir-tar")
	//log.Println(newpath)
	//os.MkdirAll(newpath, 0755)

	err := os.Remove("./copy-dir-tar")

	//err := os.Mkdir("./copy-dir-tar",0755)

	log.Println(err)
	//os.Mkdir("./copy-file-tar",0755)

}

func MD5(text string) string{
	ctx := md5.New()
	ctx.Write([]byte(text))
	return hex.EncodeToString(ctx.Sum(nil))
}

var (
	//当前的目录的拷贝目录
	copyDirTar = "copy-dir-tar"

	//当前的文件的拷贝目录
	copyFileTar = "copy-file-tar"
)

func mkdirs1()  {

	//先删除 后创建
	err :=  os.RemoveAll(copyDirTar)
	if err != nil{
		log.Println(err)
	}
	err =  os.RemoveAll(copyFileTar)
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

	//最后是否要删除？
	err =  os.RemoveAll(copyDirTar)
	if err != nil{
		log.Println(err)
	}
	err =  os.RemoveAll(copyFileTar)
	if err != nil{
		log.Println(err)
	}
}

func demo()  {

	var (
		a =1
	)
	defer func() {
		log.Println(11111)
	}()

	if a == 1{
		log.Println(55555)
		return
	}

}

func testSplit()  {
	a := (strings.Split("20181018-20181021-abc daf-/a/bc/d/asdf","-"))
	log.Println(a[2])
	os.Exit(1)
}

func testflag()  {
	var name = ""
	flag.StringVar(&name,"name","1","name usage")
	flag.Parse()

	log.Println(flag.NArg(),flag.NFlag())

	//io.readF
	//f,err := os.Open("aaa.txg")
	//ioutil.ReadAll()

	y,m,d := time.Now().Date()
	log.Println(y,int(m),d)

	funcName := func() {}
	func2 := log.Println
	//time.ParseDuration()
	//time.


	//runtime.

	//runtime.Func{}
	//获取函数名称:
	log.Println(runtime.FuncForPC(reflect.ValueOf(funcName).Pointer()).Name())
	funPc := (runtime.FuncForPC(reflect.ValueOf(func2).Pointer()))

	log.Println(funPc)
	log.Println(funPc.Entry())


	flag.PrintDefaults()

}

func teststd()  {

	//var a FloatType = 1.1
	//FloatT
	//fmt.Errorf()

}

func testParseInt()  {
	i, err := strconv.ParseInt("123", 8, 0)
	if err != nil {
		panic(err)
	}
	println(i)

}
func main()  {

	h(nil,nil)
	//testParseInt()

	//testflag()
	//demo()
	//testSplit()
}

func h(w http.ResponseWriter,r *http.Request)  {
	//r.Cookies()
	//r.Cookie("")
	//http.SetCookie(w,c)
	//strconv.Itoa()
	//strconv.Atoi()
	//strconv.ParseInt()
	a := struct{
		A string
		D string
	}{A:"ASdf",D:"wwwagagmmm"}
	nby,_ := json.MarshalIndent(a,"","\t")
	log.Println(string(nby))

	httptest.NewRecorder()
	//httptest.NewRequest()
}

func main1()  {


	log.Println(strings.Split("20181018-20181021-abc daf-/a/bc/d/asdf","-"))
	os.Exit(1)

	//destDirFile := "/gopath/src/gobible/logmanager/cli/test/ac.txt"

	destDirFile := "/gopath/src/gobible/logmanager/cli/test/cc/ac.txt"

	f,err := os.OpenFile(destDirFile, os.O_CREATE | os.O_WRONLY, 0644)
	if err != nil{
		log.Println(err)
	}
	os.Exit(0)

	ZipResultDir := "/ab/b/c/d"
	length := len(ZipResultDir)
	if ZipResultDir[length-1:] == "/"{
		log.Println(1111)
		ZipResultDir = ZipResultDir[:length-1]
		log.Println(ZipResultDir)
	}

	os.Exit(0)
	mkdirs1()

	os.Exit(1)

	//log.Println(len(nil))
	log.Println(MD5("startendcondir"))

	os.Exit(1)

	mkdirs()

	os.Exit(1)

	fullFilename := "/Users/itfanr/Documents/test.txt"
	fmt.Println("fullFilename =", fullFilename)
	var filenameWithSuffix string
	filenameWithSuffix = path.Base(fullFilename) //获取文件名带后缀
	fmt.Println("filenameWithSuffix =", filenameWithSuffix)
	//
	//
	os.Exit(0)

	//
	//util.Copy("/gopath/src/gobible/logmanager/cli/my/xx.txt","/gopath/src/gobible/logmanager/cli/my/tcp_stdout.log-20181015")


	realName := "/gopath/src/gobible/logmanager/cli/tcp_stdout.log-2018101512"

	f, err = os.Open(realName)
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

