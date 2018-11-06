package middleware

import "net/http"

//跨域
func CorsHeader(next http.Handler) http.Handler {
	f := func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")             //允许访问所有域
		w.Header().Add("Access-Control-Allow-Headers", "Content-Type") //header的类型
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		next.ServeHTTP(w, req)
	}
	return http.HandlerFunc(f)
}
