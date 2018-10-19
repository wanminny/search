package controllers

import (
	"github.com/julienschmidt/httprouter"
	"gobible/logmanager/cli/http/controllers/data"
	"gobible/logmanager/cli/http/controllers/task"
)

var (
	Router = httprouter.New()
)

//总路由
func InitRouter()  {

	Router.GET("/data",data.Search)
	
	//下发任务
	Router.POST("/data/pick",data.Pick)

	//列出正在处理中的任务
	Router.GET("/task",task.List)

	//查询某个任务是否还在运行
	Router.GET("/list/:no",task.CheckIsRunning)

}
