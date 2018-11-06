package middleware

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func ExpVarHeader(next http.Handler) http.Handler {
	f := func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		next.ServeHTTP(w, req)
	}
	return http.HandlerFunc(f)
}

func ExpVarHeader1(next http.Handler) httprouter.Handle {
	f := func(w http.ResponseWriter, req *http.Request, params httprouter.Params) {
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, req)
	}
	return f
}
