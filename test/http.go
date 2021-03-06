package main

import (
	"net/http"
	"gobible/logmanager/cli/test/middle/handler"
	"gobible/logmanager/cli/test/middle/handleFunc"
)

//
func httpDemo()  {

	//http 的添加handler的两种方式,一种是handle  一种是handleFunc
	f := func(w http.ResponseWriter,req *http.Request) {

	}
	http.Handle("/bb",http.HandlerFunc(f))

	http.HandleFunc("/aa", func(writer http.ResponseWriter, request *http.Request) {

	})

}

func main()  {

	//添加中间件1
	http.Handle("/testhandle",handler.AddMiddleWithHandle(handler.User{}))

	//添加中间件2
	http.HandleFunc("/testhandlefunc",handleFunc.AddMiddleWareWithHanlderFunc(handleFunc.CommonFunc))


	//使用便捷的方式调用 处理
	http.Handle("/usehandle",handler.Use(
		http.HandlerFunc(handleFunc.CommonFunc),
		handler.LogFunc,handler.AddMiddleWithHandle))


	// 使用这样一种就可以了 使用栈的方式执行顺序！
	http.Handle("/usehandle2",handler.Use(
		http.HandlerFunc(handleFunc.CommonFunc),
		handler.AddMiddleWithHandle,handler.LogFunc))


	http.Handle("/usehandle1",handler.Use(
		handler.User{},
		handler.LogFunc,handler.AddMiddleWithHandle))

	http.ListenAndServe(":8081",nil)
}
