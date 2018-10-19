package file

import (
	"net/http"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"gobible/logmanager/cli/http/models/data"
)

//文件目录服务
func Content(res http.ResponseWriter,req *http.Request,params httprouter.Params)()  {

	rlt := data.NewJson(0,"show result! access path /log",nil)
	fmt.Fprint(res,string(rlt))
	return
}

