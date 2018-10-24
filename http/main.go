package main

import (
	"net/http"
	"log"
	"gobible/logmanager/cli/http/controllers"
	"github.com/julienschmidt/httprouter"
	"github.com/rs/cors"
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

	go controllers.SearchdirRouter(Router)

	//c := cors.New(cors.Options{
	//	AllowedOrigins: []string{"*"},
	//	AllowCredentials: true,
	//	AllowedMethods: []string{"GET", "POST", "DELETE", "PUT", "OPTIONS"},
	//	//Debug:true,
	//})
	//跨域设置
	//handler := cors.Default().Handler(Router)

	handler := cors.AllowAll().Handler(Router)
	//log.Fatal(http.ListenAndServe(":8080",c.Handler(Router)))

	log.Fatal(http.ListenAndServe(":8080",handler))

}
