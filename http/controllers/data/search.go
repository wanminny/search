package data

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"github.com/sirupsen/logrus"
	"gobible/logmanager/cli/http/cache/redis"
	"gobible/logmanager/cli/http/config"
	"gobible/logmanager/cli/http/controllers/file"
	"gobible/logmanager/cli/http/models/data"
	"gobible/logmanager/cli/http/services/search"
	"gobible/logmanager/cli/http/utils"
	"gobible/logmanager/cli/util"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

func Search(res http.ResponseWriter, req *http.Request, params httprouter.Params) {

	bytes, err := json.Marshal(file.Test)
	if err != nil {
		log.Println(err)
	}
	fmt.Fprint(res, string(bytes))
	return
}

func errResultJson(res http.ResponseWriter, msg string, err error) {

	errMsg := ""
	if err != nil {
		errMsg = err.Error()
	}
	content := fmt.Sprintf("%s %s", msg, errMsg)
	rlt := data.NewJson(1, content, nil)
	fmt.Fprint(res, string(rlt))
	//return
}

func getGlobalDirsName(res http.ResponseWriter, dirs *[]string, startTime, endTime string) (err error) {

	ts, err := time.Parse(search.TIMEFORMAT, startTime)
	if err != nil {
		errResultJson(res, "解析开始时间格式错误:", nil)
		return err
	}
	te, err := time.Parse(search.TIMEFORMAT, endTime)
	if err != nil {
		errResultJson(res, "解析结束时间格式错误:", nil)
		return err
	}
	if te.Before(ts) {
		err = errors.New("")
		errResultJson(res, "日期不合法,结束日期比开始日期还早哦:", nil)
		return err
	}
	if ts.Equal(te) { // 日期相等
		*dirs = append(*dirs, startTime)
	} else {
		// 日期大于前者
		*dirs = append(*dirs, startTime)
		ts = ts.Add(time.Hour * 24)
		for te.After(ts) || te.Equal(ts) {
			*dirs = append(*dirs, ts.Format(search.TIMEFORMAT))
			ts = ts.Add(time.Hour * 24)
		}
	}
	logrus.Println(*dirs)
	return
}

func isProcessing(key string) bool {

	value, err := redis.GetValue(key)
	if err != nil { // nil returned 表示key不存在
		//log.Println(err)
		return false
	}
	if len(value) > 0 {
		return true
	}
	return false
}

func genDownloadDirIfInputEmpty() {

	//初始化服务器日志生成的目录；便于查看情况
	if !util.PathExist(util.GetCurrentDirectory() + "/" + config.ZipResultDir) {
		err := os.Mkdir(config.ZipResultDir, 0755)
		if err != nil {
			logrus.Println(err)
		}
	}

}

//获取数据服务
func PickDemo(res http.ResponseWriter, req *http.Request, params httprouter.Params) {

	//data.SetCROS(res)
	res.Header().Add("Access-Control-Allow-Headers", "Content-Type") //header的类型
	//startTime := req.PostFormValue("start")
	////startTime := params.ByName("start")

	content, err := ioutil.ReadAll(req.Body)
	if err != nil {
		logrus.Println(err)
	}
	defer req.Body.Close()
	//log.Println("body :",string(content))

	pickData := data.Pick{}

	err = json.Unmarshal(content, &pickData)
	if err != nil {
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
		rlt := data.NewJson(1, "获取config.json配置失败", nil)
		fmt.Fprint(res, string(rlt))
		return
	}

	down = configDown
	if len(down) == 0 {
		rlt := data.NewJson(1, "config.json down 配置为空", nil)
		fmt.Fprint(res, string(rlt))
		return
	}

	//参数校验
	if len(startTime) == 0 ||
		len(endTime) == 0 ||
		len(condition) == 0 ||
		len(dir) == 0 {
		rlt := data.NewJson(1, "参数不合法", nil)
		fmt.Fprint(res, string(rlt))
		return
	}

	if len(dir) > 0 {
		if !util.PathExist(dir) {
			rlt := data.NewJson(1, "查找路径不存在", nil)
			res.Write([]byte(string(rlt)))
			return
		}
	}

	if len(down) > 0 {
		if !util.PathExist(down) {
			rlt := data.NewJson(1, "下载目录不存在", nil)
			res.Write([]byte(string(rlt)))
			return
		}
	}

	//如果没有down参数则是download目录
	if len(down) == 0 {
		genDownloadDirIfInputEmpty()
	} else {
		//search.ZipDirSignal <- down
	}

	//如果参数合法就判断是否是重复的请求
	composeStr := fmt.Sprintf("%s-%s-%s-%s", startTime, endTime, condition, dir)
	findCondition := utils.MD5(composeStr)
	if isProcessing(findCondition) {
		rlt := data.NewJson(0, "该查找任务已经提交,正在处理中,请问重复提交,谢谢！", nil)
		res.Write([]byte(rlt))
		return
	}

	//格式化的日志列表slice
	dirs := make([]string, 0)
	//参数校验
	err = getGlobalDirsName(res, &dirs, startTime, endTime)
	if err != nil {
		return
	}

	go search.DoSearch(dirs, dir, findCondition, condition, down)

	//提交任务后马上设置值
	redis.SetValue(findCondition, composeStr)
	rlt := data.NewJson(0, "文件处理中", nil)
	res.Write([]byte(rlt))

}

