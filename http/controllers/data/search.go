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
	"errors"
	"github.com/sirupsen/logrus"
	"gobible/logmanager/cli/util"
	"os"
	"gobible/logmanager/cli/http/config"
	"gobible/logmanager/cli/http/controllers/file"
	"strings"
)

func Search(res http.ResponseWriter,req *http.Request,params httprouter.Params)  {

	bytes,err := json.Marshal(file.Test)
	if err != nil{
		log.Println(err)
	}
	fmt.Fprint(res,string(bytes))
	return
}


func errResultJson(res http.ResponseWriter,msg string,err error)  {

	errMsg := ""
	if err != nil{
		errMsg = err.Error()
	}
	content := fmt.Sprintf("%s %s",msg,errMsg)
	rlt := data.NewJson(1,content,nil)
	fmt.Fprint(res,string(rlt))
	//return
}

func getGlobalDirsName(res http.ResponseWriter,dirs *[]string,startTime,endTime string)  (err error)  {

	ts, err := time.Parse(search.TIMEFORMAT, startTime)
	if err != nil {
		errResultJson(res,"解析开始时间格式错误:",nil)
		return err
	}
	te, err := time.Parse(search.TIMEFORMAT, endTime)
	if err != nil {
		errResultJson(res,"解析结束时间格式错误:",nil)
		return err
	}
	if te.Before(ts) {
		err = errors.New("")
		errResultJson(res,"日期不合法,结束日期比开始日期还早哦:",nil)
		return err
	}
	if ts.Equal(te) { // 日期相等
		*dirs = append(*dirs, startTime)
	} else {
		// 日期大于前者
		*dirs = append(*dirs, startTime)
		ts = ts.Add(time.Hour * 24)
		for te.After(ts) || te.Equal(ts){
			*dirs = append(*dirs, ts.Format(search.TIMEFORMAT))
			ts = ts.Add(time.Hour * 24)
		}
	}
	logrus.Println(*dirs)
	return
}

func isProcessing(key string)  bool {

	value,err := redis.GetValue(key)
	if err != nil{  // nil returned 表示key不存在
		//log.Println(err)
		return false
	}
	if len(value) >0 {
		return true
	}
	return false
}

func genDownloadDirIfInputEmpty()  {

	//初始化服务器日志生成的目录；便于查看情况
	if !util.PathExist(util.GetCurrentDirectory() + "/" + config.ZipResultDir){
		err :=os.Mkdir(config.ZipResultDir,0755)
		if err != nil{
			logrus.Println(err)
		}
	}

}

//获取数据服务
func PickDemo(res http.ResponseWriter,req *http.Request,params httprouter.Params)  {

	//data.SetCROS(res)
	res.Header().Add("Access-Control-Allow-Headers", "Content-Type") //header的类型
	//startTime := req.PostFormValue("start")
	////startTime := params.ByName("start")

	content,err := ioutil.ReadAll(req.Body)
	if err !=nil{
		logrus.Println(err)
	}
	defer req.Body.Close()
	//log.Println("body :",string(content))

	pickData := data.Pick{}

	err = json.Unmarshal(content,&pickData)
	if err != nil{
		logrus.Println(err)
	}

	startTime := pickData.Start
	endTime := pickData.End
	condition := pickData.C
	//自定义要查找的目录
	dir := pickData.Dir
	down := pickData.Down

	configFile := redis.CheckRedisConf()
	configDown := configFile.Down

	if len(configDown) == 0 {
		rlt := data.NewJson(1,"获取config.json配置失败",nil)
		fmt.Fprint(res,string(rlt))
		return
	}

	down = configDown
	if len(down) == 0 {
		rlt := data.NewJson(1,"config.json down 配置为空",nil)
		fmt.Fprint(res,string(rlt))
		return
	}

	//参数校验
	if len(startTime) == 0 ||
		len(endTime) == 0 ||
		len(condition) == 0 ||
		len(dir) == 0 {
		rlt := data.NewJson(1,"参数不合法",nil)
		fmt.Fprint(res,string(rlt))
		return
	}

	if len(dir) > 0 {
		if !util.PathExist(dir){
			rlt := data.NewJson(1,"查找路径不存在",nil)
			res.Write([]byte(string(rlt)))
			return
		}
	}

	if len(down) > 0 {
		if !util.PathExist(down){
			rlt := data.NewJson(1,"下载目录不存在",nil)
			res.Write([]byte(string(rlt)))
			return
		}
	}

	//如果没有down参数则是download目录
	if len(down) == 0 {
		genDownloadDirIfInputEmpty()
	}else{
		//search.ZipDirSignal <- down
	}

	//如果参数合法就判断是否是重复的请求
	composeStr := fmt.Sprintf("%s-%s-%s-%s",startTime,endTime,condition,dir)
	findCondition := utils.MD5(composeStr)
	if isProcessing(findCondition) {
		rlt := data.NewJson(0,"该查找任务已经提交,正在处理中,请问重复提交,谢谢！",nil)
		res.Write([]byte(rlt))
		return
	}

	//格式化的日志列表slice
	dirs := make([]string,0)
	//参数校验
	err = getGlobalDirsName(res,&dirs,startTime,endTime)
	if err != nil{
		return
	}

	go search.DoSearch(dirs,dir,findCondition,condition,down)

	//提交任务后马上设置值
	redis.SetValue(findCondition,composeStr)
	rlt := data.NewJson(0,"文件处理中",nil)
	res.Write([]byte(rlt))

}



