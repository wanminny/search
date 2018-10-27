package main

import (
	"net/http"
	"log"
	"gobible/logmanager/cli/http/controllers"
	"github.com/julienschmidt/httprouter"
	"github.com/rs/cors"
	"gobible/logmanager/cli/http/utils"
	"flag"
	"gobible/logmanager/cli/http/cache/redis"
	"gobible/logmanager/cli/util"
	"gobible/logmanager/cli/http/services/search"
	"gobible/logmanager/cli/http/controllers/data"
)

var (
	Router  *httprouter.Router
	Port string
)

func GetGlobalRouter() *httprouter.Router  {
	Router = httprouter.New()
	return Router
}

//自定义端口
func initEnv()  {
	flag.StringVar(&Port,"port",":8080","server port .")
	flag.Parse()
}

func initConfigJson()  {
	configFile := redis.CheckRedisConf()
	configDown := configFile.Down
	if len(configDown) == 0 {
		log.Fatal("获取config.json配置失败")
	}
	down := configDown
	if len(down) > 0 {
		if !util.PathExist(down){
			log.Fatal("config.json down 中的下载目录不存在")
		}
	}

	search.DownloadDir = down
}

func main()  {

	initEnv()

	GetGlobalRouter()

	controllers.InitRouter(Router)

	initConfigJson()

	controllers.InitLog()

	log.Println("service start on :8080,ok!")

	handler := cors.AllowAll().Handler(Router)

	go utils.DeleteOverSomeTime()

	// 队列处理;
	go data.DoWork()

	log.Fatal(http.ListenAndServe(":8080",handler))

}
