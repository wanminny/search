package data

import (
	"net/http"
	"github.com/julienschmidt/httprouter"
	"fmt"
	"log"
	"io/ioutil"
	"gobible/logmanager/cli/http/models/data"
	"encoding/json"
	"time"
	"gobible/logmanager/cli/http/services/search"
)


func Search(res http.ResponseWriter,req *http.Request,params httprouter.Params)  {
	res.Write([]byte("search service!"))
}


func getGlobalDirsName(dirs []string,startTime,endTime string)  {

	ts, err := time.Parse(search.TIMEFORMAT, startTime)
	if err != nil {
		log.Fatal("解析开始时间格式错误:", err, startTime)
	}
	te, err := time.Parse(search.TIMEFORMAT, endTime)
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
			dirs = append(dirs, ts.Format(search.TIMEFORMAT))
			ts = ts.Add(time.Hour * 24)
		}
	}
	log.Println(dirs)
}

//获取数据服务
func Pick(res http.ResponseWriter,req *http.Request,params httprouter.Params)  {

	//startTime := req.PostFormValue("start")
	//endTime := req.PostFormValue("end")
	//condition := req.PostFormValue("c")
	//dir := req.PostFormValue("dir")

	////startTime := params.ByName("start")
	//log.Println(startTime,endTime,condition,dir)

	content,err := ioutil.ReadAll(req.Body)
	if err !=nil{
		log.Println(err)
	}
	defer req.Body.Close()
	log.Println("body :",string(content))

	pickData := data.Pick{}

	err = json.Unmarshal(content,&pickData)
	if err != nil{
		log.Println(err)
	}
	//log.Println(pickData)
	fmt.Fprint(res,string(content))

	startTime := pickData.Start
	endTime := pickData.End
	condition := pickData.C
	//自定义要查找的目录
	dir := pickData.Dir

	//参数校验
	if len(startTime) == 0 ||
		len(endTime) ==0 ||
		len(condition) == 0 ||
		len(dir) == 0 {
		fmt.Fprint(res,"参数不合法！")
		return
	}
	//格式化的日志列表slice
	dirs := make([]string,0)
	//参数校验
	getGlobalDirsName(dirs,startTime,endTime)

	search.DoSearch(dirs,dir)
	fmt.Fprintf(res,"文件较大,处理中！")
}