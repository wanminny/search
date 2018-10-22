package main

import (
	"gobible/logmanager/cli/http/controllers"
	"net/http"
	"github.com/sirupsen/logrus"
	"log"
)

func main()  {

	controllers.InitRouter()

	controllers.InitLog()

	log.Println("service start on :8080,ok!")

	logrus.Fatal(http.ListenAndServe(":8080",controllers.Router))

}
