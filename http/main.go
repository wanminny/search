package main

import (
	"net/http"
	"log"
	"gobible/logmanager/cli/http/controllers"
	"github.com/julienschmidt/httprouter"
)

var (
	Router  *httprouter.Router
)

func GetGlobalRouter() *httprouter.Router  {
	Router = httprouter.New()
	return Router
}


func main()  {

	GetGlobalRouter()

	controllers.InitRouter(Router)

	go controllers.SearchdirRouter(Router)

	controllers.InitLog()

	log.Println("service start on :8080,ok!")

	log.Fatal(http.ListenAndServe(":8080",Router))

}
