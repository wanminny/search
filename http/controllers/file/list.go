package file

import (
	"net/http"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"gobible/logmanager/cli/http/models/data"
	"encoding/json"
	"github.com/qiniu/log"
)



var (
	Test map[string]string = make(map[string]string,2)

)


//文件目录服务
func Content(res http.ResponseWriter,req *http.Request,params httprouter.Params)()  {

	Test["aaa"] = "test"
	//res.Header().Add("Access-Control-Allow-Headers", "Content-Type") //header的类型
	data.SetCROS(res)
	rlt := data.NewJson(0,"show result! access path /log",nil)
	fmt.Fprint(res,string(rlt))
	return
}

//内存诊断
func Memory(res http.ResponseWriter,req *http.Request,params httprouter.Params)()  {
	res.Header().Add("Access-Control-Allow-Headers", "Content-Type") //header的类型
	//data.SetCROS(res)
	bytes,err := json.Marshal(Test)
	if err != nil{
		log.Println(err)
	}
	fmt.Fprint(res,string(bytes))
	return

}