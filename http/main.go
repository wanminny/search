package main

import (
	"gobible/logmanager/cli/http/controllers"
	"log"
	"net/http"
)

func main()  {

	controllers.InitRouter()
	log.Fatal(http.ListenAndServe(":8080",controllers.Router))

}
