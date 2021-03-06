package file

import (
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"github.com/qiniu/log"
	"gobible/logmanager/cli/http/models/data"
	"net/http"
)

var (
	Test map[string]string = make(map[string]string, 2)
)

//文件目录服务
func Content(res http.ResponseWriter, req *http.Request, params httprouter.Params) {

	Test["aaa"] = "test"
	//res.Header().Add("Access-Control-Allow-Headers", "Content-Type") //header的类型
	data.SetCROS(res)
	//nil 与 struct{}{}的区别
	rlt := data.NewJson(0, "show result! access path /log", struct{}{})
	fmt.Fprint(res, string(rlt))
	return
}

//内存诊断
func Memory(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
	res.Header().Add("Access-Control-Allow-Headers", "Content-Type") //header的类型
	//data.SetCROS(res)
	bytes, err := json.Marshal(Test)
	if err != nil {
		log.Println(err)
	}
	fmt.Fprint(res, string(bytes))
	return

}
