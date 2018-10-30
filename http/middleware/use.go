package middleware

import "net/http"

//简化专中间件的调用方式! 防止过长  使用这样一种就可以了
func Use(h http.Handler,middle ...func(http http.Handler)http.Handler) http.Handler  {

	for _,v := range middle{
		h = v(h)
	}
	return h
}


