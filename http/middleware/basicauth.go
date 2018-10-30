package middleware

import (
	"net/http"
	"log"
	"github.com/julienschmidt/httprouter"
	"strings"
	"encoding/base64"
	"gobible/logmanager/cli/test/response"
	"gobible/logmanager/cli/http/config"
	"io"
)

const UNAME  = "test"
const UPASSWD  = "test"

func AuthHeader(next http.Handler) http.Handler{
	f :=  func(w http.ResponseWriter,req *http.Request) {
		authH := w.Header().Get("Authentication")
		log.Println(authH)
		next.ServeHTTP(w,req)
	}
	return http.HandlerFunc(f)
}


///因为是 httprouter自己的 处理器（handle）所以不能同上处理了！
func AuthHeaderWithHttpRouter(next httprouter.Handle) httprouter.Handle{
	return  func(w http.ResponseWriter,req *http.Request,params httprouter.Params) {
		authH := req.Header.Get("Authorization")

		if len(authH) == 0 {
			msg := response.NewJson(1,config.AuthFail.Error(),nil)
			io.WriteString(w,msg)
			return
		}
		//否则会crash
		if len(authH) >0 {
			splitStr := strings.SplitN(authH," ",2)
			c,err := base64.StdEncoding.DecodeString(splitStr[1])
			if err != nil{
				log.Println(err)
			}
			token := strings.Split(string(c),":")

			uname := token[0]
			passwd := token[1]
			if uname != UNAME || passwd != UPASSWD{
				msg := response.NewJson(1,config.AuthFail.Error(),nil)
				io.WriteString(w,msg)
				return
			}
		}
		next(w,req,params)
	}


}


