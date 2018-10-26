package main

import (
	"net/http"
	"log"
	"gobible/logmanager/cli/http/controllers"
	"github.com/julienschmidt/httprouter"
	"github.com/rs/cors"
	"gobible/logmanager/cli/http/utils"
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

	controllers.InitLog()

	log.Println("service start on :8080,ok!")
	
	handler := cors.AllowAll().Handler(Router)

	go utils.DeleteOverSomeTime()

	log.Fatal(http.ListenAndServe(":8080",handler))

}
