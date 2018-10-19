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


//检查某个任务是否处理完成。
func CheckIsRunning(res http.ResponseWriter,req *http.Request,params httprouter.Params)  {

	no := params.ByName("no")

	if len(no) == 0{
		rlt := data.NewJson(1,"参数不合法",nil)
		fmt.Fprint(res,string(rlt))
		return
	}

	v,err := redis.GetValue(no)
	if err != nil{
		//log.Println(err)
		rlt := data.NewJson(1, "任务已经处理完成或者不存在该任务：" + err.Error(),nil)
		fmt.Fprint(res,string(rlt))
		return
	}

	rlt := data.NewJson(0,"正在处理中,查询的条件是:"+ v,nil)
	fmt.Fprint(res,string(rlt))
	return
}



