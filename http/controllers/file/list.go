package file

import (
	"net/http"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"gobible/logmanager/cli/http/models/data"
)

//文件目录服务
func Content(res http.ResponseWriter,req *http.Request,params httprouter.Params)()  {

	res.Header().Add("Access-Control-Allow-Headers", "Content-Type") //header的类型
	//data.SetCROS(res)
	rlt := data.NewJson(0,"show result! access path /log",nil)
	fmt.Fprint(res,string(rlt))
	return
}

