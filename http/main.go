package main

import (
	"gobible/logmanager/cli/http/controllers"
	"log"
	"net/http"
)

func main()  {

	controllers.InitRouter()

	log.Println("service start on :8080,ok!")

	log.Fatal(http.ListenAndServe(":8080",controllers.Router))

}
