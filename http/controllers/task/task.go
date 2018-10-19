package task

import (
	"net/http"
	"github.com/julienschmidt/httprouter"
	"gobible/logmanager/cli/http/cache/redis"
	"gobible/logmanager/cli/http/models/data"
	"fmt"
)

// 列出正在处理的任务
func List(res http.ResponseWriter,req *http.Request,params httprouter.Params)  {

	m := redis.AllKeys()
	rlt := data.NewJson(0,"所有任务列表",m)
	fmt.Fprint(res,string(rlt))
	return

}


