package task

import (
	"fmt"
	"github.com/julienschmidt/httprouter"
	"gobible/logmanager/cli/http/cache/redis"
	"gobible/logmanager/cli/http/config"
	"gobible/logmanager/cli/http/models/data"
	"net/http"
)

// 列出正在处理的任务
func List(res http.ResponseWriter, req *http.Request, params httprouter.Params) {

	//data.SetCROS(res)
	res.Header().Add("Access-Control-Allow-Headers", "Content-Type") //header的类型
	m := redis.AllKeys()
	rlt := data.NewJson(0, "所有任务列表", m)
	fmt.Fprint(res, string(rlt))
	return

}

//检查某个任务是否处理完成。
func CheckIsRunning1(res http.ResponseWriter, req *http.Request, params httprouter.Params) {

	//data.SetCROS(res)
	res.Header().Add("Access-Control-Allow-Headers", "Content-Type") //header的类型
	no := params.ByName("no")
	if len(no) == 0 {
		rlt := data.NewJson(1, "参数不合法", nil)
		fmt.Fprint(res, string(rlt))
		return
	}

	v, err := redis.GetValue(no)
	if err != nil {
		rlt := data.NewJson(1, "任务已经处理完成或者不存在该任务："+err.Error(), nil)
		fmt.Fprint(res, string(rlt))
		return
	}

	rlt := data.NewJson(0, "正在处理中,查询的条件是:"+v, nil)
	fmt.Fprint(res, string(rlt))
	return
}

//检查某个任务是否处理完成。
func CheckIsRunning(res http.ResponseWriter, req *http.Request, params httprouter.Params) {

	//data.SetCROS(res)
	res.Header().Add("Access-Control-Allow-Headers", "Content-Type") //header的类型
	no := params.ByName("no")
	if len(no) == 0 {
		rlt := data.NewJson(1, "参数不合法", nil)
		fmt.Fprint(res, string(rlt))
		return
	}

	v, err := redis.HGetAll(no)
	//v,err = redis.HGet(no,"status")
	//v,err := redis.GetValue(no)
	if err != nil {
		rlt := data.NewJson(1, "查询失败:"+err.Error(), nil)
		fmt.Fprint(res, string(rlt))
		return
	}

	if v.Status == 0 {
		rlt := data.NewJson(0, config.HashNotExists.Error(), v.DownLoad)
		fmt.Fprint(res, string(rlt))
		return
	}
	msg := config.RedisStatus(config.RedisTaskStatus(v.Status))
	rlt := data.NewJson(0, msg, v.DownLoad)
	fmt.Fprint(res, string(rlt))
	return
}
