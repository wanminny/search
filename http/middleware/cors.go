package middleware

import "net/http"

func CorsHeader(next http.Handler) http.Handler{
	f :=  func(w http.ResponseWriter,req *http.Request) {
		//w.Header().Set("Content-Type","application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")             //允许访问所有域
		w.Header().Add("Access-Control-Allow-Headers", "Content-Type") //header的类型
		w.Header().Set("access-control-allow-credentials","true")
		//res.Header().Set("content-type", "application/json")             //返回数据格式是json
		next.ServeHTTP(w,req)
	}
	return http.HandlerFunc(f)
}


