package handler

import "net/http"
import (
	"log"
	"gobible/logmanager/cli/test/response"
)

type User struct {

}


func AddMiddleWithHandle(next http.Handler) http.Handler{

	f :=  func(w http.ResponseWriter,req *http.Request) {
		log.Print(123)
		v := req.URL.Query()
		user := v["user"]
		if user[0] != "ccc" {
			response.NewJson(1,"auth failed", struct {
			}{})
			return
		}
		log.Println(user)
		next.ServeHTTP(w,req)
	}
	return http.HandlerFunc(f)
}

// User类型实现了ServeHTTP方法
func(u User) ServeHTTP(w http.ResponseWriter,req *http.Request)  {
	w.Write([]byte("user handler...."))
}



