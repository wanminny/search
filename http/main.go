package main

import (
	"gobible/logmanager/cli/http/controllers"
	"net/http"
	"log"
)

func main()  {

	controllers.InitRouter()

	controllers.InitLog()

	log.Println("service start on :8080,ok!")

	log.Fatal(http.ListenAndServe(":8080",controllers.Router))

}
