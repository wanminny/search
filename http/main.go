package main

import (
	"flag"
	"github.com/julienschmidt/httprouter"
	"github.com/rs/cors"
	"github.com/sirupsen/logrus"
	"gobible/logmanager/cli/http/cache/redis"
	"gobible/logmanager/cli/http/controllers"
	"gobible/logmanager/cli/http/controllers/data"
	"gobible/logmanager/cli/http/middleware"
	"gobible/logmanager/cli/http/services/search"
	"gobible/logmanager/cli/http/utils"
	"gobible/logmanager/cli/util"
	"log"
	"net/http"
)

var (
	Router *httprouter.Router
	Port   string
)

func GetGlobalRouter() *httprouter.Router {
	Router = httprouter.New()
	return Router
}

//自定义端口
func initEnv() {
	flag.StringVar(&Port, "port", ":8080", "server port :6000")
	flag.Parse()
}

func initConfigJson() {
	configFile := redis.CheckRedisConf()
	configDown := configFile.Down
	downDomain := configFile.DownDomain
	if len(configDown) == 0 {
		log.Fatal("获取config.json配置失败")
	}
	down := configDown
	if len(down) > 0 {
		if !util.PathExist(down) {
			log.Fatal("config.json down 中的下载目录不存在")
		}
	}

	if len(downDomain) == 0 {
		log.Fatal("获取config.json 下载Domain 配置失败")
	}

	search.DownloadDir = down
	search.DownLoadDomain = downDomain
}

func main() {

	//catch 全局异常 ？
	//excepiton.Finally()

	defer func() {
		if err := recover(); err != nil {
			log.Println("异常 :", err)
			logrus.Println(err)
			//http.Error()
		}
	}()

	initEnv()

	GetGlobalRouter()

	controllers.InitRouter(Router)

	initConfigJson()

	controllers.InitLog()

	log.Printf("service start on %s,ok!",Port)

	//跨域支持！
	handler := cors.AllowAll().Handler(Router)

	//对接口进行Json化
	handler = middleware.JsonHeader(handler)

	go utils.DeleteOverSomeTime()

	// 队列处理;
	go data.DoWork()

	//gracehttp.Serve()
	//log.Fatal(grace.ListenAndServe(":8080", handler))

	if len(Port) == 0{
		Port =":8080"
	}else{
		//Port = ":"+Port
	}

	//方法一
	log.Fatal(http.ListenAndServe(Port, handler))

	// 实际上就是替换MUX的过程！ 【默认 	http.DefaultServeMux or 自定义的mux 】
	// 方法二
	// log.Fatal(http.ListenAndServe(Port,Router))


	//  方法三 （另外一种使用方法 : 自定义 server）
	//server := http.Server{
	//	Addr:Port,
	//	Handler:handler,
	//	ErrorLog:log.New(os.Stdout,"my-server",log.LstdFlags),
	//	ReadTimeout: 5 *time.Second,
	//	WriteTimeout:10 *time.Second,
	//	IdleTimeout:15 * time.Second,
	//}
	//
	//log.Fatal(server.ListenAndServe())



	// 方法四
	//或者 这样处理
	//listener, err := net.Listen("tcp",Port)
	//if err != nil{
	//	log.Fatal(err)
	//}
	//log.Fatal(server.Serve(listener))

}
