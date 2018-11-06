package controllers

import (
	"expvar"
	"github.com/julienschmidt/httprouter"
	"github.com/sirupsen/logrus"
	"gobible/logmanager/cli/http/cache/redis"
	"gobible/logmanager/cli/http/config"
	"gobible/logmanager/cli/http/controllers/data"
	"gobible/logmanager/cli/http/controllers/file"
	"gobible/logmanager/cli/http/controllers/jsonp"
	"gobible/logmanager/cli/http/controllers/task"
	"gobible/logmanager/cli/http/middleware"
	"gobible/logmanager/cli/http/services/search"
	"gobible/logmanager/cli/util"
	"log"
	"net/http"
	"os"
	"sync"
)

//取巧；解决不能在控制器中在设置静态路由---》造成循环引用；
// 这里只会执行一次！而且不会导致后续的请求中的阻塞。
func SearchdirRouter(router *httprouter.Router) {

	once := sync.Once{}
	for {
		ZipResultDir := <-search.ZipDirSignal
		//log.Print(ZipResultDir)
		//目录仅仅执行一次
		once.Do(func() {
			//log.Println("once......")
			//ZipResultDir := <-search.ZipDirSignal
			// TODO 测试 如果 ZipResultDir是非本地目录是否可以？
			router.ServeFiles("/log/*filepath", http.Dir(ZipResultDir))
			logrus.Println("zip result set ok: ", config.ZipResultDir)
		})
	}
}

//总路由
func InitRouter(router *httprouter.Router) {

	//主页
	//文件目录
	router.GET("/", file.Content)

	router.GET("/search", data.Search)

	configFile := redis.CheckRedisConf()

	//http://localhost:8080/log/ 可以成功！
	// 注释： *filepath是固定的！否则报错
	//文件目录列表服务
	router.ServeFiles("/log/*filepath", http.Dir(configFile.Down))

	//服务器列表目录 [服务器诊断日志]
	router.ServeFiles("/server_log/*filepath", http.Dir(config.ServerLogDir))

	//下发任务
	//router.POST("/data/pick",data.Pick)

	//下发任务
	router.POST("/data/pick", data.Pick)

	//列出正在处理中的任务
	router.GET("/task", task.List)

	//查询某个任务是否还在运行
	router.GET("/list/:no", task.CheckIsRunning)

	//内存 诊断
	router.GET("/memory", file.Memory)

	router.GET("/jsonp", jsonp.JsonpHandler)

	//某个函数处理 auth
	router.GET("/auth", middleware.AuthHeaderWithHttpRouter(file.Content))

	// 这样使用不管用;并没有并自定义的router 或者是mux 接管；
	//http.Handle("/stat",expvar.Handler())

	// http://localhost:8080/stat [使用的默认的 expvar ]
	// not http://localhost:8080/debug/vars
	router.GET("/stat", middleware.ExpVarHeader1(expvar.Handler()))

}

func initDir() {

	//初始化服务器日志生成的目录；便于查看情况
	//if !util.PathExist(util.GetCurrentDirectory() + "/" + search.ZipResultDir){
	//	err :=os.Mkdir(search.ZipResultDir,0755)
	//	if err != nil{
	//		logrus.Println(err)
	//	}
	//}

	if !util.PathExist(util.GetCurrentDirectory() + "/" + config.TmpTransferDir) {
		err := os.Mkdir(config.TmpTransferDir, 0755)
		if err != nil {
			logrus.Println(err)
		}
	}

	//初始化服务器日志生成的目录；便于查看情况
	if !util.PathExist(util.GetCurrentDirectory() + "/" + config.ServerLogDir) {
		err := os.Mkdir(config.ServerLogDir, 0755)
		if err != nil {
			logrus.Println(err)
		}
	}

}

func InitLog() (err error) {

	f, _ := os.OpenFile(util.GetCurrentDirectory()+"/"+config.ServerLogDir+"/"+"server.log", os.O_WRONLY|os.O_CREATE|os.O_SYNC|os.O_APPEND, 0755)

	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetOutput(f)

	initDir()

	//读取配置并 检查是否可以连上redis服务器
	err = redis.PING()
	if err != nil {
		log.Println(err)
		panic(err)
	}
	return

}