//获取数据服务
func Pick(res http.ResponseWriter,req *http.Request,params httprouter.Params)  {

	res.Header().Add("Access-Control-Allow-Headers", "Content-Type") //header的类型
	content,err := ioutil.ReadAll(req.Body)
	if err !=nil{
		logrus.Println(err)
	}
	defer req.Body.Close()
	pickData := data.Pick{}

	err = json.Unmarshal(content,&pickData)
	if err != nil{
		logrus.Println(err)
	}

	startTime := pickData.Start
	endTime := pickData.End
	condition := pickData.C
	//自定义要查找的目录
	dir := pickData.Dir
	//down := pickData.Down

	//参数校验
	if len(startTime) == 0 ||
		len(endTime) == 0 ||
		len(condition) == 0 ||
		len(dir) == 0 {
		rlt := data.NewJson(1,"参数不合法",nil)
		fmt.Fprint(res,string(rlt))
		return
	}
	if len(dir) > 0 {
		if !util.PathExist(dir){
			rlt := data.NewJson(1,"查找路径不存在",nil)
			res.Write([]byte(string(rlt)))
			return
		}
	}

	//如果参数合法就判断是否是重复的请求
	composeStr := fmt.Sprintf("%s-%s-%s-%s",startTime,endTime,condition,dir)
	findCondition := utils.MD5(composeStr)

	//if isProcessing(findCondition) {
	//	rlt := data.NewJson(0,"该查找任务已经提交,正在处理中,请问重复提交,谢谢！",nil)
	//	res.Write([]byte(rlt))
	//	return
	//}

	//格式化的日志列表slice
	dirs := make([]string,0)
	//参数校验
	err = getGlobalDirsName(res,&dirs,startTime,endTime)
	if err != nil{
		return
	}

	// 如果任务存在着不需要再下发了;并显示状态 ？ 或者干脆不管！


	err = redis.LPush(config.RedisTaskName,findCondition)
	if err != nil{
		rlt := data.NewJson(1,"下发任务失败",err.Error())
		res.Write([]byte(string(rlt)))
		return
	}
	task := map[string] interface{} {
		"status":config.RedisStatusNotStart,
		"condition":composeStr,
	}
	redis.HMSet(findCondition,task)
	msg := fmt.Sprintf("%s,%s",findCondition,"任务下发成功;请去查询接口获取结果文件地址")
	rlt := data.NewJson(1,msg, struct {}{})
	res.Write([]byte(string(rlt)))
	return
}



func DoWork()  {

	for {

		listLen,err := redis.LLen(config.RedisTaskName)
		if err != nil{
			log.Println(err)
			continue
		}
		if listLen > 0 {

			v,err := redis.RRop(config.RedisTaskName)
			if err != nil{
				log.Println(err)
				continue
			}
			//根据v查找 hash中的条件和状态
			job,err := redis.HGetAll(v)
			if err != nil{
				log.Println(err)
				continue
			}

			//status := job.Status
			condtion := job.Condition

			fields := strings.Split(condtion,config.ConditionSplitChar)
			//格式化的日志列表slice
			dirs := make([]string,0)
			startTime := fields[0]
			endTime := fields[1]

			ts, err := time.Parse(search.TIMEFORMAT, startTime)
			if err != nil {
				log.Println(err)
			}
			te, err := time.Parse(search.TIMEFORMAT, endTime)
			if err != nil {
				log.Println(err)
			}
			if ts.Equal(te) { // 日期相等
				dirs = append(dirs, startTime)
			} else {
				// 日期大于前者
				dirs = append(dirs, startTime)
				ts = ts.Add(time.Hour * 24)
				for te.After(ts) || te.Equal(ts){
					dirs = append(dirs, ts.Format(search.TIMEFORMAT))
					ts = ts.Add(time.Hour * 24)
				}
			}
			logrus.Println(dirs)

			//findCondition := ""
			condition := fields[2]
			dir :=  fields[3]
			down := search.DownloadDir
			//获取到所有的值
			search.DoSearch(dirs,dir,"",condition,down)

			//成功后设置标志位

			//提交任务后马上设置值
			//redis.SetValue(findCondition,composeStr)
			//rlt := data.NewJson(0,"文件处理中",nil)
			//res.Write([]byte(rlt))
		}

	}

}