package handleFunc

import (
	"net/http"
	"io"
	"fmt"
	"log"
)

// 通过handlerFunc方式添加中间件
func AddMiddleWareWithHanlderFunc(next http.HandlerFunc) http.HandlerFunc  {

	return func(writer http.ResponseWriter, request *http.Request) {
		//获取get参数
		query := request.URL.Query()
		user := query["hobby"]

		//获取post参数

		//Form：存储了post、put和get参数，在使用之前需要调用ParseForm方法。
		//
		//PostForm：存储了post、put参数，在使用之前需要调用ParseForm方法。
		//
		//MultipartForm：存储了包含了文件上传的表单的post参数，在使用前需要调用ParseMultipartForm方法
		if err := request.ParseForm(); err != nil {
			log.Print(err)
		}

		//[{"key":"Content-Type","name":"Content-Type","value":"application/x-www-form-urlencoded","description":"","type":"text"}]
		log.Println(request.PostFormValue("user"))

		fmt.Println(user)
		next(writer,request)
	}
}

// 简单的添加 普通方法;
func CommonFunc(w http.ResponseWriter,req * http.Request)  {
	log.Println("common func with handler ")
	//time.Sleep(time.Second)
	io.WriteString(w,"common func with handler func !")
}

func Log(w http.ResponseWriter,req *http.Request)  {
	log.Println("log here.....")
	io.WriteString(w,"log func")
}

// 防止普通方式一个调用嵌套一个调用 链太长
func Use(httpFunc http.HandlerFunc,middle ...func(httpFunc http.HandlerFunc) http.HandlerFunc) http.HandlerFunc  {

	for _,v := range middle {
		httpFunc = v(httpFunc)
	}
	return httpFunc
}