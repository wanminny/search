package middleware

import (
	"net/http"
)

func JsonHeader(next http.Handler) http.Handler{
	f :=  func(w http.ResponseWriter,req *http.Request) {
		w.Header().Set("Content-Type","application/json")
		next.ServeHTTP(w,req)
	}
	return http.HandlerFunc(f)
}