//获取数据服务
func Pick(res http.ResponseWriter, req *http.Request, params httprouter.Params) {

	res.Header().Add("Access-Control-Allow-Headers", "Content-Type") //header的类型
	content, err := ioutil.ReadAll(req.Body)
	if err != nil {
		logrus.Println(err)
	}
	defer req.Body.Close()
	pickData := data.Pick{}

	err = json.Unmarshal(content, &pickData)
	if err != nil {
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
		rlt := data.NewJson(1, "参数不合法", nil)
		fmt.Fprint(res, string(rlt))
		return
	}
	if len(dir) > 0 {
		if !util.PathExist(dir) {
			rlt := data.NewJson(1, "查找路径不存在", nil)
			res.Write([]byte(string(rlt)))
			return
		}
	}
	//如果参数合法就判断是否是重复的请求
	timeNoa := time.Now().Nanosecond()
	composeStr := fmt.Sprintf("%s-%s-%s-%s-%d", startTime, endTime, condition, dir, timeNoa)
	findCondition := utils.MD5(composeStr)

	//格式化的日志列表slice
	dirs := make([]string, 0)
	//参数校验
	err = getGlobalDirsName(res, &dirs, startTime, endTime)
	if err != nil {
		return
	}

	v, err := redis.HGetAll(findCondition)
	if err != nil {
		rlt := data.NewJson(1, "查询失败:"+err.Error(), nil)
		fmt.Fprint(res, string(rlt))
		return
	}

	if v.Status == 0 {
		//不存在的时候才去下发
		err = redis.LPush(config.RedisTaskName, findCondition)
		if err != nil {
			rlt := data.NewJson(1, "下发任务失败", err.Error())
			res.Write([]byte(string(rlt)))
			return
		}
		task := map[string]interface{}{
			"status":    config.RedisStatusNotStart,
			"condition": composeStr,
			"download":  "",
		}
		redis.HMSet(findCondition, task)
		redis.Expire(findCondition, config.OVERTIME)
		msg := fmt.Sprintf("%s,%s", findCondition, "任务下发成功;请去查询接口获取结果文件地址")
		rlt := data.NewJson(1, msg, struct{}{})
		res.Write([]byte(string(rlt)))
		return
	}
	//如果存在该条记录
	if v.Status != int(config.RedisStatusFailure) {
		msg := config.RedisStatus(config.RedisTaskStatus(v.Status))
		rlt := data.NewJson(0, findCondition+","+msg, v.DownLoad)
		fmt.Fprint(res, string(rlt))
		return
	}

}

func DoWork() {

	for {

		listLen, err := redis.LLen(config.RedisTaskName)
		if err != nil {
			log.Println(err)
			time.Sleep(time.Second)
			continue
		}
		if listLen > 0 {

			v, err := redis.RRop(config.RedisTaskName)
			if err != nil {
				log.Println(err)
				continue
			}
			//根据v查找 hash中的条件和状态
			job, err := redis.HGetAll(v)
			if err != nil {
				log.Println(err)
				continue
			}
			//status := job.Status
			RedisCondtion := job.Condition
			fields := strings.Split(RedisCondtion, config.ConditionSplitChar)

			//格式化的日志列表slice
			dirs := make([]string, 0)
			startTime := fields[0]
			endTime := fields[1]

			// 获取需要遍历的目录 控制器部分已经做了合法判断
			ts, err := time.Parse(search.TIMEFORMAT, startTime)
			if err != nil {
				log.Println(err)
				continue
			}
			te, err := time.Parse(search.TIMEFORMAT, endTime)
			if err != nil {
				log.Println(err)
				continue
			}
			if ts.Equal(te) { // 日期相等
				dirs = append(dirs, startTime)
			} else {
				// 日期大于前者
				dirs = append(dirs, startTime)
				ts = ts.Add(time.Hour * 24)
				for te.After(ts) || te.Equal(ts) {
					dirs = append(dirs, ts.Format(search.TIMEFORMAT))
					ts = ts.Add(time.Hour * 24)
				}
			}
			//findCondition := ""
			condition := fields[2]
			dir := fields[3]
			down := search.DownloadDir
			//获取到所有的值

			//log.Println(dirs,dir,down,condition)

			hashKey := v

			//开始处理
			task := map[string]interface{}{
				"status": config.RedisStautsRunning,
				//"condition":composeStr,
				//"download":"",
			}
			//只有第一次设置的时候对过期时间有影响；无需再次设置！
			redis.HMSet(hashKey, task)

			search.DoSearch(dirs, dir, hashKey, condition, down)

		}

		time.Sleep(time.Second)

	}

}
