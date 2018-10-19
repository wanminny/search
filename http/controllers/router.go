package controllers

import (
	"github.com/julienschmidt/httprouter"
	"gobible/logmanager/cli/http/controllers/data"
	"gobible/logmanager/cli/http/controllers/task"
	"gobible/logmanager/cli/http/controllers/file"
	"net/http"
)

var (
	Router = httprouter.New()
)

//总路由
func InitRouter()  {

	//主页
	Router.GET("/",data.Search)

	//文件目录
	Router.GET("/catalog",file.Content)

	//http://localhost:8080/log/ 可以成功！
	// 注释： *filepath是固定的！否则报错
	//文件目录列表服务
	Router.ServeFiles("/log/*filepath",http.Dir("download"))

	//下发任务
	Router.POST("/data/pick",data.Pick)

	//列出正在处理中的任务
	Router.GET("/task",task.List)

	//查询某个任务是否还在运行
	Router.GET("/list/:no",task.CheckIsRunning)

}
