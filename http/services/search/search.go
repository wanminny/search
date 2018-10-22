package search

import (
	"fmt"
	"os"
	"bufio"
	"io"
	"log"
	"gobible/logmanager/cli/http/services/util"
	"gobible/logmanager/cli/http/cache/redis"
)

func init()  {
	log.SetFlags(log.Llongfile | log.Ltime)
}

func DoSearch(dirs []string,directory,findCondition,deviceId string)  {

	randSeed()
	mkDirs()

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

	// 获取 不满足 条件的slice
	unsatisfy := make([]string,0)
	for k,v := range fileMaps{
		if v.empty {
			unsatisfy = append(unsatisfy,k)
		}
	}

	fmt.Printf("%#v",unsatisfy)
	fmt.Printf("%#v\n",fileMaps)
	//os.Exit(1)

	for k,v := range fileMaps{
		//先处理本地的所有的已知的文件；【包括两种情况 1.是有没有扩展名称的； 2.一种是有扩展名称的】
		realNameIte := util.GetCurrentDirectory() + "/" + prefix + k
		realNameGzIt := util.GetCurrentDirectory() + "/" + prefix + k + extName
		if !v.empty{  //非空 即不需要去指定目录去查找的情况
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
					timeDeltaAndDeviceIdOK(line,deviceId)
				}
			} else{   ////一种情况是压缩 v.gz == true

				//先拷贝在解压;然后在查找；
				destFile := util.GetCurrentDirectory() +"/"+ copyFileTar + "/" + util.GetFileName(realNameGzIt)+ ".gz"
				util.SimpleCopyFile(destFile,realNameGzIt)

				util.UnGzipFile(destFile,util.GetCurrentDirectory() +"/"+ copyFileTar +"/" + util.GetFileName(realNameGzIt))

				findTextInFile(getDestFileDir() + util.GetFileName(realNameGzIt),deviceId)
			}
		}
	}

	// 先把前面满足条件的一次性跑完；
	// 这里专门处理不满足条件的；  即这里是所有结构体为空的情况empty= true 【去指定目录去查找的情况】
	for _,v := range unsatisfy{

		realNameGzIt := directory + "/" + prefix + v + extName
		//是当前目录没有的；需要去指定目录 处理的
		log.Println("开始===>去指定目录中查找.")
		//直接将制定目录的.gz文件解压到指定文件然后查找处理

		fileName := directory + "/" + util.GetFileName(realNameGzIt) + extName

		tmpDir := util.GetCurrentDirectory() + "/" + copyDirTar + "/"

		destFileName := tmpDir +  util.GetFileName(realNameGzIt) + extName

		if util.PathExist(fileName) {
			//复制文件 到指定目录
			util.SimpleCopyFile(destFileName,fileName)
			//先解压文件；
			util.UnGzipFile(destFileName,tmpDir + util.GetFileName(realNameGzIt))

			findTextInFile(tmpDir + util.GetFileName(realNameGzIt),deviceId)
		}
	}

	gzipOK <- struct{}{}
	go gzipFile(deviceId)
	<-end

	delDirs()
	//处理完成后清空key
	redis.DelKey(findCondition)

}

