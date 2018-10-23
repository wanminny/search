package controllers

import (
	"github.com/julienschmidt/httprouter"
	"gobible/logmanager/cli/http/controllers/data"
	"gobible/logmanager/cli/http/controllers/task"
	"gobible/logmanager/cli/http/controllers/file"
	"net/http"
	"os"
	"github.com/sirupsen/logrus"
	"gobible/logmanager/cli/util"
	"gobible/logmanager/cli/http/services/search"
	"gobible/logmanager/cli/http/cache/redis"
)

var (
	Router = httprouter.New()
)

//总路由
func InitRouter()  {

	//主页
	//文件目录
	Router.GET("/",file.Content)

	//http://localhost:8080/log/ 可以成功！
	// 注释： *filepath是固定的！否则报错
	//文件目录列表服务
	Router.ServeFiles("/log/*filepath",http.Dir(search.ZipResultDir))

	//服务器列表目录 [服务器诊断日志]
	Router.ServeFiles("/server_log/*filepath",http.Dir(search.ServerLogDir))

	//下发任务
	Router.POST("/data/pick",data.Pick)

	//列出正在处理中的任务
	Router.GET("/task",task.List)

	//查询某个任务是否还在运行
	Router.GET("/list/:no",task.CheckIsRunning)

}

func initDir()  {

	//初始化服务器日志生成的目录；便于查看情况
	if !util.PathExist(util.GetCurrentDirectory() + "/" + search.ZipResultDir){
		err :=os.Mkdir(search.ZipResultDir,0755)
		if err != nil{
			logrus.Println(err)
		}
	}

	//初始化服务器日志生成的目录；便于查看情况
	if !util.PathExist(util.GetCurrentDirectory() + "/" + search.ServerLogDir){
		err :=os.Mkdir(search.ServerLogDir,0755)
		if err != nil{
			logrus.Println(err)
		}
	}

}

func InitLog()  {

	f, _ := os.OpenFile(util.GetCurrentDirectory() + "/" + search.ServerLogDir + "/"+"server.log", os.O_WRONLY|os.O_CREATE|os.O_SYNC|os.O_APPEND,0755)

	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetOutput(f)

	initDir()

	//读取配置并 检查是否可以连上redis服务器
	redis.PING()

}