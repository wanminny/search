package handler

import "net/http"
import (
	"log"
	"gobible/logmanager/cli/test/response"
)

type User struct {

}

type Auth struct {

}

func AddMiddleWithHandle(next http.Handler) http.Handler{

	f :=  func(w http.ResponseWriter,req *http.Request) {
		log.Print(123)
		v := req.URL.Query()
		user := v["user"]
		if len(user) > 0 && user[0] != "ccc" {
			response.NewJson(1,"auth failed", struct {
			}{})
			return
		}
		log.Println(user)
		next.ServeHTTP(w,req)
	}
	return http.HandlerFunc(f)
}

func LogFunc(next http.Handler) http.Handler{

	f :=  func(w http.ResponseWriter,req *http.Request) {
		log.Print("start.......")
		//time.Sleep(time.Second)
		log.Println("end....")
		next.ServeHTTP(w,req)
	}
	return http.HandlerFunc(f)
}
// User类型实现了ServeHTTP方法
func(u User) ServeHTTP(w http.ResponseWriter,req *http.Request)  {
	w.Write([]byte("user handler...."))
}

// Auth 类型实现了ServeHTTP方法
func(auth Auth) ServeHTTP(w http.ResponseWriter,req *http.Request)  {
	w.Write([]byte("auth handler...."))
}

//简化专中间件的调用方式! 防止过长  使用这样一种就可以了
func Use(h http.Handler,middle ...func(http http.Handler)http.Handler) http.Handler  {

	for _,v := range middle{
		h = v(h)
	}
	return h
}
