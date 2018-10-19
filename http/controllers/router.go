package controllers

import (
	"github.com/julienschmidt/httprouter"
	"gobible/logmanager/cli/http/controllers/data"
)

var (
	Router = httprouter.New()
)

//总路由
func InitRouter()  {

	Router.GET("/data",data.Search)

	Router.POST("/data/pick",data.Pick)

}
