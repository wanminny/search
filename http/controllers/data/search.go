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
	"gobible/logmanager/cli/http/cache/redis"
	"gobible/logmanager/cli/http/utils"
)


func Search(res http.ResponseWriter,req *http.Request,params httprouter.Params)  {
	res.Write([]byte("search service!"))
}


func errResultJson(res http.ResponseWriter,msg string,err error)  {

	errMsg := ""
	if err != nil{
		errMsg = err.Error()
	}
	content := fmt.Sprintf("%s,%s",msg,errMsg)
	rlt := data.NewJson(1,content,nil)
	fmt.Fprint(res,string(rlt))
	//return
}

func getGlobalDirsName(res http.ResponseWriter,dirs []string,startTime,endTime string)  (err error)  {

	ts, err := time.Parse(search.TIMEFORMAT, startTime)
	if err != nil {
		errResultJson(res,"解析开始时间格式错误:",nil)
		return err
		//log.Fatal("解析开始时间格式错误:", err, startTime)
	}
	te, err := time.Parse(search.TIMEFORMAT, endTime)
	if err != nil {
		errResultJson(res,"解析结束时间格式错误:",nil)
		return err
		//log.Fatal("解析结束时间格式错误:", err, endTime)
	}
	if te.Before(ts) {
		errResultJson(res,"日期不合法,结束日期比开始日期还早哦:",nil)
		return err
		//log.Fatal("日期不合法,结束日期比开始日期还早哦.")
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
	return
}

func isProcessing(key string)  bool {

	value,err := redis.GetValue(key)
	if err != nil{  // nil returned 表示key不存在
		log.Println(err)
		return false
	}
	if len(value) >0 {
		return true
	}
	return false
}


//获取数据服务
func Pick(res http.ResponseWriter,req *http.Request,params httprouter.Params)  {

	//startTime := req.PostFormValue("start")
	////startTime := params.ByName("start")
	//log.Println(startTime,endTime,condition,dir)


	content,err := ioutil.ReadAll(req.Body)
	if err !=nil{
		log.Println(err)
	}
	defer req.Body.Close()
	//log.Println("body :",string(content))

	pickData := data.Pick{}

	err = json.Unmarshal(content,&pickData)
	if err != nil{
		log.Println(err)
	}
	//log.Println(pickData)
	//fmt.Fprint(res,string(content))

	startTime := pickData.Start
	endTime := pickData.End
	condition := pickData.C
	//自定义要查找的目录
	dir := pickData.Dir

	//参数校验
	if len(startTime) == 0 ||
		len(endTime) == 0 ||
		len(condition) == 0 ||
		len(dir) == 0 {
		rlt := data.NewJson(1,"参数不合法",nil)
		fmt.Fprint(res,string(rlt))
		return
	}

	//如果参数合法就判断是否是重复的请求
	composeStr := fmt.Sprintf("%s-%s-%s-%s",startTime,endTime,condition,dir)
	findCondition := utils.MD5(composeStr)
	if isProcessing(findCondition) {
		rlt := data.NewJson(0,"该查找任务已经提交,正在处理中,请问重复提交,谢谢！",nil)
		fmt.Fprint(res,string(rlt))
		return
	}

	//格式化的日志列表slice
	dirs := make([]string,0)
	//参数校验
	err = getGlobalDirsName(res,dirs,startTime,endTime)
	if err != nil{
		return
	}


	//os.Exit(0)
	//go search.DoSearch(dirs,dir)

	//提交任务后马上设置值
	redis.SetValue(findCondition,composeStr)

	rlt := data.NewJson(0,"文件处理中",nil)
	fmt.Fprint(res,string(rlt))
	return

}